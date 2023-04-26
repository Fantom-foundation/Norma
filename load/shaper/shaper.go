package shaper

import (
	"time"
)

// Shaper defines delays between produced txs to ensure desired produced traffic profile.
type Shaper interface {
	// GetNextWaitTime provides the time to wait before the next tx should be sent
	GetNextWaitTime() time.Duration
}
