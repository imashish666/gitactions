package datatypes

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ConnectionString struct {
	AtRiskReadDB     string
	AtRiskWriteDB    string
	SchoolsReadDB    string
	SchoolsWriteDB   string
	AtRiskReadRedis  string
	AtRiskWriteRedis string
	WWWReadRedis     string
	WWWWriteRedis    string
}

type Connections struct {
	DB      map[string]*sqlx.DB
	Redis   map[string]*redis.Client
	Elastic map[string]*elasticsearch.Client
}
