package parser

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

// Scenario is the root element of a scenario description. It defines basic
// scenario properties and lists a set of nodes and transaction source.
type Scenario struct {
	Name          string
	Duration      float32
	NumValidators *int          `yaml:"num_validators,omitempty"` // nil == 1
	Nodes         []Node        `yaml:",omitempty"`
	Applications  []Application `yaml:",omitempty"`
}

// Node is a configuration for a group of nodes with similar properties.
// Each node has a name, a set of features (e.g. 'validator', 'archve'),
// and a start and end time. Furthermore, nodes may be instantiated multiple
// times to create larger, homogenious groups easier.
type Node struct {
	Name      string
	Features  []string
	Instances *int     `yaml:",omitempty"` // nil is interpreted as 1
	Start     *float32 `yaml:",omitempty"` // nil is interpreted as 0
	End       *float32 `yaml:",omitempty"` // nil is interpreted as end-of-scenario
}

// Application is a load generator in the simulated network. Each application defines
// a type application load is generated for, a start and end time, a traffic
// shape (see Rate below), and a number of instances.
type Application struct {
	Name      string
	Instances *int     `yaml:",omitempty"` // nil is interpreted as 1
	Accounts  *int     `yaml:",omitempty"` // nil is interpreted as 1
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
	Slope    *float32 `yaml:",omitempty"`
	Wave     *Wave    `yaml:",omitempty"`
}

// Wave defines the parameters of a sin-wave traffic pattern.
type Wave struct {
	Min    *float32 `yaml:",omitempty"` // Tx/s, nil = 0
	Max    float32  // Tx/s
	Period float32  // seconds
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
