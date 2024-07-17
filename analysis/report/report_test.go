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

package report

import (
	"os"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"go.uber.org/mock/gomock"
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
