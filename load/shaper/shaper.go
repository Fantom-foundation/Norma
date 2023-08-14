package shaper

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver/parser"
)

//go:generate mockgen -source shaper.go -destination shaper_mock.go -package shaper

// Shaper defines the shape of traffic to be produced by an application.
type Shaper interface {
	// Start notifies the shaper that processing is started at the given time
	// and provides a source for fetching load information.
	Start(time.Time, LoadInfoSource)

	// GetNumMessagesInInterval provides the number of messages to be produced
	// in the given time interval. The result is expected to be >= 0.
	GetNumMessagesInInterval(start time.Time, duration time.Duration) float64
}

// LoadInfoSource defines an interface for load-sensitive traffic shapes to
// request load state information.
type LoadInfoSource interface {
	GetSentTransactions() (uint64, error)
	GetReceivedTransactions() (uint64, error)
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
	if rate.Auto != nil {
		increase := 1.0
		if rate.Auto.Increase != nil {
			increase = float64(*rate.Auto.Increase)
		}
		decrease := 0.2
		if rate.Auto.Decrease != nil {
			decrease = float64(*rate.Auto.Decrease)
		}
		return NewAutoShaper(increase, decrease), nil
	}
	if rate.Wave != nil {
		min := float32(0)
		if rate.Wave.Min != nil {
			min = *rate.Wave.Min
		}
		return NewWaveShaper(min, rate.Wave.Max, rate.Wave.Period), nil
	}

	return nil, fmt.Errorf("unknown rate type")
}
