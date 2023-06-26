package shaper

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver/parser"
)

// Shaper defines delays between produced txs to ensure desired produced traffic profile.
type Shaper interface {
	// GetNextWaitTime provides the time to wait before the next tx should be sent
	GetNextWaitTime() time.Duration
}

// ParseRate parses rate from the parser.
func ParseRate(rate *parser.Rate) (Shaper, error) {
	// return default constant shaper if rate is not specified
	if rate == nil {
		return NewConstantShaper(0), nil
	}

	if rate.Constant != nil {
		return NewConstantShaper(*rate.Constant), nil
	}
	if rate.Slope != nil {
		return NewSlopeShaper(rate.Slope.Start, rate.Slope.Increment), nil
	}

	// TODO: add wave shaper

	return nil, fmt.Errorf("unknown rate type")
}
