package nodemon

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Fantom-foundation/Norma/driver"
	opera "github.com/Fantom-foundation/Norma/driver/node"
)

type PprofData []byte

func GetPprofData(node driver.Node, duration time.Duration) (PprofData, error) {
	url := node.GetHttpServiceUrl(&opera.OperaPprofService)
	if url == nil {
		return nil, fmt.Errorf("node does not offer the pprof service")
	}
	resp, err := http.Get(fmt.Sprintf("%s/debug/pprof/profile?seconds=%d", *url, int(duration.Seconds())))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch result: %v", resp)
	}
	return io.ReadAll(resp.Body)
}
