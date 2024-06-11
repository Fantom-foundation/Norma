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
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/parser"
	pq "github.com/jupp0r/go-priority-queue"
)

// Run executes the given scenario on the given network using the provided clock
// as a time source. Execution will fail (fast) if the scenario is not valid (see
// Scenario's Check() function).
func Run(clock Clock, network driver.Network, scenario *parser.Scenario) error {
	if err := scenario.Check(); err != nil {
		return err
	}

	queue := newEventQueue()

	// Schedule end of simulation as a dummy event.
	endTime := Seconds(scenario.Duration)
	queue.add(toSingleEvent(endTime, "shutdown", func() error {
		return nil
	}))

	// Schedule validation at the end of the simulation
	if scenario.Validate {
		queue.add(toSingleEvent(endTime, "[NOT IMPLEMENTED] validate", func() error {
			return nil
		}))
	}

	// Schedule all operations listed in the scenario.
	for _, node := range scenario.Nodes {
		scheduleNodeEvents(&node, queue, network, endTime)
	}
	for _, app := range scenario.Applications {
		if err := scheduleApplicationEvents(&app, queue, network, endTime); err != nil {
			return err
		}
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

		log.Printf("processing '%s' at time %v ...\n", event.name(), event.time())

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
func scheduleNodeEvents(node *parser.Node, queue *eventQueue, net driver.Network, end Time) {
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
	for i := 0; i < instances; i++ {
		name := fmt.Sprintf("%s-%d", node.Name, i)
		var instance = new(driver.Node)
		queue.add(toSingleEvent(startTime, fmt.Sprintf("starting node %s", name), func() error {
			newNode, err := net.CreateNode(&driver.NodeConfig{
				Name: name,
			})
			*instance = newNode
			return err
		}))
		queue.add(toSingleEvent(endTime, fmt.Sprintf("stopping node %s", name), func() error {
			if instance == nil {
				return nil
			}
			if err := net.RemoveNode(*instance); err != nil {
				return err
			}
			if err := (*instance).Stop(); err != nil {
				return err
			}
			return (*instance).Cleanup()
		}))
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
	}
	return nil
}
