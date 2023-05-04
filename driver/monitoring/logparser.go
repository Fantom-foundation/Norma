package monitoring

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// LogParserThroughput implementation gets information by parsing Opera logs
type LogParserThroughput struct {
	file string // a file with the logs

	data     map[int]block // map block number -> block, for now it stores all data in-memory
	lastSize int64         // size of the file when it  was last read
}

// block contains data of one block
type block struct {
	time    time.Time // timestamp of the block
	txs     int       // number of transactions in block
	gasUsed int       // gas used in the block
}

// NewLogParserThroughput creates a new instance of the log parser
func NewLogParserThroughput(file string) *LogParserThroughput {
	return &LogParserThroughput{
		file:     file,
		lastSize: -1,
		data:     make(map[int]block, 10000),
	}
}

func (p *LogParserThroughput) GetTransactions(block int) (int, error) {
	b, err := p.getBlock(block)
	if err != nil {
		return 0, err
	}

	return b.txs, nil
}

func (p *LogParserThroughput) GetGas(block int) (int, error) {
	b, err := p.getBlock(block)
	if err != nil {
		return 0, err
	}

	return b.gasUsed, nil
}

func (p *LogParserThroughput) GetBlockTime(block int) (time.Time, error) {
	b, err := p.getBlock(block)
	if err != nil {
		return time.UnixMilli(0), err
	}

	return b.time, nil
}

func (p *LogParserThroughput) GetBlockDelay(block int) (time.Duration, error) {
	b, err := p.getBlock(block)
	if err != nil {
		return 0, err
	}

	prev, err := p.getBlock(block - 1)
	if err != nil {
		return 0, err
	}

	return b.time.Sub(prev.time), nil
}

// getBlock reads block information from internal cache first. If the block is not there,
// and the log file was modified since last read, it parses the log and tries to read the block again.
func (p *LogParserThroughput) getBlock(blockNum int) (block block, err error) {
	if b, exists := p.data[blockNum]; exists {
		return b, nil
	}

	if err := p.readFile(); err != nil {
		return block, err
	}

	// get from cache again
	if b, exists := p.data[blockNum]; exists {
		return b, nil
	}

	return block, ErrNotFound
}

// readFile reads the log file from start to end and converts lines to block information
// which is stored in-memory
func (p *LogParserThroughput) readFile() error {
	file, err := os.Open(p.file)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	size := stat.Size()
	if p.lastSize == size {
		return ErrNotFound // block not found, file was already read before
	}

	p.lastSize = size

	timestampReg := regexp.MustCompile(`\[\S*\]`)
	blockReg := regexp.MustCompile(`index=\d*`)
	gasReg := regexp.MustCompile(`gas_used=\S*`)
	txsReg := regexp.MustCompile(`txs=\d+`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// example line: "INFO [05-04|09:34:15.537] New block index=3 id=3:1:3d6fb6 gas_used=117,867 txs=1/0 age=343.255ms t=1.579ms
		line := scanner.Text()
		if strings.Contains(line, "New block") {
			timestampStr := timestampReg.FindString(line)
			blockNumberStr := strings.Split(blockReg.FindString(line), "=")[1]
			gasUsedStr := strings.ReplaceAll(strings.Split(gasReg.FindString(line), "=")[1], ",", "")
			txsStr := strings.Split(txsReg.FindString(line), "=")[1]

			blockNumber, err := strconv.Atoi(blockNumberStr)
			if err != nil {
				return err
			}

			timestamp, err := parseTime(timestampStr)
			if err != nil {
				return err
			}

			txs, err := strconv.Atoi(txsStr)
			if err != nil {
				return err
			}

			gasUsed, err := strconv.Atoi(gasUsedStr)
			if err != nil {
				return err
			}

			p.data[blockNumber] = block{
				time:    timestamp,
				txs:     txs,
				gasUsed: gasUsed,
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// parseTime convert time from log format into Time type
func parseTime(str string) (time.Time, error) {

	return time.Parse("[01-02|15:04:05.000]", str)
}
