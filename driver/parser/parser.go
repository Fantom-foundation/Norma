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
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DefaultInstance = 1
	// MaxBlockGas, MaxEpochGas defaults borrowed from example-genesis.json
	DefaultMaxBlockGas = 20_500_000_000
	DefaultMaxEpochGas = 1_500_000_000_000
)

// Scenario is the root element of a scenario description. It defines basic
// scenario properties and lists a set of nodes and transaction source.
type Scenario struct {
	Name             string
	Duration         float32
	NumValidators    *int          `yaml:"num_validators,omitempty"` // nil == 1
	GenesisGasLimits GasLimits     `yaml:"genesis_gas_limit,omitempty"`
	Nodes            []Node        `yaml:",omitempty"`
	Applications     []Application `yaml:",omitempty"`
	Cheats           []Cheat       `yaml:",omitempty"`
}

// GasLimits is a configuration group for gas limit rules
type GasLimits struct {
	MaxBlockGas *uint64 `yaml:"max_block_gas,omitempty"`
	MaxEpochGas *uint64 `yaml:"max_epoch_gas,omitempty"`
}

// GetMaxBlockGas returns MaxBlockGas
func (s *Scenario) GetMaxBlockGas() uint64 {
	if s.GenesisGasLimits.MaxBlockGas != nil {
		return *s.GenesisGasLimits.MaxBlockGas
	}

	return DefaultMaxBlockGas
}

// GetMaxEpochGas returns MaxEpochGas
func (s *Scenario) GetMaxEpochGas() uint64 {
	if s.GenesisGasLimits.MaxEpochGas != nil {
		return *s.GenesisGasLimits.MaxEpochGas
	}

	return DefaultMaxEpochGas
}

func (s *Scenario) GetNumValidators() int {
	if s.NumValidators != nil {
		return *s.NumValidators
	}
	return 1
}

// Node is a configuration for a group of nodes with similar properties.
// Each node has a name, a set of features (e.g. 'validator', 'archve'),
// and a start and end time. Furthermore, nodes may be instantiated multiple
// times to create larger, homogenious groups easier.
type Node struct {
	Name      string
	Features  []string
	Instances *int       `yaml:",omitempty"` // nil is interpreted as 1
	Start     *float32   `yaml:",omitempty"` // nil is interpreted as 0
	End       *float32   `yaml:",omitempty"` // nil is interpreted as end-of-scenario
	Client    ClientType `yaml:",omitempty"`
	Mount     *string    `yaml:",omitempty"`
}

// IsValidator returns true if the node is defined as validator in Features
func (n *Node) IsValidator() bool {
	return n.Client.Type == "validator"
}

// IsStaticValidator returns true if the node is defined as validator in Features
func (n *Node) IsStaticValidator(s *Scenario) bool {
	start := float32(0)
	if n.Start != nil {
		start = *n.Start
	}

	end := s.Duration
	if n.End != nil {
		end = *n.End
	}

	return n.IsValidator() && start == float32(0) && end == s.Duration
}

// IsCheater returns true if the node is defined as cheater in Features
func (n *Node) IsCheater() bool {
	for _, item := range n.Features {
		if item == "cheater" {
			return true
		}
	}
	return false
}

// ClientType is an optional configuration for Node.
// ImageName can be used to specify the exact client version for the defined Node.
// Type can be used to configure the launching command of the client
type ClientType struct {
	ImageName string `yaml:",omitempty"` // nil is interpreted as main
	Type      string `yaml:",omitempty"` // nil is interpreted as observer
}

// Application is a load generator in the simulated network. Each application defines
// a type application load is generated for, a start and end time, a traffic
// shape (see Rate below), and a number of instances.
type Application struct {
	Name      string
	Type      string   `yaml:",omitempty"` // empty is interpreted as the default app type
	Instances *int     `yaml:",omitempty"` // nil is interpreted as 1
	Users     *int     `yaml:",omitempty"` // nil is interpreted as 1
	Start     *float32 `yaml:",omitempty"` // nil is interpreted as 0
	End       *float32 `yaml:",omitempty"` // nil is interpreted as end-of-scenario
	Rate      Rate
}

// Rate defines the shape of traffic to be generated. There are three types
// currently supported:
//   - constant ... traffic is created at a constant rate
//   - slope    ... traffic rate starts at 0 and is linearly increased
//   - wave     ... traffic rate follows a sin-wave pattern
//
// Only one of those options can be set for a single source.
type Rate struct {
	// Only one of the next fields may be set.
	Constant *float32 `yaml:",omitempty"`
	Slope    *Slope   `yaml:",omitempty"`
	Wave     *Wave    `yaml:",omitempty"`
	Auto     *Auto    `yaml:",omitempty"`
}

// Slope defines the parameters of a linearly increasing traffic pattern.
// The pattern is defined by a starting Tx/s rate and an increment per second.
type Slope struct {
	Start     float32 // starting Tx/s
	Increment float32 // increment by given Tx/s per second
}

// Wave defines the parameters of a sin-wave traffic pattern.
type Wave struct {
	Min    *float32 `yaml:",omitempty"` // Tx/s, nil = 0
	Max    float32  // Tx/s
	Period float32  // seconds
}

// A load pattern automatically maxing out throughput.
type Auto struct {
	Increase *float32 `yaml:",omitempty"` // increase in non-overload case per second in Tx/s, nil = 1
	Decrease *float32 `yaml:",omitempty"` // decrease in overload case in percent, nil = 0.2 (=20%)
}

// Cheat is a configuration to simulate cheating at a particular timing.
// For example, 2 validators with the same keys started at the same time can be considered
// an attempt to cheat.
type Cheat struct {
	Name  string
	Start *float32
}

// Parse parses a YAML based scenario description from the given reader.
// The parsing will fail if there are syntactic issues in the YAML file
// or if there are unknown keys. However, no semantic checks on the resulting
// scenariou will be conducted.
func Parse(reader io.Reader) (Scenario, error) {
	var res Scenario
	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)
	err := decoder.Decode(&res)
	return res, err
}

// ParseBytes parses the YAML encoded scenario in the given byte slice.
func ParseBytes(data []byte) (Scenario, error) {
	return Parse(bytes.NewReader(data))
}

// ParseFile parses the YAML encoded scenario in the given file.
func ParseFile(path string) (Scenario, error) {
	if reader, err := os.Open(path); err == nil {
		return Parse(reader)
	} else {
		return Scenario{}, err
	}
}
