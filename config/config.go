package config

import (
	_ "embed"
	"log"
	"os"
	"www-api/internal/constants"
	"www-api/internal/logger"
	"www-api/utils"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Region     string
	Deployment string
	Server     server
	Mysql      map[string]mysqlDatabase
	Redis      map[string]redisDatabase
	Elastic    elasticvariables
	RedisPort  string
}

type elasticvariables struct {
	Username string
	Password string
	Host     string
	Port     string
}

type dbvariables struct {
	User     string
	Password string
	Host     string
	Port     string
	DBName   string
}

type mysqlDatabase struct {
	Read  dbvariables
	Write dbvariables
}

type redisvariables struct {
	Host string
	Port string
}

type redisDatabase struct {
	Read  redisvariables
	Write redisvariables
}

type server struct {
	Host string
	Port string
}

// LoadConfig returns Config struct after reading the config file
func LoadConfig(filePath, region, deployment, secretName string, logger logger.ZapLogger) (Config, error) {
	if region != "" && secretName != "" {
		secrets := utils.FetchAWSSecrets(region, secretName, logger)

		return Config{
			Region:     region,
			Deployment: deployment,
			Server: server{
				Host: "",
				Port: secrets["server-port"],
			},
			Mysql: map[string]mysqlDatabase{
				constants.AtRiskDBKey: {
					Read: dbvariables{
						User:     secrets["globals-securly_atrisk_read_username"],
						Password: secrets["globals-securly_atrisk_read_password"],
						Host:     secrets["globals-securly_atrisk_read_host"],
						Port:     secrets["globals-securly_atrisk_read_port"],
						DBName:   constants.AtRiskDBName,
					},
					Write: dbvariables{
						User:     secrets["globals-securly_atrisk_write_username"],
						Password: secrets["globals-securly_atrisk_write_password"],
						Host:     secrets["globals-securly_atrisk_write_host"],
						Port:     secrets["globals-securly_atrisk_write_port"],
						DBName:   constants.AtRiskDBName,
					},
				},
				constants.SchoolsDBKey: {
					Read: dbvariables{
						User:     secrets["globals-securly_schools_read_username"],
						Password: secrets["globals-securly_schools_read_password"],
						Host:     secrets["globals-securly_schools_read_host"],
						Port:     secrets["globals-securly_schools_read_port"],
						DBName:   constants.SchoolsDBName,
					},
					Write: dbvariables{
						User:     secrets["globals-securly_schools_write_username"],
						Password: secrets["globals-securly_schools_write_password"],
						Host:     secrets["globals-securly_schools_write_host"],
						Port:     secrets["globals-securly_schools_write_port"],
						DBName:   constants.SchoolsDBName,
					},
				},
			},
			Redis: map[string]redisDatabase{
				constants.AtRiskRedisKey: {
					Read: redisvariables{
						Host: secrets["globals-www_redis_host"],
						Port: secrets["globals-www_redis_port"],
					},
					Write: redisvariables{
						Host: secrets["globals-www_redis_host"],
						Port: secrets["globals-www_redis_port"],
					},
				},
			},
			Elastic: elasticvariables{
				Username: secrets["globals-elastic_cloud_user"],
				Password: secrets["globals-elastic_cloud_password"],
				Host:     secrets["globals-elastic_cloud_host"],
				Port:     secrets["globals-elastic_cloud_port"],
			},
		}, nil
	}

	var config Config
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Config file not found at given location, %v\n", err)
		return config, err
	}

	err = yaml.Unmarshal(content, &config)
	if err != nil {
		log.Printf("Unable to decode into struct, %v\n", err)
		return config, err
	}
	return config, nil
}
