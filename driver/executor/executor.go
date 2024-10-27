// Copyright 2024 Fantom Foundation
// This file is part of Norma System Testing Infrastructure for Sonic.
//
// Norma is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Norma is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Norma. If not, see <http://www.gnu.org/licenses/>.

package executor

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/parser"
	"github.com/Fantom-foundation/go-opera/cmd/sonictool/chain"
	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	pq "github.com/jupp0r/go-priority-queue"
)

// Run executes the given scenario on the given network using the provided clock
// as a time source. Execution will fail (fast) if the scenario is not valid (see
// Scenario's Check() function).
func Run(clock Clock, network driver.Network, scenario *parser.Scenario, outputDir string, epochTracker map[monitoring.Node]string) error {
	if err := scenario.Check(); err != nil {
		return err
	}

	queue := newEventQueue()

	// Schedule end of simulation as a dummy event.
	endTime := Seconds(scenario.Duration)
	queue.add(toSingleEvent(endTime, "shutdown", func() error {
		return nil
	}))

	// Schedule all operations listed in the scenario.
	for _, node := range scenario.Nodes {
		scheduleNodeEvents(&node, queue, network, endTime, outputDir, epochTracker)
	}
	for _, app := range scenario.Applications {
		if err := scheduleApplicationEvents(&app, queue, network, endTime); err != nil {
			return err
		}
	}
	for _, cheat := range scenario.Cheats {
		scheduleCheatEvents(&cheat, queue, network, endTime)
	}

	// Register a handler for Ctrl+C events.
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, os.Interrupt)
	defer signal.Stop(abort)

	// restart clock as network initialization could time considerable amount of time.
	clock.Restart()
	// Run all events.
	for !queue.empty() {
		event := queue.getNext()
		if event == nil {
			break
		}

		// Wait until the event is going to occure ...
		select {
		case <-clock.NotifyAt(event.time()):
			// continue processing
		case <-abort:
			// abort processing
			log.Printf("Received user abort, ending execution ...")
			return fmt.Errorf("aborted by user")
		}

		delay := clock.Delay(event.time())
		// display delay if it exceeds over 1 second
		if delay > time.Second {
			log.Printf("processing '%s' at time %v (delay: %v)...\n", event.name(), event.time(), delay.Round(time.Second/10).Seconds())
		} else {
			log.Printf("processing '%s' at time %v...\n", event.name(), event.time())
		}

		// Execute the event and schedule successors.
		successors, err := event.run()
		if err != nil {
			return err
		}
		queue.addAll(successors)
	}

	return nil
}

// event is a single action required to happen at (approximately) a given time.
type event interface {
	// The time at which the event is to be processed.
	time() Time
	// A short name describing the event for logging.
	name() string
	// Executes the event's action, potentially triggering successor events.
	run() ([]event, error)
}

// eventQueue is a type-safe wrapper of a priority queue to organize events
// to be scheduled and executed during a scenario run.
type eventQueue struct {
	queue pq.PriorityQueue
}

func newEventQueue() *eventQueue {
	return &eventQueue{pq.New()}
}

func (q *eventQueue) empty() bool {
	return q.queue.Len() == 0
}

func (q *eventQueue) add(event event) {
	q.queue.Insert(event, float64(event.time()))
}

func (q *eventQueue) addAll(events []event) {
	for _, event := range events {
		q.add(event)
	}
}

func (q *eventQueue) getNext() event {
	res, err := q.queue.Pop()
	if err != nil {
		log.Printf("Warning: event queue error encountered: %v", err)
		return nil
	}
	return res.(event)
}

// genericEvent is an implementation of an event combining an action-defining
// lambda with a time stamp determining its execution time.
type genericEvent struct {
	eventTime Time
	eventName string
	action    func() ([]event, error)
}

func (e *genericEvent) time() Time {
	return e.eventTime
}

func (e *genericEvent) name() string {
	return e.eventName
}

func (e *genericEvent) run() ([]event, error) {
	return e.action()
}

func toEvent(time Time, name string, action func() ([]event, error)) event {
	return &genericEvent{time, name, action}
}

func toSingleEvent(time Time, name string, action func() error) event {
	return toEvent(time, name, func() ([]event, error) {
		return nil, action()
	})
}

// scheduleNodeEvents schedules a number of events covering the life-cycle of a class of
// nodes during the scenario execution. The nature of the scheduled nodes is taken from the
// given node description, and actions are applied to the given network.
// Node Lifecycle: create -> timer sim events {start, end, kill, restart} -> remove
func scheduleNodeEvents(node *parser.Node, queue *eventQueue, net driver.Network, end Time, outputDir string, epochTracker map[monitoring.Node]string) {
	instances := 1
	if node.Instances != nil {
		instances = *node.Instances
	}
	startTime := Time(0)
	if node.Start != nil {
		startTime = Seconds(*node.Start)
	}
	endTime := end
	if node.End != nil {
		endTime = Seconds(*node.End)
	}
	nodeImportEvent := ""
	if node.Event.Import != nil {
		nodeImportEvent = *node.Event.Import
	}
	nodeExportEvent := ""
	if node.Event.Export != nil {
		nodeExportEvent = *node.Event.Export
	}
	nodeImportInitialGenesis := ""
	if node.Genesis.ImportInitial != nil {
		nodeImportInitialGenesis = *node.Genesis.ImportInitial
	}
	// export by default at <output>/genesis
	nodeExportInitialGenesis := filepath.Join(outputDir, "genesis")
	if node.Genesis.ExportInitial != nil {
		nodeExportInitialGenesis = *node.Genesis.ExportInitial
	}
	nodeExportFinalGenesis := ""
	if node.Genesis.ExportFinal != nil {
		nodeExportFinalGenesis = *node.Genesis.ExportFinal
	}
	nodeMount := ""
	if node.Mount != nil {
		if *node.Mount == "tmp" { // shorthand to bundle this into outputdir
			nodeMount = outputDir
		} else {
			nodeMount = *node.Mount
		}
	}

	for i := 0; i < instances; i++ {
		name := fmt.Sprintf("%s-%d", node.Name, i)
		var instance = new(driver.Node)

		// make sure all mount points are created
		var pathToDatadir string = ""
		if nodeMount != "" {
			pathToDatadir = filepath.Join(nodeMount, name)
		}

		var pathToGenesis string = ""
		if nodeExportInitialGenesis != "" {
			pathToDatadir = filepath.Join(nodeExportInitialGenesis, name)
		}

		for _, path := range []string{pathToDatadir, pathToGenesis} {
			if path != "" {
				os.MkdirAll(path, os.ModePerm)
			}
		}

		// 1. Queue Creation of Node
		// create node -> import genesis if any -> import event if any
		var importEvent event = toSingleEvent(
			startTime,
			fmt.Sprintf("[%s] Check Import Event", name),
			func() error {
				if nodeImportEvent != "" {
					fmt.Sprintf("NOT IMPLEMENTED ! [%s] Importing event from %s\n", name, nodeImportEvent)
				}
				return nil
			},
		)

		var importGenesis event = toEvent(
			startTime,
			fmt.Sprintf("[%s] Check Import Genesis", name),
			func() ([]event, error) {
				if nodeImportInitialGenesis != "" {
					fmt.Sprintf("NOT IMPLEMENTED ! [%s] Importing genesis from %s\n", name, nodeImportInitialGenesis)
				}
				return []event{importEvent}, nil
			},
		)

		var nodeCreate event = toEvent(
			startTime,
			fmt.Sprintf("[%s] Creating node", name),
			func() ([]event, error) {
				var mountGenesis *string = nil
				if pathToGenesis != "" {
					mountGenesis = &pathToGenesis
				}

				var mountDatadir *string = nil
				if pathToDatadir != "" {
					mountDatadir = &pathToDatadir
				}

				newNode, err := net.CreateNode(&driver.NodeConfig{
					Name:         name,
					Validator:    node.IsValidator(),
					Cheater:      node.IsCheater(),
					MountDatadir: mountGenesis,
					MountGenesis: mountDatadir,
				})

				*instance = newNode
				return []event{importGenesis}, err
			},
		)

		queue.add(nodeCreate)

		// 2. Queue Timer SimEvents
		if &node.Timer != nil {
			for timing, evt := range node.Timer {
				switch evt {
				case "start":
					queue.add(toSingleEvent(
						Seconds(timing),
						fmt.Sprintf("[%s] Starting node", name),
						func() error {
							_, err := net.StartNode(*instance)
							return err
						},
					))
				case "end":
					queue.add(toSingleEvent(
						Seconds(timing),
						fmt.Sprintf("[%s] Ending node", name),
						func() error {
							if instance == nil {
								return nil
							}
							if err := net.RemoveNode(*instance); err != nil {
								return err
							}
							if err := (*instance).Stop(); err != nil {
								return err
							}
							return nil
						},
					))
				case "kill":
					queue.add(toSingleEvent(
						Seconds(timing),
						fmt.Sprintf("[%s] Killing node", name),
						func() error {
							return net.KillNode(*instance)
						},
					))
				case "restart":
					queue.add(toEvent(
						Seconds(timing),
						fmt.Sprintf("[%s] Restarting node, ending", name),
						func() ([]event, error) {
							if instance == nil {
								return []event{}, nil
							}
							if err := net.RemoveNode(*instance); err != nil {
								return []event{}, err
							}
							if err := (*instance).Stop(); err != nil {
								return []event{}, err
							}
							return []event{
								toSingleEvent(
									Seconds(timing)+30, // 30 seconds grace period
									fmt.Sprintf("[%s] Restarting node, starting", name),
									func() error {
										_, err := net.StartNode(*instance)
										return err
									},
								),
							}, nil
						},
					))
				}
			}
		}

		// 3. Queue Removal of Node
		// stop node -> export event if any -> export genesis if any -> remove node
		var nodeRemove event = toSingleEvent(
			endTime,
			fmt.Sprintf("[%s] Removing node", name),
			func() error {
				if instance == nil {
					return nil
				}
				if err := (*instance).Cleanup(); err != nil {
					return err
				}
				return nil
			},
		)

		// export genesis if any
		var exportFinalGenesis event = toEvent(
			endTime,
			fmt.Sprintf("[%s] Export Genesis", name),
			func() ([]event, error) {
				if nodeExportFinalGenesis != "" {
					fmt.Printf("Not Implemented ! [%s] Exporting genesis to %s\n", name, nodeExportFinalGenesis)
				}
				return []event{nodeRemove}, nil
			},
		)

		// export event if any
		var exportEvent event = toEvent(
			endTime,
			fmt.Sprintf("[%s] Export Event", name),
			func() ([]event, error) {
				if nodeExportEvent != "" && pathToDatadir != "" {
					pathToOutput := filepath.Join(nodeMount, fmt.Sprintf("%s_%s", name, nodeExportEvent))
					f, err := os.OpenFile(pathToOutput, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
					if err != nil {
						return nil, err
					}
					defer f.Close()

					var writer io.Writer = f
					if strings.HasSuffix(pathToOutput, ".gz") {
						writer = gzip.NewWriter(writer)
						defer writer.(*gzip.Writer).Close()
					}

					if epochTracker == nil {
						return nil, fmt.Errorf("failed to export events; epochTracker == nil")
					}

					if id == nil {
						return nil, fmt.Errorf("failed to export events; id == nil")
					}

					ep, exists := epochTracker[monitoring.Node(name)]
					if !exists {
						return nil, fmt.Errorf("failed to export events; failed to track %d in epochTracker %v\n", name, epochTracker)
					}

					epoch, err := strconv.ParseInt(ep, 10, 32)
					if err != nil {
						return nil, fmt.Errorf("failed to export events; failed to convert epoch to int; %w\n", err)
					}

					fmt.Printf("[%s] Exporting events up to epoch %d to path %s\n", name, epoch, path)
					err = chain.ExportEvents(writer, nodeMount, idx.Epoch(1), idx.Epoch(epoch))
					if err != nil {
						return nil, fmt.Errorf("failed to export events; %w\n", err)
					}
				}
				return []event{exportFinalGenesis}, nil
			},
		)

		// stop node
		var stopNode event = toEvent(
			endTime,
			fmt.Sprintf("[%s] Stop Node", name),
			func() ([]event, error) {
				if instance == nil {
					return nil, nil
				}
				if err := net.RemoveNode(*instance); err != nil {
					return nil, err
				}
				if err := (*instance).Stop(); err != nil {
					return nil, err
				}
				return []event{exportEvent}, nil
			},
		)

		queue.add(stopNode)
	}
}

// scheduleApplicationEvents schedules a number of events covering the life-cycle of a class of
// applications during the scenario execution. The nature of the scheduled applications is taken from the
// given application description, and actions are applied to the given network.
func scheduleApplicationEvents(source *parser.Application, queue *eventQueue, net driver.Network, end Time) error {
	instances := 1
	if source.Instances != nil {
		instances = *source.Instances
	}
	users := 1
	if source.Users != nil {
		users = *source.Users
	}
	startTime := Time(0)
	if source.Start != nil {
		startTime = Seconds(*source.Start)
	}
	endTime := end
	if source.End != nil {
		endTime = Seconds(*source.End)
	}

	for i := 0; i < instances; i++ {
		name := fmt.Sprintf("%s-%d", source.Name, i)
		// TODO add deployment time of contract to config
		queue.add(toSingleEvent(Seconds(5), fmt.Sprintf("deploying contract app %s", name), func() error {
			return startApplication(net, source, name, users, startTime, endTime, queue)
		}))
	}
	return nil
}

// startApplication creates and starts a new application on the network.
func startApplication(net driver.Network, source *parser.Application, name string, users int, startTime, endTime Time, queue *eventQueue) error {
	if newApp, err := net.CreateApplication(&driver.ApplicationConfig{
		Name:  name,
		Type:  source.Type,
		Rate:  &source.Rate,
		Users: users,
	}); err == nil { // schedule application only when it could be created
		queue.add(toSingleEvent(startTime, fmt.Sprintf("starting app %s", name), func() error {
			return newApp.Start()
		}))
		queue.add(toSingleEvent(endTime, fmt.Sprintf("stopping app %s", name), func() error {
			return newApp.Stop()
		}))
	} else {
		return err
	}
	return nil
}

// scheduleCheatEvents schedules a number of events covering the life-cycle of a class of
// cheats during the scenario execution. Currently, a cheat is defined a simultaneous start
// of multiple validator nodes with the same key.
func scheduleCheatEvents(cheat *parser.Cheat, queue *eventQueue, net driver.Network, end Time) {
	startTime := Time(0)
	if cheat.Start != nil {
		startTime = Seconds(*cheat.Start)
	}

	queue.add(toSingleEvent(startTime, fmt.Sprintf("Attempting Cheat %s - currently unsupported cheat, nothing happens", cheat.Name), func() error {
		return nil
	}))
}
