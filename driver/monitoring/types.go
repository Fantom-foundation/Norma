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

// App is an identifier of an application deployed in the network
type App string

// User is an identifier of a user/account interacting with an application.
type User struct {
	App App // The application being a part of.
	Id  int // A unique identifier of the user/account.
}

func (a *User) Less(b *User) bool {
	return a.App < b.App || (a.App == b.App && a.Id < b.Id)
}

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

// BlockStatus encapsulates epoch, blockheight
type BlockStatus struct {
	Epoch       int
	BlockHeight int
}

func (b BlockStatus) String() string {
	return fmt.Sprintf("%d/%d", b.Epoch, b.BlockHeight)
}

// Percent is used to represent a percentage of some value. Internaly it is
// represented as a float value, typically in the range between [0,1] denoting
// values between 0% and 100%. However, values exceeding those boundaries are
// allowed (e.g. 120% or -30%).
type Percent float32

func (p Percent) String() string {
	return fmt.Sprintf("%.1f%%", p*100)
}
