package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/load/app"
	"github.com/golang/mock/gomock"
)

func TestApplicationRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	// generate test mock applications
	size := 1000
	appsCount := 11
	apps := make(map[monitoring.App][]driver.Application, appsCount)
	appsList := make([]driver.Application, 0)
	for i := 0; i < size; i++ {
		application := driver.NewMockApplication(ctrl)
		appName := fmt.Sprintf("app-%d", i%appsCount)
		application.EXPECT().Config().AnyTimes().Return(&driver.ApplicationConfig{
			Name:     appName,
			Accounts: i + 1,
		})
		application.EXPECT().GetTransactionCounts().AnyTimes().Return(app.TransactionCounts{
			SentTxs:     uint64(i * 10),
			ReceivedTxs: uint64(i * 20),
		}, nil)
		appsList = append(appsList, application)

		arr, exists := apps[monitoring.App(appName)]
		if !exists {
			arr = make([]driver.Application, 0, size/appsCount+1)
		}

		arr = append(arr, application)
		apps[monitoring.App(appName)] = arr
	}
	net.EXPECT().GetActiveApplications().AnyTimes().Return(appsList)

	// simulate applications received
	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: t.TempDir()})
	if err != nil {
		t.Fatalf("failed to start monitor instance: %v", err)
	}
	err = errors.Join(
		monitoring.InstallSourceFor(SentTransactions, monitor),
		monitoring.InstallSourceFor(ReceivedTransactions, monitor),
	)
	if err != nil {
		t.Fatalf("failed to install metric sources: %v", err)
	}

	// shutdown causes calculation of data
	if err := monitor.Shutdown(); err != nil {
		t.Fatalf("cannot shutdown: %s", err)
	}

	metrics := []monitoring.Metric[monitoring.App, monitoring.Series[int, int]]{
		SentTransactions, ReceivedTransactions,
	}

	for _, metric := range metrics {
		t.Run(metric.Name, func(t *testing.T) {

			// verify results
			subjects := monitoring.GetSubjects(monitor, metric)
			for _, subject := range subjects {
				_, exists := apps[subject]
				if !exists {
					t.Errorf("subject does not exist within expected subjects: %v", subject)
				}
			}

			if got, want := len(subjects), len(apps); got != want {
				t.Errorf("amount of subjects do not match: %d != %d", got, want)
			}

			for subject, app := range apps {
				series, exists := monitoring.GetData(monitor, subject, metric)
				if !exists {
					t.Errorf("data for subject dos not exist: %v", subject)
					continue
				}

				for i, point := range series.GetRange(0, 100000) {
					txsCount, err := app[i].GetTransactionCounts()
					if err != nil {
						t.Fatalf("failed to get txs counts; %v", err)
					}
					want := int(txsCount.SentTxs)
					if metric.Name == ReceivedTransactions.Name {
						want = int(txsCount.ReceivedTxs)
					}
					if point.Value != want {
						t.Errorf("data series contain unexpected value: %v != %v", point.Value, want)
					}
					if got, want := point.Position, app[i].Config().Accounts; got != want {
						t.Errorf("positions do not match: %v != %v", got, want)
					}
				}
			}
		})
	}
}

func TestApplicationPrinted(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	application := driver.NewMockApplication(ctrl)
	application.EXPECT().Config().AnyTimes().Return(&driver.ApplicationConfig{
		Name:     fmt.Sprintf("app-%d", 666),
		Accounts: 999,
	})
	application.EXPECT().GetTransactionCounts().AnyTimes().Return(app.TransactionCounts{
		SentTxs:     uint64(15),
		ReceivedTxs: uint64(16),
	}, nil)
	net.EXPECT().GetActiveApplications().AnyTimes().Return([]driver.Application{application})

	outDir := t.TempDir()
	monitor, err := monitoring.NewMonitor(net, monitoring.MonitorConfig{OutputDir: outDir})
	if err != nil {
		t.Fatalf("failed to start monitoring instance for test: %v", err)
	}

	err = errors.Join(
		monitoring.InstallSourceFor(SentTransactions, monitor),
		monitoring.InstallSourceFor(ReceivedTransactions, monitor),
	)
	if err != nil {
		t.Fatalf("failed to install metric sources: %v", err)
	}

	if err := monitor.Shutdown(); err != nil {
		t.Fatalf("failed to shut down monitoring: %v", err)
	}

	content, _ := os.ReadFile(monitor.GetMeasurementFileName())

	expected := []string{
		"SentTransactions, network, , app-666, , , 999, 15\n",
		"ReceivedTransactions, network, , app-666, , , 999, 16\n",
	}
	for _, line := range expected {
		if got, want := string(content), line; !strings.Contains(got, want) {
			t.Errorf("unexpected export: %v != %v", got, want)
		}
	}
}
