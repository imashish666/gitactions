package server

import (
	"fmt"
	atRisk "www-api/api/at-risk"
	"www-api/api/customer"
	info "www-api/api/student"
	"www-api/config"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// implement different api routes
func AddRoutes(router *gin.Engine, config config.Config, log logger.ZapLogger) {
	connectionStrings := getConnectionString(config)

	at_risk_read_redis := redisConnection(connectionStrings.AtRiskReadRedis, log)
	at_risk_write_redis := redisConnection(connectionStrings.AtRiskWriteRedis, log)

	www_read_redis := redisConnection(connectionStrings.WWWReadRedis, log)
	www_write_redis := redisConnection(connectionStrings.WWWWriteRedis, log)

	at_risk_read_db := databaseConnection(connectionStrings.AtRiskReadDB, log)
	at_risk_write_db := databaseConnection(connectionStrings.AtRiskWriteDB, log)

	schools_read_db := databaseConnection(connectionStrings.SchoolsReadDB, log)
	schools_write_db := databaseConnection(connectionStrings.SchoolsWriteDB, log)

	connections := &datatypes.Connections{
		DB: map[string]*sqlx.DB{
			constants.AtRiskReadDBKey:   at_risk_read_db,
			constants.AtRiskWriteDBKey:  at_risk_write_db,
			constants.SchoolsReadDBKey:  schools_read_db,
			constants.SchoolsWriteDBKey: schools_write_db,
		},
		Redis: map[string]*redis.Client{
			constants.AtRiskReadRedisKey:  at_risk_read_redis,
			constants.AtRiskWriteRedisKey: at_risk_write_redis,
			constants.WWWReadRedisKey:     www_read_redis,
			constants.WWWWriteRedisKey:    www_write_redis,
		},
		Elastic: nil,
	}

	//create instance of NewRiskAPI
	risk := atRisk.NewRiskAPI(config, log, connections)
	//create instance of NewInfoAPI
	student := info.NewInfoAPI(config, log, connections)
	//create instance of CustomerAPI
	cust := customer.NewCustomerAPI(config, log, connections)

	//create main router group
	api := router.Group("/api")
	{
		//create router sub group & attach hanlder functions
		atRisk := api.Group("/atRisk")
		{
			atRisk.POST("/cache/create", risk.CreateCache)
			atRisk.DELETE("/cache/delete", risk.DeleteCache)
			atRisk.POST("/extend-ttl", risk.ExtendTTL)
			atRisk.GET("/score", risk.Score)
			atRisk.GET("/event-score-details", risk.EventScore)
		}

		//create router sub group & attach hanlder functions
		customer := api.Group("/customer")
		{
			customer.GET("/privacy/status", cust.PrivacyStatus)
			customer.GET("/timezone", cust.Timezone)
			customer.GET("/notification/config/aware", cust.Notification)
			customer.GET("/filter-type", cust.FilterType)
		}

		api.GET("/user", student.GetInfo)

	}
}

func getConnectionString(config config.Config) datatypes.ConnectionString {
	riskReadDB := config.Mysql[constants.AtRiskDBKey].Read
	riskWriteDB := config.Mysql[constants.AtRiskDBKey].Write
	schoolReadDB := config.Mysql[constants.SchoolsDBKey].Read
	schoolWriteDB := config.Mysql[constants.SchoolsDBKey].Write
	atRiskReadRedis := config.Redis[constants.AtRiskRedisKey].Read
	atRiskWriteRedis := config.Redis[constants.AtRiskRedisKey].Write
	wwwReadRedis := config.Redis[constants.WWWRedisKey].Read
	wwwWriteRedis := config.Redis[constants.WWWRedisKey].Write

	return datatypes.ConnectionString{
		AtRiskReadDB:     fmt.Sprintf(constants.DBConnectionString, riskReadDB.User, riskReadDB.Password, riskReadDB.Host, riskReadDB.Port, riskReadDB.DBName),
		AtRiskWriteDB:    fmt.Sprintf(constants.DBConnectionString, riskWriteDB.User, riskWriteDB.Password, riskWriteDB.Host, riskWriteDB.Port, riskWriteDB.DBName),
		SchoolsReadDB:    fmt.Sprintf(constants.DBConnectionString, schoolReadDB.User, schoolReadDB.Password, schoolReadDB.Host, schoolReadDB.Port, schoolReadDB.DBName),
		SchoolsWriteDB:   fmt.Sprintf(constants.DBConnectionString, schoolWriteDB.User, schoolWriteDB.Password, schoolWriteDB.Host, schoolWriteDB.Port, schoolWriteDB.DBName),
		AtRiskReadRedis:  fmt.Sprintf(constants.RedisConnectionString, atRiskReadRedis.Host, atRiskReadRedis.Port),
		AtRiskWriteRedis: fmt.Sprintf(constants.RedisConnectionString, atRiskWriteRedis.Host, atRiskWriteRedis.Port),
		WWWReadRedis:     fmt.Sprintf(constants.RedisConnectionString, wwwReadRedis.Host, wwwReadRedis.Port),
		WWWWriteRedis:    fmt.Sprintf(constants.RedisConnectionString, wwwWriteRedis.Host, wwwWriteRedis.Port),
	}
}
