package monitoring

import (
	"io"
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
	blockReader := NewLogReader(createTestLog())

	time2, _ := parseTime("[05-04|09:34:15.080]")
	time3, _ := parseTime("[05-04|09:34:15.537]")
	time4, _ := parseTime("[05-04|09:34:16.027]")
	time5, _ := parseTime("[05-04|09:34:16.512]")
	time6, _ := parseTime("[05-04|09:34:17.003]")

	s1 := "INFO [05-04|09:34:15.080] New block                                index=2 id=2:1:247c79       gas_used=417,928 txs=2/0 age=7.392s t=3.686ms \n"
	s2 := "INFO [05-04|09:34:15.537] New block                                index=3 id=3:1:3d6fb6       gas_used=117,867 txs=1/0 age=343.255ms t=1.579ms \n"
	s3 := "INFO [05-04|09:34:16.027] New block                                index=4 id=3:4:9bb789       gas_used=43426   txs=1/0 age=380.470ms t=1.540ms \n"
	s4 := "INFO [05-04|09:34:16.512] New block                                index=5 id=3:7:a780ce       gas_used=138,470 txs=5/0 age=374.251ms t=3.796ms \n"
	s5 := "INFO [05-04|09:34:17.003] New block                                index=6 id=3:10:d7da0b      gas_used=105,304 txs=4/0 age=381.575ms t=3.249ms \n"

	logToBlock := map[string]Block{
		s1: {2, time2, 2, 417_928},
		s2: {3, time3, 1, 117_867},
		s3: {4, time4, 1, 43426},
		s4: {5, time5, 5, 138_470},
		s5: {6, time6, 4, 105_304},
	}

	blockToLog := make(map[int]string, len(logToBlock))
	for k, v := range logToBlock {
		blockToLog[v.height] = k
	}

	for b := range blockReader {
		s, exists := blockToLog[b.height]
		if !exists {
			t.Errorf("unknow block: %d", b.height)
		}

		val, exists := logToBlock[s]
		if !exists {
			t.Errorf("unknow log: %s", s)
		}

		if val.txs != b.txs {
			t.Errorf("values do not match: %v != %v", val.txs, b.txs)
		}
		if val.gasUsed != b.gasUsed {
			t.Errorf("values do not match: %v != %v", val.gasUsed, b.gasUsed)
		}
		if val.time != b.time {
			t.Errorf("values do not match: %v != %v", val.time, b.time)
		}

		delete(blockToLog, b.height)
	}

	if len(blockToLog) != 0 {
		t.Errorf("not all keys were visited")
	}
}

func TestParseLogStream(t *testing.T) {
	blockReader := NewLogReader(createTestLog())

	time2, _ := parseTime("[05-04|09:34:15.080]")
	time3, _ := parseTime("[05-04|09:34:15.537]")
	time4, _ := parseTime("[05-04|09:34:16.027]")
	time5, _ := parseTime("[05-04|09:34:16.512]")
	time6, _ := parseTime("[05-04|09:34:17.003]")

	expected := map[int]Block{
		2: {2, time2, 2, 417_928},
		3: {3, time3, 1, 117_867},
		4: {4, time4, 1, 43426},
		5: {5, time5, 5, 138_470},
		6: {6, time6, 4, 105_304},
	}

	for b := range blockReader {
		val, exists := expected[b.height]
		if !exists {
			t.Errorf("unknow block: %d", b.height)
		}

		if val.txs != b.txs {
			t.Errorf("values do not match: %v != %v", val.txs, b.txs)
		}
		if val.gasUsed != b.gasUsed {
			t.Errorf("values do not match: %v != %v", val.gasUsed, b.gasUsed)
		}
		if val.time != b.time {
			t.Errorf("values do not match: %v != %v", val.time, b.time)
		}

		delete(expected, b.height)
	}

	if len(expected) != 0 {
		t.Errorf("not all keys were visited")
	}
}

func createTestLog() io.Reader {
	testLog :=
		"INFO [05-04|09:34:15.080] New block                                index=2 id=2:1:247c79       gas_used=417,928 txs=2/0 age=7.392s t=3.686ms \n" +
			"INFO [05-04|09:34:15.537] New block                                index=3 id=3:1:3d6fb6       gas_used=117,867 txs=1/0 age=343.255ms t=1.579ms \n" +
			"INFO [05-04|09:34:16.027] New block                                index=4 id=3:4:9bb789       gas_used=43426   txs=1/0 age=380.470ms t=1.540ms \n" +
			"INFO [05-04|09:34:16.512] New block                                index=5 id=3:7:a780ce       gas_used=138,470 txs=5/0 age=374.251ms t=3.796ms \n" +
			"INFO [05-04|09:34:17.003] New block                                index=6 id=3:10:d7da0b      gas_used=105,304 txs=4/0 age=381.575ms t=3.249ms \n"

	return strings.NewReader(testLog)
}
