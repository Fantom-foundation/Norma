package parser

import (
	"errors"
	"fmt"
	"strings"
)

// Check tests semantic constraints on the configuration of a scenario.
func (s *Scenario) Check() error {
	errs := []error{}
	if strings.TrimSpace(s.Name) == "" {
		errs = append(errs, fmt.Errorf("scenario name must not be empty"))
	}
	if s.Duration <= 0 {
		errs = append(errs, fmt.Errorf("scenario duration must be > 0"))
	}
	for _, node := range s.Nodes {
		if err := node.Check(s); err != nil {
			errs = append(errs, err)
		}
	}
	for _, source := range s.Sources {
		if err := source.Check(s); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// Check tests semantic constraints on the node configuration of a scenario.
func (n *Node) Check(scenario *Scenario) error {
	errs := []error{}
	if strings.TrimSpace(n.Name) == "" {
		errs = append(errs, fmt.Errorf("node name must not be empty"))
	}
	if n.Instances != nil && *n.Instances < 0 {
		errs = append(errs, fmt.Errorf("number of instances must be >= 0, is %d", *n.Instances))
	}

	if err := checkTimeInterval(n.Start, n.End, scenario.Duration); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

// Check tests semantic constraints on the source configuration of a scenario.
func (s *Source) Check(scenario *Scenario) error {
	errs := []error{}

	if strings.TrimSpace(s.Application) == "" {
		errs = append(errs, fmt.Errorf("sources must name an application"))
	}

	if s.Instances != nil && *s.Instances < 0 {
		errs = append(errs, fmt.Errorf("number of instances must be >= 0, is %d", *s.Instances))
	}

	if err := checkTimeInterval(s.Start, s.End, scenario.Duration); err != nil {
		errs = append(errs, err)
	}

	if err := s.Rate.Check(scenario); err != nil {
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
	if count != 1 {
		return fmt.Errorf("source must specify exactly one load shape, got %d", count)
	}

	if r.Constant != nil && *r.Constant < 0 {
		return fmt.Errorf("constant transaction rate must be >= 0, got %f", *r.Constant)
	}
	if r.Slope != nil && *r.Slope < 0 {
		return fmt.Errorf("slope transaction increase rate must be >= 0, got %f", *r.Slope)
	}
	if r.Wave != nil {
		return r.Wave.Check()
	}
	return nil
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
