package monitoring

import (
	"strings"
	"testing"
)

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

func TestParseBlock(t *testing.T) {
	blockReader := NewLogReader(strings.NewReader(Node1TestLog))

	var count int
	for got := range blockReader {
		want, exists := BlockHeight1TestMap[got.Height]
		if !exists {
			t.Errorf("unknow block: %d", got.Height)
		}

		if want.Txs != got.Txs {
			t.Errorf("values do not match: %v != %v", want.Txs, got.Txs)
		}
		if want.GasUsed != got.GasUsed {
			t.Errorf("values do not match: %v != %v", want.GasUsed, got.GasUsed)
		}
		if want.Time != got.Time {
			t.Errorf("values do not match: %v != %v", want.Time, got.Time)
		}
		if want.ProcessingTime != got.ProcessingTime {
			t.Errorf("values do not match: %s != %s", want.Time, got.Time)
		}

		count++
	}

	if len(BlockHeight1TestMap) != count {
		t.Errorf("not all keys were visited")
	}
}

func TestParseLogStream(t *testing.T) {
	blockReader := NewLogReader(strings.NewReader(Node1TestLog))

	var count int
	for got := range blockReader {
		want, err := BlockHeight1TestMap[got.Height]
		if !err {
			t.Errorf("unknow block: %d", got.Height)
		}

		if want.Txs != got.Txs {
			t.Errorf("values do not match: %v != %v", want.Txs, got.Txs)
		}
		if want.GasUsed != got.GasUsed {
			t.Errorf("values do not match: %v != %v", want.GasUsed, got.GasUsed)
		}
		if want.Time != got.Time {
			t.Errorf("values do not match: %v != %v", want.Time, got.Time)
		}
		if want.ProcessingTime != got.ProcessingTime {
			t.Errorf("values do not match: %s != %s", want.Time, got.Time)
		}

		count++
	}

	if len(BlockHeight1TestMap) != count {
		t.Errorf("not all keys were visited")
	}
}
