# ElasticWriter

ElasticWriter implements an io.Writer that writes to [ElasticSearch](https://www.elastic.co/products/elasticsearch).

[![Build
Status](https://travis-ci.com/jellehelsen/elasticwriter.svg?branch=develop)](https://travis-ci.com/jellehelsen/elasticwriter) [![Go Report Card](https://goreportcard.com/badge/github.com/jellehelsen/elasticwriter)](https://goreportcard.com/report/github.com/jellehelsen/elasticwriter)

## Installation

Add the package to your `go.mod` file:

    require github.com/jellehelsen/elasticwriter
    

## Usage

Export the `ELASTICSEARCH_URL` environment variable to the correct URL.

``` golang
import "github.com/jellehelsen/elasticwriter"

writer, err := elasticwriter.New("index_to_write_to")

if err != nil {
  log.Fatalf("Error creating the client: %s", err)
}

```

## License

(c) 2020 Jelle Helsen. Licensed under the Apache License, Version 2.0.
