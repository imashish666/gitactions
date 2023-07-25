package server

import (
	lo "log"
	"www-api/internal/constants"
	"www-api/internal/logger"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// databaseConnection takes a db connection string and returns a db instance
func databaseConnection(connectionString string, log logger.ZapLogger) *sqlx.DB {
	conn, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		lo.Println("error connection database", err)
		log.Fatal("connection to database failed", map[string]interface{}{"error": err})
	}
	conn.SetMaxIdleConns(constants.MaxDBIdleConnections)
	conn.SetMaxOpenConns(constants.MaxDBOpenConnections)
	return conn
}

// redisConnection takes a redis address and returns redis client
func redisConnection(redisAddress string, log logger.ZapLogger) *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	if conn == nil {
		log.Fatal("connection to redis server failed", map[string]interface{}{"error": "consider checking redis adderss", "address": redisAddress})
	}
	return conn
}

// redisClusterConnection takes a list redis address and returns redis cluster client
func redisClusterConnection(redisAddresses []string, log logger.ZapLogger) *redis.ClusterClient {
	conn := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: redisAddresses,
	})
	if conn == nil {
		log.Fatal("connection to redis server failed", map[string]interface{}{"error": "consider checking redis addersses", "addresses": redisAddresses})
	}
	return conn
}

// elasticConnection takes a url, username and password and returns elsatic client
func elasticConnection(url, username, password string, log logger.ZapLogger) *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
		Username: "foo",
		Password: "bar",
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		lo.Println("error connection elastic colud", err)
		log.Fatal("connection to elastic cloud failed", map[string]interface{}{"error": err})
	}

	return es
}
