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

package monitoring

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timestampReg      = regexp.MustCompile(`\[\S*\]`)
	blockReg          = regexp.MustCompile(`index=\d*`)
	gasReg            = regexp.MustCompile(`gas_used=\S*`)
	gasRateReg        = regexp.MustCompile(`gas_rate=\d+(\.\d*)?`)
	baseFeeReg        = regexp.MustCompile(`base_fee=\d+`)
	txsReg            = regexp.MustCompile(`txs=\d+`)
	processingTimeReg = regexp.MustCompile(`t=\S*`)
)

// NewLogReader creates a channel and reads logs from the input reader, sending it to the channel.
// The reader is expected to contain Opera Log stream, which is parsed and converted into Block struct.
// The Blocks are sent to the output channel.
func NewLogReader(reader io.Reader) <-chan Block {
	ch := make(chan Block, 10)
	go func() {
		defer close(ch)
		if err := readBlocks(reader, ch); err != nil {
			log.Printf("error: %s", err)
		}
	}()

	return ch
}

// readBlocks reads the input reader, which is expected to contain Opera Log.
// The log is parsed and information about produced blocks is sent to the input channel.
func readBlocks(reader io.Reader, ch chan<- Block) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "New block") {
			block, err := parseBlock(line)
			if err != nil {
				return err
			}

			ch <- block
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// parseTime convert time from log format into Time type.
func parseTime(str string) (time.Time, error) {
	year := time.Now().Year()
	return time.Parse("2006-[01-02|15:04:05.000]", fmt.Sprintf("%d-%s", year, str))
}

// parseBlock parses block information from the log line. It is expected the log line is well-formed.
func parseBlock(line string) (block Block, err error) {
	// example line: "INFO [05-04|09:34:15.537] New block index=3 id=3:1:3d6fb6 gas_used=117,867 base_fee=123 txs=1/0 age=343.255ms t=1.579ms
	timestampStr := timestampReg.FindString(line)
	blockNumberStr := strings.Split(blockReg.FindString(line), "=")[1]
	gasUsedStr := strings.ReplaceAll(strings.Split(gasReg.FindString(line), "=")[1], ",", "")
	gasRateStr := strings.Split(gasRateReg.FindString(line), "=")[1]
	baseFeeStr := strings.Split(baseFeeReg.FindString(line), "=")[1]
	txsStr := strings.Split(txsReg.FindString(line), "=")[1]
	processingTimeStr := strings.Trim(strings.Split(processingTimeReg.FindString(line), "=")[1], "\"")

	blockNumber, err := strconv.Atoi(blockNumberStr)
	if err != nil {
		return block, err
	}

	timestamp, err := parseTime(timestampStr)
	if err != nil {
		return block, err
	}

	txs, err := strconv.Atoi(txsStr)
	if err != nil {
		return block, err
	}

	gasUsed, err := strconv.Atoi(gasUsedStr)
	if err != nil {
		return block, err
	}

	gasRate, err := strconv.ParseFloat(gasRateStr, 64)
	if err != nil {
		return block, err
	}

	baseFeeUsed, err := strconv.Atoi(baseFeeStr)
	if err != nil {
		return block, err
	}

	processingTime, err := time.ParseDuration(processingTimeStr)
	if err != nil {
		return block, err
	}

	return Block{
		Height:         blockNumber,
		Txs:            txs,
		GasUsed:        gasUsed,
		Time:           timestamp,
		ProcessingTime: processingTime,
		GasBaseFee:     baseFeeUsed,
		GasRate:        gasRate,
	}, nil
}

// Block contains data of one block
type Block struct {
	Height         int
	Time           time.Time     // timestamp of the block
	Txs            int           // number of transactions in block
	GasUsed        int           // gas used in the block
	ProcessingTime time.Duration // block processing time
	GasBaseFee     int           // gas base fee in the block
	GasRate        float64       // gas rate in gas/s in the block
}
