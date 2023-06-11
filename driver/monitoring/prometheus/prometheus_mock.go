package prometheusmon

import (
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/docker"
)

type MockPrometheusRunner struct{}

type MockPrometheus struct{}

func (m *MockPrometheusRunner) StartPrometheus(_ driver.Network, _ *docker.Network) (Prometheus, error) {
	return &MockPrometheus{}, nil
}

func (m *MockPrometheus) AddNode(_ driver.Node) error {
	return nil
}

func (m *MockPrometheus) Shutdown() error {
	return nil
}
