package nodemon

import (
	"fmt"
	"github.com/Fantom-foundation/Norma/driver"
	"github.com/Fantom-foundation/Norma/driver/monitoring"
	"github.com/golang/mock/gomock"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestTransactionsThroughputSource(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	source := NewTransactionsThroughputSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))

	now := time.Now()
	seconds := now.Unix()
	loops := 100
	nodes := []monitoring.Node{"A", "B", "C"}

	expected := make(map[monitoring.Node][]float32, len(nodes))
	for _, node := range nodes {
		timeGrow := rand.Intn(10) + 1
		expectedTxsList := make([]float32, 0, loops)
		// insert certain transactions in the same controlled delay between each
		for i := 0; i < loops; i++ {
			// progressively growing time
			timeStamp := time.Unix(seconds+int64(i*timeGrow), 0)
			txs := rand.Intn(1000)
			expectedTxs := float32(txs) / float32(int64(i*timeGrow)-int64((i-1)*timeGrow))
			expectedTxsList = append(expectedTxsList, expectedTxs)

			b := monitoring.Block{Height: i, Time: timeStamp, Txs: txs}
			source.OnBlock(node, b)
		}
		expected[node] = expectedTxsList
	}

	for node, txs := range expected {
		t.Run(fmt.Sprintf("node-%s", node), func(t *testing.T) {
			series, exists := source.GetData(node)
			if !exists {
				t.Errorf("data should exist")
			}

			// skip first block which is off
			for i := 1; i < loops; i++ {
				if got, want := series.GetRange(monitoring.BlockNumber(i), monitoring.BlockNumber(i+1))[0].Value, txs[i]; got != want {
					t.Errorf("transaction througput incorect: %3.2f != %3.2f", got, want)
				}
			}
		})
	}

}

func TestTransactionsTimeDiffBelowSec(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	writer := monitoring.NewMockWriterChain(ctrl)
	writer.EXPECT().Add(gomock.Any()).AnyTimes()

	source := NewTransactionsThroughputSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))

	seconds := time.Now().Unix()
	nsDiff := int64(50)
	secDif := 50 / 1e9

	// time diff only 50ns
	source.OnBlock("A", monitoring.Block{Height: 10, Time: time.Unix(seconds, 0), Txs: 10})
	source.OnBlock("A", monitoring.Block{Height: 11, Time: time.Unix(seconds, nsDiff), Txs: 10})

	series, exists := source.GetData("A")
	if !exists {
		t.Errorf("data should exist")
	}

	if got, want := series.GetLatest().Value, float32(10)/float32(secDif); got != want {
		t.Errorf("transaction througput incorect: %3.2f != %3.2f", got, want)
	}
}

func TestTransactionsCsvExport(t *testing.T) {
	ctrl := gomock.NewController(t)

	net := driver.NewMockNetwork(ctrl)
	net.EXPECT().RegisterListener(gomock.Any()).AnyTimes()
	net.EXPECT().GetActiveNodes().AnyTimes().Return([]driver.Node{})

	csvFile, _ := os.CreateTemp(t.TempDir(), "file.csv")
	writer := monitoring.NewWriterChain(csvFile)
	source := NewTransactionsThroughputSource(monitoring.NewMonitor(net, monitoring.MonitorConfig{}, writer))

	seconds := time.Now().Unix()

	// time diff only 50ns
	source.OnBlock("A", monitoring.Block{Height: 10, Time: time.Unix(seconds, 0), Txs: 10})
	source.OnBlock("A", monitoring.Block{Height: 11, Time: time.Unix(seconds+1, 0), Txs: 10})
	_ = writer.Close()
	content, _ := os.ReadFile(csvFile.Name())
	if got, want := string(content), "TransactionsThroughput, network, A, , 11, 10\n"; got != want {
		t.Errorf("unexpected export: %v != %v", got, want)
	}
}
