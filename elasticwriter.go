package elasticwriter

import (
	"bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"time"
)

type ElasticWriter struct {
	*elasticsearch.Client
	Index string
}

type elasticLogLine struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
}

func New(index string) (*ElasticWriter, error) {
	es7, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, err
	}
	return &ElasticWriter{es7, index}, nil
}

func (ew *ElasticWriter) Write(data []byte) (int, error) {
	logline := elasticLogLine{time.Now(), string(data)}
	payload, err := json.Marshal(logline)
	if err != nil {
		return 0, err
	}
	_, err = ew.Client.Index(ew.Index, bytes.NewReader(payload))
	if err != nil {
		return 0, err
	}
	return len(data), nil
}
