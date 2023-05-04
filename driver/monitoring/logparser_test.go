package monitoring

import (
	"fmt"
	"os"
	"testing"
)

func TestImplements(t *testing.T) {
	var inst LogParserThroughput
	var _ Throughput = &inst
}

func TestParseTime(t *testing.T) {
	tt, err := parseTime("[05-04|09:34:16.512]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	if val := tt.Month(); val != 5 {
		t.Errorf("wrong time parsed: %d", val)
	}
	if val := tt.Day(); val != 4 {
		t.Errorf("wrong time parsed: %d", val)
	}
	if val := tt.Hour(); val != 9 {
		t.Errorf("wrong time parsed: %d", val)
	}
	if val := tt.Minute(); val != 34 {
		t.Errorf("wrong time parsed: %d", val)
	}
	if val := tt.Second(); val != 16 {
		t.Errorf("wrong time parsed: %d", val)
	}
	if val := tt.Nanosecond(); val != 512*1e6 {
		t.Errorf("wrong time parsed: %d", val)
	}

}

func TestGetCorrectLogData(t *testing.T) {
	parser := createLogParser(t)

	if val, err := parser.GetTransactions(3); val != 1 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := parser.GetTransactions(6); val != 4 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := parser.GetGas(3); val != 117_867 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := parser.GetGas(5); val != 138_470 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	time1, err := parseTime("[05-04|09:34:15.537]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	if val, err := parser.GetBlockTime(3); val != time1 || err != nil {
		t.Errorf("wrong value: %s, err: %s", val, err)
	}

	time2, err := parseTime("[05-04|09:34:16.512]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	if val, err := parser.GetBlockTime(5); val != time2 || err != nil {
		t.Errorf("wrong value: %s, err: %s", val, err)
	}

	time3, err := parseTime("[05-04|09:34:17.003]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	diff := time3.Sub(time2)
	if val, err := parser.GetBlockDelay(6); val != diff || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}
}

func TestGetNonExistingBlocks(t *testing.T) {
	parser := createLogParser(t)

	if _, err := parser.GetTransactions(1); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := parser.GetGas(10); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := parser.GetBlockTime(0); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := parser.GetBlockDelay(7); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	// non-existing previous block
	if _, err := parser.GetBlockDelay(2); err != ErrNotFound {
		t.Errorf("block should not exist")
	}
}

func TestLogAppended(t *testing.T) {
	parser := createLogParser(t)

	// block does not exist here
	if _, err := parser.GetTransactions(1); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	// new log line has arrived here
	f, err := os.OpenFile(parser.file, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("cannot open file: %s", err)
	}

	if _, err := f.WriteString("INFO [05-04|09:35:17.003] New block   index=7 id=3:10:d7da0b      gas_used=111,223 txs=666/0 age=310.575ms t=5.249ms \n"); err != nil {
		t.Fatalf("cannot write file: %s", err)
	}

	_ = f.Close()

	if val, err := parser.GetTransactions(7); val != 666 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}
}

func createLogParser(t *testing.T) *LogParserThroughput {
	dir := t.TempDir()

	testLog :=
		"INFO [05-04|09:34:15.080] New block                                index=2 id=2:1:247c79       gas_used=417,928 txs=2/0 age=7.392s t=3.686ms \n" +
			"INFO [05-04|09:34:15.537] New block                                index=3 id=3:1:3d6fb6       gas_used=117,867 txs=1/0 age=343.255ms t=1.579ms \n" +
			"INFO [05-04|09:34:16.027] New block                                index=4 id=3:4:9bb789       gas_used=43426   txs=1/0 age=380.470ms t=1.540ms \n" +
			"INFO [05-04|09:34:16.512] New block                                index=5 id=3:7:a780ce       gas_used=138,470 txs=5/0 age=374.251ms t=3.796ms \n" +
			"INFO [05-04|09:34:17.003] New block                                index=6 id=3:10:d7da0b      gas_used=105,304 txs=4/0 age=381.575ms t=3.249ms \n"

	file, err := os.Create(fmt.Sprintf("%s/file.txt", dir))
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		t.Fatalf("error writing to file: %v", err)
	}

	data := []byte(testLog)
	if _, err = file.Write(data); err != nil {
		t.Fatalf("error writing to file: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("error closing file: %v", err)
	}

	return NewLogParserThroughput(file.Name())
}
