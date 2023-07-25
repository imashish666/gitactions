package elastic

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type ElasticClient struct {
	con *elasticsearch.Client
}

func NewElasticClient(e *elasticsearch.Client) ElasticClient {
	return ElasticClient{
		con: e,
	}
}

type ElasticActions interface {
	Info() (*esapi.Response, error)
	Search() (*esapi.Response, error)
}

func (c ElasticClient) Info() (*esapi.Response, error) {
	return c.con.Info()
}

func (c ElasticClient) Search() (*esapi.Response, error) {
	return c.con.Search()
}
