package prometheusmon

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
)

type MockPrometheus struct{}

type MockPrometheusNode struct{}

func (m *MockPrometheus) Start(_ driver.Network, _ *docker.Network) (PrometheusNode, error) {
	return &MockPrometheusNode{}, nil
}

func (m *MockPrometheusNode) AddNode(_ driver.Node) error {
	return nil
}

func (m *MockPrometheusNode) Shutdown() error {
	return nil
}

func (m *MockPrometheusNode) GetUrl() string {
	return ""
}
