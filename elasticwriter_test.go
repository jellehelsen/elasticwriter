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

func TestWriteSplitsMultilineData(t *testing.T) {
	assert := assert.New(t)
	ja := jsonassert.New(t)
	requests := make([]*http.Request, 0)
	bodies := make([][]byte, 0, 0)
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)
		body, _ := ioutil.ReadAll(r.Body)
		bodies = append(bodies, body)
		handler(w, r)
	}))
	defer ts.Close()
	client, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{ts.URL},
	})
	assert.NoError(err)

	writer := elasticwriter.ElasticWriter{client, "index"}
	msg := `This is a multiline
message
Hello from line 3`
	size, err := writer.Write([]byte(msg))
	assert.NoError(err)
	assert.Equal(len(requests), 3)
	assert.NoError(err)
	assert.Equal(len(bodies), 3)
	assert.Equal(len(msg), size)
	ja.Assertf(string(bodies[0]), `{
		"timestamp": "<<PRESENCE>>",
		"message": "This is a multiline"
	}`)
}

func Example() {
	writer, err := elasticwriter.New("index")
	if err != nil {
		panic(err)
	}
	_, err = writer.Write([]byte("Hello World!"))
	if err != nil {
		panic(err)
	}
}
