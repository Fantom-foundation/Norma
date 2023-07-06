package shaper

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver/parser"
)

// Shaper defines the shape of traffic to be produced by an application.
type Shaper interface {
	// GetNumMessagesInInterval provides the number of messages to be produced
	// in the given time interval. The result is expected to be >= 0.
	GetNumMessagesInInterval(start time.Time, duration time.Duration) float64
}

// ParseRate parses rate from the parser.
func ParseRate(rate *parser.Rate) (Shaper, error) {
	// return default constant shaper if rate is not specified
	if rate == nil {
		return NewConstantShaper(0), nil
	}

	if rate.Constant != nil {
		return NewConstantShaper(float64(*rate.Constant)), nil
	}
	if rate.Slope != nil {
		return NewSlopeShaper(float64(rate.Slope.Start), float64(rate.Slope.Increment)), nil
	}

	// TODO: add wave shaper

	return nil, fmt.Errorf("unknown rate type")
}
