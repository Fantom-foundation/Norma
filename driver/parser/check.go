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

package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/Fantom-foundation/Norma/load/app"
)

const namePatternStr = "^[A-Za-z0-9-]+$"

var namePattern = regexp.MustCompile(namePatternStr)

// Check tests semantic constraints on the configuration of a scenario.
func (s *Scenario) Check() error {
	errs := []error{}
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, fmt.Errorf("scenario name must not be empty"))
	}
	if s.Duration <= 0 {
		errs = append(errs, fmt.Errorf("scenario duration must be > 0"))
	}
	if s.NumValidators != nil && *s.NumValidators != 0 {
		errs = append(errs, fmt.Errorf("scenario contains deprecated expression NumValidator"))
	}
	if err := s.checkValidatorConstraints(); err != nil {
		errs = append(errs, err)
	}

	names := map[string]bool{}
	for _, node := range s.Nodes {
		if err := node.Check(s); err != nil {
			errs = append(errs, err)
		}
		if _, exists := names[node.Name]; exists {
			errs = append(errs, fmt.Errorf("node names must be unique, %s encountered multiple times", node.Name))
		} else {
			names[node.Name] = true
		}
	}
	names = map[string]bool{}
	for _, application := range s.Applications {
		if err := application.Check(s); err != nil {
			errs = append(errs, err)
		}
		if _, exists := names[application.Name]; exists {
			errs = append(errs, fmt.Errorf("application names must be unique, %s encountered multiple times", application.Name))
		} else {
			names[application.Name] = true
		}
	}
	names = map[string]bool{}
	for _, cheat := range s.Cheats {
		if err := cheat.Check(s); err != nil {
			errs = append(errs, err)
		}
		if _, exists := names[cheat.Name]; exists {
			errs = append(errs, fmt.Errorf("cheat names must be unique, %s encountered multiple times", cheat.Name))
		} else {
			names[cheat.Name] = true
		}
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the node configuration of a scenario.
func (n *Node) Check(scenario *Scenario) error {
	errs := []error{}
	if !namePattern.Match([]byte(n.Name)) {
		errs = append(errs, fmt.Errorf("node name must match %v, got %v", namePatternStr, n.Name))
	}
	if n.Instances != nil && *n.Instances < 0 {
		errs = append(errs, fmt.Errorf("number of instances must be >= 0, is %d", *n.Instances))
	}
	if n.Client.Type == "" {
		n.Client.Type = "observer"
	}
	if n.Timer == nil {
		n.Timer = make(map[float32]string, 10)
	}

	// Event import/export, Genesis import/export are being refactored.
	// The check "checkTimeNodeAlive" is now obsolete and thus removed.
	// TODO: Remove this comment once refactoring is completed and
	// Event import/export Genesis import/export check is in place.
	if n.Genesis.Import != "" {
		if err := isGenesisFile(n.Genesis.Import, true); err != nil {
			errs = append(errs, err)
		}
	}

	if n.Genesis.Export != "" {
		if err := isGenesisFile(n.Genesis.Export, false); err != nil {
			errs = append(errs, err)
		}
	}

	if err := checkTimeInterval(n.Start, n.End, scenario.Duration); err != nil {
		errs = append(errs, err)
	}

	if err := n.isTypeValid(); err != nil {
		errs = append(errs, err)
	}

	if err := n.isTimerEventValid(); err != nil {
		errs = append(errs, err)
	}

	if err := n.isTimerSequenceValid(scenario); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// isTypeValid returns true if the node has valid type, false otherwise
func (n *Node) isTypeValid() error {
	return isTypeValid(n.Client.Type)
}

func isTypeValid(t string) error {
	switch t {
	case
		"validator",
		"rpc",
		"observer":
		return nil
	}
	return fmt.Errorf("type of node must be observer, rpc or validator, was set to %s", t)
}

// isTimerEventValid returns true if the timer event has valid type, false otherwise
func (n *Node) isTimerEventValid() error {
	errs := []error{}
	for _, event := range n.Timer {
		if err := isTimerEventValid(event); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func isTimerEventValid(e string) error {
	switch e {
	case
		"start",
		"end",
		"kill",
		"restart":
		return nil
	}
	return fmt.Errorf("timer event of node must be start, end, kill, or restart; was set to %s", e)
}

// isTimerSequenceValid returns true if the timer sequence make sense e.g. start only happens if the node is off, end only happens if the node is on, etc.
func (n *Node) isTimerSequenceValid(scenario *Scenario) error {
	var now bool = true // assumed to be started

	start := float32(0)
	if n.Start != nil {
		start = *n.Start
	}

	end := scenario.Duration
	if n.End != nil {
		end = *n.End
	}

	// sort timer sequence
	timings := make([]float32, 0, len(n.Timer))
	for t := range n.Timer {
		timings = append(timings, t)
	}
	slices.Sort(timings)

	for ix, t := range timings {
		if t < start {
			return fmt.Errorf("Node %s has event at time %f < start=%f", t, start)
		}

		if t > end {
			return fmt.Errorf("Node %s has event at time %f > end=%f", t, end)
		}

		next, err := isTimerSequenceValid(now, n.Timer[t])
		if err != nil {
			return fmt.Errorf("Node %s at time %f: %v", n.Name, t, err)
		}

		// if event is "kill", then it must be last
		if n.Timer[t] == "kill" && ix < len(timings) {
			return fmt.Errorf("Node %s has kill at time %f but there is more event queued.", n.Name, t)
		}

		now = next
	}

	return nil
}

func isTimerSequenceValid(on bool, event string) (bool, error) {
	switch event {
	case "start":
		if on {
			return false, fmt.Errorf("asked to start, already on")
		}
		return true, nil
	case "end":
		if !on {
			return false, fmt.Errorf("asked to end, already off")
		}
		return false, nil
	case "kill":
		if !on {
			return false, fmt.Errorf("asked to kill, already off")
		}
		return false, nil
	case "restart":
		if !on {
			return false, fmt.Errorf("asked to restart, already off")
		}
		return true, nil
	}
	return false, fmt.Errorf("event not recognized: %s", event)
}

// GetStaticValidatorCount returns the number of validator that begins at time 0
// and last the entire duration.
// Static Validator = Validator that lasts the entire duration of the run.
func (s *Scenario) GetStaticValidatorCount() int {
	var count int = 0
	for _, n := range s.Nodes {
		count += n.GetStaticValidatorCount(s)
	}
	return count
}

func (n *Node) GetStaticValidatorCount(scenario *Scenario) int {
	var count int = 0
	if n.Instances != nil {
		count = *n.Instances
	}

	if n.IsStaticValidator(scenario) {
		return count
	}

	return 0
}

// isGenesisFile checks if a file exist at a given path and that it is a ".g" extension
func isGenesisFile(path string, isImport bool) error {
	errs := []error{}
	_, err := os.Stat(path)

	if errors.Is(err, os.ErrNotExist) && isImport {
		errs = append(errs, fmt.Errorf("provided genesis file does not exist: %s", path))
	} else if err == nil && !isImport {
		errs = append(errs, fmt.Errorf("provided genesis file already exists: %s", path))
	}

	if ext := filepath.Ext(path); ext != ".g" {
		errs = append(errs, fmt.Errorf("provided path is not a genesis file: %s", path))
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the application configuration of a scenario.
func (a *Application) Check(scenario *Scenario) error {
	errs := []error{}

	if !namePattern.Match([]byte(a.Name)) {
		errs = append(errs, fmt.Errorf("application name must match %v, got %v", namePatternStr, a.Name))
	}

	if a.Type == "" {
		errs = append(errs, fmt.Errorf("application type must be specified"))
	} else if !app.IsSupportedApplicationType(a.Type) {
		errs = append(errs, fmt.Errorf("unknown application type: %v", a.Type))
	}

	if a.Instances != nil && *a.Instances < 0 {
		errs = append(errs, fmt.Errorf("number of instances must be >= 0, is %d", *a.Instances))
	}

	if a.Users != nil && *a.Users < 1 {
		errs = append(errs, fmt.Errorf("number of users must be >= 1, is %d", *a.Users))
	}

	if err := checkTimeInterval(a.Start, a.End, scenario.Duration); err != nil {
		errs = append(errs, err)
	}

	if err := a.Rate.Check(scenario); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the cheat configuration of a scenario.
func (c *Cheat) Check(scenario *Scenario) error {
	errs := []error{}

	if !namePattern.Match([]byte(c.Name)) {
		errs = append(errs, fmt.Errorf("cheat name must match %v, got %v", namePatternStr, c.Name))
	}

	if err := checkTimeInterval(c.Start, nil, scenario.Duration); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the traffic shape configuration of a source.
func (r *Rate) Check(scenario *Scenario) error {
	count := 0
	if r.Constant != nil {
		count++
	}
	if r.Slope != nil {
		count++
	}
	if r.Wave != nil {
		count++
	}
	if r.Auto != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("application must specify exactly one load shape, got %d", count)
	}

	if r.Constant != nil && *r.Constant < 0 {
		return fmt.Errorf("constant transaction rate must be >= 0, got %f", *r.Constant)
	}
	if r.Slope != nil {
		return r.Slope.Check()
	}
	if r.Wave != nil {
		return r.Wave.Check()
	}
	if r.Auto != nil {
		return r.Auto.Check()
	}
	return nil
}

// Check tests semantic constraints on the configuration of a slope traffic pattern.
func (s *Slope) Check() error {
	errs := []error{}

	if s.Start < 0 {
		errs = append(errs, fmt.Errorf("initial transaction rate must be >= 0, got %f", s.Start))
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the configuration of a wave-shaped traffic pattern.
func (w *Wave) Check() error {
	errs := []error{}

	min := float32(0.0)
	if w.Min != nil {
		min = *w.Min
	}
	max := w.Max

	if min < 0 {
		errs = append(errs, fmt.Errorf("minimum transaction rate must be >= 0, got %f", min))
	}
	if max < 0 {
		errs = append(errs, fmt.Errorf("maximum transaction rate must be >= 0, got %f", max))
	}
	if min > max {
		errs = append(errs, fmt.Errorf("minimum transaction rate must be <= maximum rate, got %f > %f", min, max))
	}

	if w.Period <= 0 {
		errs = append(errs, fmt.Errorf("wave priode must be > 0, got %f", w.Period))
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the configuration of a auto-shaped traffic pattern.
func (a *Auto) Check() error {
	errs := []error{}

	if a.Increase != nil {
		if *a.Increase <= 0 {
			errs = append(errs, fmt.Errorf("traffic rate increase per second must be positive, got %f", *a.Increase))
		}
	}
	if a.Decrease != nil {
		if *a.Decrease < 0 || *a.Decrease > 1 {
			errs = append(errs, fmt.Errorf("traffic decrease rate must be between 0 and 1, got %f", *a.Decrease))
		}
	}

	return errors.Join(errs...)
}

// checkTimeInterval is a utility function checking the validity of a start/end time pair.
func checkTimeInterval(start, end *float32, duration float32) error {
	realStart := float32(0.0)
	if start != nil {
		realStart = *start
	}
	realEnd := duration
	if end != nil {
		realEnd = *end
	}
	errs := []error{}
	if realStart < 0 {
		errs = append(errs, fmt.Errorf("start time must be >= 0, is %f", realStart))
	}
	if realStart > duration {
		errs = append(errs, fmt.Errorf("start time must be <= scenario duration (=%fs), is %f", duration, realStart))
	}
	if realEnd < realStart {
		errs = append(errs, fmt.Errorf("end time must be >= start time,  end=%fs, start=%fs", realEnd, realStart))
	} else {
		if realEnd > duration {
			errs = append(errs, fmt.Errorf("end time must be <= scenario duration, end=%fs, duration=%fs", realEnd, duration))
		}
	}
	return errors.Join(errs...)
}

// checkValidatorConstraints makes sure that there is correct number of validators.
// if during whole run there are just genesis validators, then one validator is enough
// if there is creation of new validators trough sfc, then at all times there has to be at least two validators validating at every time
// note during run validator won't be immediately registered for epoch, because it only happens during epoch seal,
// therefore this check is not sufficient on its own
func (s *Scenario) checkValidatorConstraints() error {

	// count static validators within the node
	gvCount := s.GetStaticValidatorCount()
	if s.NumValidators == nil {
		s.NumValidators = &gvCount
	}

	// This must be check here for test to pass
	// _must_ allow case where validator count = 0
	if *s.NumValidators < 0 {
		return fmt.Errorf("invalid number of validators: %d <= 0", *s.NumValidators)
	}

	// check if there are 2 genesis validators if there is dynamic validator
	var dynamicValidatorCount int = 0
	for _, node := range s.Nodes {
		if node.IsValidator() && !node.IsStaticValidator(s) {
			instances := 1
			if node.Instances != nil {
				instances = *node.Instances
			}
			dynamicValidatorCount += instances
		}
	}

	if dynamicValidatorCount > 0 && gvCount < 2 {
		return fmt.Errorf("Dynamic Validator count = %d; Number of static validators should have been at least 2: %d < 2", dynamicValidatorCount, *s.NumValidators)
	}

	// TODO add check for dynamic validators to have always at least two running at any time
	// needs to be implemented before enabling to shut down genesis validators
	return nil
}
