package report

import (
	"os"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
)

func TestRenderReport(t *testing.T) {
	reports := []Report{
		SingleEvalReport,
	}

	for _, report := range reports {
		t.Run(report.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			net := driver.NewMockNetwork(ctrl)
			net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
			net.EXPECT().GetActiveNodes().AnyTimes().Return(nil)

			// Create a dummy measurement data file.
			outputdir := t.TempDir()
			monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{
				OutputDir: outputdir,
			})
			if err != nil {
				t.Fatalf("failed to start monitor: %v", err)
			}

			if err := monitor.Shutdown(); err != nil {
				t.Fatalf("failed to shut down monitor: %v", err)
			}

			data := monitor.GetMeasurementFileName()
			if res, err := report.Render(data, outputdir); err != nil {
				t.Errorf("failed to render report: %v", err)
			} else if _, err := os.Stat(outputdir + "/" + res); err != nil {
				t.Errorf("generated report file %s does not exist: %v", res, err)
			}
		})
	}
}
