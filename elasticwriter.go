package elasticwriter

import (
	"bytes"
	"strings"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"time"
)

// ElasticWriter implements io.Writer while writing to elasticsearch
type ElasticWriter struct {
	*elasticsearch.Client
	Index string
}

type elasticLogLine struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

// New creates a new writer
//
// It will use http://localhost:9200 as the default elasticsearch address.
//
// It will use the ELASTICSEARCH_URL environment variable, if set,
// to configure the addresses; use a comma to separate multiple URLs.
//
// index is the name of the elasticsearch index to write to
//
func New(index string) (*ElasticWriter, error) {
	es7, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &ElasticWriter{es7, index}, nil
}

func (ew *ElasticWriter) Write(data []byte) (int, error) {
	bytesWritten := 0
	for _, line := range(strings.Split(string(data), "\n")) {
		logline := elasticLogLine{time.Now(), line}
		payload, err := json.Marshal(logline)
		if err != nil {
			return 0, err
		}
		_, err = ew.Client.Index(ew.Index, bytes.NewReader(payload))
		if err != nil {
			return bytesWritten, err
		}
		bytesWritten += len(line)
	}
	// Add the newlines we've stripped
	bytesWritten += strings.Count(string(data), "\n")
	return bytesWritten, nil
}
