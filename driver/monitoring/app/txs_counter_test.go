package app

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/Fantom-foundation/Norma/load/generator"
	"github.com/golang/mock/gomock"
	"os"
	"testing"
)

func TestApplicationRegistered(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	// generate test mock applications
	size := 1000
	appsCount := 11
	apps := make(map[monitoring.App][]driver.Application, appsCount)
	for i := 0; i < size; i++ {
		app := driver.NewMockApplication(ctrl)
		appName := fmt.Sprintf("app-%d", i%appsCount)
		app.EXPECT().Config().AnyTimes().Return(&driver.ApplicationConfig{
			Name:     appName,
			Rate:     0,
			Accounts: i + 1,
		})
		txsCount := generator.NewMockTransactionCounts(ctrl)
		txsCount.EXPECT().GetAmountOfSentTxs().AnyTimes().Return(uint64(i * 10))
		txsCount.EXPECT().GetAmountOfReceivedTxs().AnyTimes().Return(uint64(i*20), nil)

		app.EXPECT().GetTransactionCounts().AnyTimes().Return(txsCount, true)

		arr, exists := apps[monitoring.App(appName)]
		if !exists {
			arr = make([]driver.Application, 0, size/appsCount+1)
		}

		arr = append(arr, app)
		apps[monitoring.App(appName)] = arr
	}

	// simulate applications received
	source1 := NewSentTransactionsSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))
	source2 := NewReceivedTransactionsSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))

	for _, source := range []*TxsCounter{source1, source2} {
		t.Run(fmt.Sprintf("%s", source.metric.Name), func(t *testing.T) {

			// fill-in data
			for name := range apps {
				for _, app := range apps[name] {
					source.AfterApplicationCreation(app)
				}
			}

			// shutdown causes calculation of data
			if err := source.Shutdown(); err != nil {
				t.Fatalf("cannot shutdown: %s", err)
			}

			// verify results
			for _, subject := range source.GetSubjects() {
				_, exists := apps[subject]
				if !exists {
					t.Errorf("subject does not exist within expected subjects: %v", subject)
				}
			}

			if got, want := len(source.GetSubjects()), len(apps); got != want {
				t.Errorf("amount of subjects do not mathc: %d != %d", got, want)
			}

			for subject := range apps {
				series, exists := source.GetData(subject)
				if !exists {
					t.Errorf("data for subject dos not exist: %v", subject)
				}

				for i, point := range series.GetRange(0, 100000) {
					txsCount, _ := apps[subject][i].GetTransactionCounts()
					want, _ := source.getter(txsCount)
					if point.Value != want {
						t.Errorf("data series contain unexpected value: %v != %v", point.Value, want)
					}
					if got, want := point.Position, apps[subject][i].Config().Accounts; got != want {
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

	app := driver.NewMockApplication(ctrl)
	app.EXPECT().Config().AnyTimes().Return(&driver.ApplicationConfig{
		Name:     fmt.Sprintf("app-%d", 666),
		Rate:     0,
		Accounts: 999,
	})
	txsCount := generator.NewMockTransactionCounts(ctrl)
	txsCount.EXPECT().GetAmountOfSentTxs().AnyTimes().Return(uint64(15))
	txsCount.EXPECT().GetAmountOfReceivedTxs().AnyTimes().Return(uint64(16), nil)

	app.EXPECT().GetTransactionCounts().AnyTimes().Return(txsCount, true)

	csvFile1, _ := os.CreateTemp(t.TempDir(), "file.csv")
	writer1 := monitoring.NewWriterChain(csvFile1)
	source1 := NewSentTransactionsSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer1))

	csvFile2, _ := os.CreateTemp(t.TempDir(), "file.csv")
	writer2 := monitoring.NewWriterChain(csvFile2)
	source2 := NewReceivedTransactionsSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer2))

	csvFiles := []*os.File{csvFile1, csvFile2}
	expected := []string{
		"SentTransactions, network, , app-666, , , 999, 15\n",
		"ReceivedTransactions, network, , app-666, , , 999, 16\n",
	}

	for i, source := range []*TxsCounter{source1, source2} {
		t.Run(fmt.Sprintf("%s", source.metric.Name), func(t *testing.T) {
			// insert data
			source.AfterApplicationCreation(app)

			// shutdown causes calculation of data
			if err := source.Shutdown(); err != nil {
				t.Fatalf("cannot shutdown: %s", err)
			}
			_ = source.monitor.Writer().Close()

			content, _ := os.ReadFile(csvFiles[i].Name())
			if got, want := string(content), expected[i]; got != want {
				t.Errorf("unexpected export: %v != %v", got, want)
			}
		})
	}

}
