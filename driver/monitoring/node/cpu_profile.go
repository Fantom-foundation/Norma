package nodemon

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver/monitoring/export"
	"io"
	"net/http"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	opera "github.com/Fantom-foundation/Norma/driver/node"
)

type PprofData []byte

func GetPprofData(node driver.Node, duration time.Duration) (PprofData, error) {
	url := node.GetHttpServiceUrl(&opera.OperaPprofService)
	if url == nil {
		return nil, fmt.Errorf("node does not offer the pprof service")
	}
	resp, err := http.Get(fmt.Sprintf("%s/debug/pprof/profile?seconds=%d", *url, int(duration.Seconds())))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch result: %v", resp)
	}
	return io.ReadAll(resp.Body)
}

// NodeCpuProfile periodically collects CPU profiles from individual nodes.
var NodeCpuProfile = mon.Metric[mon.Node, mon.TimeSeries[PprofData]]{
	Name:        "NodeCpuProfile",
	Description: "CpuProfile samples of a node at various times.",
}

func init() {
	if err := mon.RegisterSource(NodeCpuProfile, NewNodeCpuProfileSource); err != nil {
		panic(fmt.Sprintf("failed to register metric source: %v", err))
	}
}

// NewNodeCpuProfileSource creates a new data source periodically collecting
// CPU profiling data at configured sampling periods.
func NewNodeCpuProfileSource(monitor *mon.Monitor) mon.Source[mon.Node, mon.TimeSeries[PprofData]] {
	return newPeriodicNodeDataSource[PprofData](
		NodeCpuProfile,
		monitor,
		10*time.Second, // Sampling period; TODO: make customizable
		&cpuProfileSensorFactory{},
		export.DirectConverter[PprofData]{},
	)
}

type cpuProfileSensorFactory struct{}

func (f *cpuProfileSensorFactory) CreateSensor(node driver.Node) (Sensor[PprofData], error) {
	return &cpuProfileSensor{
		node,
		5 * time.Second, // the duration of the CPU profile collection; TODO: make configurable
	}, nil
}

type cpuProfileSensor struct {
	node     driver.Node
	duration time.Duration
}

func (s *cpuProfileSensor) ReadValue() (PprofData, error) {
	return GetPprofData(s.node, s.duration)
}
