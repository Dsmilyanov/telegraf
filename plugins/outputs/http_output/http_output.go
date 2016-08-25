package http_output

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/outputs"
)

type HttpOutput struct {
	URLs []string `toml:"urls"`
}

type requestStruct struct {
	Timestamp   time.Time
	Measurement string
	Tags        map[string]string
	Fields      map[string]interface{}
}

var sampleConfig = `
	## An array of endpoint URLs, where telegraf.Metric points will be handled
  urls = ["http://localhost:8080"] # required
`

func (h *HttpOutput) Description() string {
	return "Configuration for sending metrics via http"
}

func (h *HttpOutput) SampleConfig() string {
	return sampleConfig
}

func (h *HttpOutput) Connect() error {
	var urls []string
	for _, u := range h.URLs {
		urls = append(urls, u)
	}
	return nil
}

func (h *HttpOutput) Close() error {
	return nil
}

func (h *HttpOutput) Write(metrics []telegraf.Metric) error {
	for _, u := range h.URLs {
		var rawData = make([]requestStruct, len(metrics))

		for key, metric := range metrics {
			rawData[key].Measurement = metric.Name()
			rawData[key].Timestamp = metric.Time()
			rawData[key].Tags = metric.Tags()
			rawData[key].Fields = metric.Fields()
		}

		marshaledData, err := json.Marshal(rawData)
		if err != nil {
			return err
		}

		resp, err := http.Post(u, "json", bytes.NewBuffer(marshaledData))
		if err != nil {
			return err
		}
		defer resp.Body.Close()
	}
	return nil
}

func init() {
	outputs.Add("http_output", func() telegraf.Output {
		return &HttpOutput{}
	})
}
