package prometheusmon

import (
	"bytes"
	"text/template"

	"github.com/Fantom-foundation/Norma/driver"
)

// promCfg is the default Prometheus configuration.
const promCfg = `
global:
  scrape_interval: 5s
  evaluation_interval: 5s
scrape_configs:
  - job_name: "opera"
    metrics_path: "/debug/metrics/prometheus"
    file_sd_configs:
      - files:
         - "/etc/prometheus/opera-*.json"
`

// promTargetCfgTmpl is the Prometheus target configuration template.
const promTargetCfgTmpl = `
[
  {
    "targets": ["{{.Host}}:{{.Port}}"],
    "labels": {
      "job": "opera",
      "label": "{{.Label}}"
    }
  }
]
`

// promTargetConfig is the Prometheus target configuration.
type promTargetConfig struct {
	Host  string
	Port  int
	Label string
}

// renderConfigForNode renders the Prometheus configuration for a node.
func renderConfigForNode(node driver.Node) (string, error) {
	cfg := promTargetConfig{
		Host:  node.Hostname(),
		Port:  node.MetricsPort(),
		Label: node.GetLabel(),
	}
	tmpl, err := template.New("promTargetCfg").Parse(promTargetCfgTmpl)
	if err != nil {
		return "", err
	}
	var configBuffer bytes.Buffer
	err = tmpl.Execute(&configBuffer, cfg)
	if err != nil {
		return "", err
	}
	return configBuffer.String(), nil
}
