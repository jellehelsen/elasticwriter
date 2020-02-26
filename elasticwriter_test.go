package elasticwriter_test

import (
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/jellehelsen/elasticwriter"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
)

func TestCreation(t *testing.T) {
	var (
		writer io.Writer
		err    error
	)
	writer, err = elasticwriter.New("index")
	assert.NotNil(t, writer)
	assert.NoError(t, err)
}

func TestWrite(t *testing.T) {
	assert := assert.New(t)
	ja := jsonassert.New(t)
	requests := make([]*http.Request, 0)
	body := make([]byte, 0)
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)
		body, _ = ioutil.ReadAll(r.Body)
		handler(w, r)
	}))
	defer ts.Close()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{ts.URL},
	})
	assert.NoError(err)

	writer := elasticwriter.ElasticWriter{client, "index"}
	size, err := writer.Write([]byte("hello"))
	assert.NoError(err)
	assert.NotEqual(len(requests), 0)
	assert.NoError(err)
	assert.NotEqual(len(body), 0)
	assert.NotEqual(size, 0)
	ja.Assertf(string(body), `{
		"timestamp": "<<PRESENCE>>",
		"message": "hello"
	}`)
}
