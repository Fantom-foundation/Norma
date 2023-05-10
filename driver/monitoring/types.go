package monitoring

import (
	"fmt"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
)

// Node identifies a node in the network.
type Node driver.NodeID

// Network is a unit type to reference the full managed network in a scenario
// as the subject of a metric.
type Network struct{}

// Time is the time used in time series. The value represents UnixNanos.
// Note, time.Time cannot be used since it doesn't satisfy constraints.Ordered.
type Time uint64

func NewTime(time time.Time) Time {
	return Time(time.UnixNano())
}

func (t Time) Time() time.Time {
	return time.Unix(int64(t/1e9), int64(t%1e9))
}

func (t Time) String() string {
	return t.Time().String()
}

// BlockNumber is the type used to identify a block.
type BlockNumber int

// Percent is used to represent a percentage of some value. Internaly it is
// represented as a float value, typically in the range between [0,1] denoting
// values between 0% and 100%. However, values exceeding those boundaries are
// allowed (e.g. 120% or -30%).
type Percent float32

func (p Percent) String() string {
	return fmt.Sprintf("%.1f%%", p*100)
}
