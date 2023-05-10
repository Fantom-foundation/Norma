package monitoring

import (
	"testing"
)

func TestProcessesLog(t *testing.T) {
	mon := createBlockMetrics()

	if blocks, _ := mon.GetBlockHeight(); blocks != 6 {
		t.Errorf("wrong number of blocks captured: %d", blocks)
	}
}

func TestGetCorrectLogData(t *testing.T) {
	mon := createBlockMetrics()

	if val, err := mon.GetNumberOfTransactions(3); val != 1 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := mon.GetNumberOfTransactions(6); val != 4 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := mon.GetGas(3); val != 117_867 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	if val, err := mon.GetGas(5); val != 138_470 || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}

	time1, err := parseTime("[05-04|09:34:15.537]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	if val, err := mon.GetBlockTime(3); val != time1 || err != nil {
		t.Errorf("wrong value: %s, err: %s", val, err)
	}

	time2, err := parseTime("[05-04|09:34:16.512]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	if val, err := mon.GetBlockTime(5); val != time2 || err != nil {
		t.Errorf("wrong value: %s, err: %s", val, err)
	}

	time3, err := parseTime("[05-04|09:34:17.003]")
	if err != nil {
		t.Errorf("cannot parse: %s", err)
	}

	diff := time3.Sub(time2)
	if val, err := mon.GetBlockDelay(6); val != diff || err != nil {
		t.Errorf("wrong value: %d, err: %s", val, err)
	}
}

func TestGetNonExistingBlocks(t *testing.T) {
	mon := createBlockMetrics()

	if _, err := mon.GetNumberOfTransactions(1); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := mon.GetGas(10); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := mon.GetBlockTime(0); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	if _, err := mon.GetBlockDelay(7); err != ErrNotFound {
		t.Errorf("block should not exist")
	}

	// non-existing previous block
	if _, err := mon.GetBlockDelay(2); err != ErrNotFound {
		t.Errorf("block should not exist")
	}
}

func createBlockMetrics() *NodeMetrics {
	blockReader := NewLogReader(createTestLog())

	mon := CreateNodeMetrics(blockReader)
	mon.drain()

	return mon
}
