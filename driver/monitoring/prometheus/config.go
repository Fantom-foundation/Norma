package prometheusmon

// promCfg is the default PrometheusDockerNode configuration.
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

// promTargetCfgTmpl is the PrometheusDockerNode target configuration template.
const promTargetCfgTmpl = `
[
  {
    "targets": ["%s:%d"],
    "labels": {
      "job": "opera",
      "label": "%s"
    }
  }
]
`
