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

	return errors.Join(errs...)
}

// GetGenesisValidatorCount returns the number of validator that begins at time 0
// and last the entire duration.
func (s *Scenario) GetGenesisValidatorCount() int {
	var count int = 0
	for _, n := range s.Nodes {
		count += n.GetGenesisValidatorCount(s)
	}
	return count
}

func (n *Node) GetGenesisValidatorCount(scenario *Scenario) int {
	if n.IsGenesisValidator(scenario) {
		return *n.Instances
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

// checkValidatorConstrains makes sure that there is correct number of validators.
// if during whole run there are just genesis validators, then one validator is enough
// if there is creation of new validators trough sfc, then at all times there has to be at least two validators validating at every time
// note during run validator won't be immediately registered for epoch, because it only happens during epoch seal,
// therefore this check is not sufficient on its own
func (s *Scenario) checkValidatorConstraints() error {
	
	// check genesis validators within the node
	gvCount := s.GetGenesisValidatorCount()
	if s.NumValidators == nil {
		s.NumValidators = &gvCount
	}

	// error if found more genesis validator than specified in NumValidators
	if *s.NumValidators < gvCount {
		return fmt.Errorf("mismatched number of genesis validators in scenario: NumValidator=%d < %d found in node list", *s.NumValidators, gvCount)
	}

	// NumValidators are used as short-hand to create gv nodes
	if *s.NumValidators > gvCount {
		gvCount = *s.NumValidators
	}
	
	// At least one genesis validator expected
	if gvCount <= 0 {
		return fmt.Errorf("invalid number of genesis validators in scenario: %d <= 0", gvCount)
	}


	// remove all GV from nodes. These will be initialized separately vs non-gv nodes.
	// NOTE: once we can specify more information about GV, this will need to change.
	s.removeGenesisValidator()

	// since all GVs are removed, the remaining validators are dynamic validators.
	// no dynamic validator = at least 1 = caught above
	dynamicValidatorCount := len(s.Nodes)
	if dynamicValidatorCount > 0 && gvCount < 2 {
		return fmt.Errorf("invalid number of genesis validators for sfc createValidator scenario: %d < 2", *s.NumValidators)
	}

	// TODO add check for dynamic validators to have always at least two running at any time
	// needs to be implemented before enabling to shut down genesis validators
	return nil
}

// removeGenesisValidator removes genesis validator from the list of nodes
func (s *Scenario) removeGenesisValidator() {
	s.Nodes = removeGenesisValidator(s, s.Nodes)
}

func removeGenesisValidator(scenario *Scenario, nodes []Node) []Node {
	var ret []Node
	for _, n := range nodes {
		if !n.IsGenesisValidator(scenario) {
			ret = append(ret, n)
		}
	}
	return ret
}
