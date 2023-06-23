package nodemon

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	mon "github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/driver/monitoring/utils"
	opera "github.com/Fantom-foundation/Norma/driver/node"
)

type PprofData []byte

func GetPprofData(node driver.Node, duration time.Duration) (PprofData, error) {
	url := node.GetServiceUrl(&opera.OperaDebugService)
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
var NodeCpuProfile = mon.Metric[mon.Node, mon.Series[mon.Time, string]]{
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
func NewNodeCpuProfileSource(monitor *monitoring.Monitor) mon.Source[mon.Node, mon.Series[mon.Time, string]] {
	return newPeriodicNodeDataSource[string](
		NodeCpuProfile,
		monitor,
		10*time.Second, // Sampling period; TODO: make customizable
		&cpuProfileSensorFactory{
			outputDir: monitor.Config().OutputDir,
		},
	)
}

type cpuProfileSensorFactory struct {
	outputDir string
}

func (f *cpuProfileSensorFactory) CreateSensor(node driver.Node) (utils.Sensor[string], error) {
	return &cpuProfileSensor{
		node:      node,
		duration:  5 * time.Second, // the duration of the CPU profile collection; TODO: make configurable
		outputDir: f.outputDir,
	}, nil
}

type cpuProfileSensor struct {
	node        driver.Node
	duration    time.Duration
	outputDir   string
	numProfiles int
}

func (s *cpuProfileSensor) ReadValue() (string, error) {
	data, err := GetPprofData(s.node, s.duration)
	if err != nil {
		return "", err
	}
	dir := fmt.Sprintf("%s/cpu_profiles/%s", s.outputDir, s.node.GetLabel())
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s/%06d.prof", dir, s.numProfiles)
	s.numProfiles++
	if err := os.WriteFile(filename, data, 0600); err != nil {
		return "", err
	}
	return filename, nil
}
