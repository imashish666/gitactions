package atrisk

import (
	"context"
	"strconv"
	"strings"
	"time"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/cache"
	"www-api/pkg/model"

	"www-api/pkg/database"

	_ "github.com/go-sql-driver/mysql"
)

type RiskService struct {
	log            logger.ZapLogger
	redis          cache.RedisOps
	getAtRiskScore func(email string) ([]datatypes.RiskScore, error)
}

// NewRiskService returns an instance of RiskService struct
func NewRiskService(log logger.ZapLogger, connections *datatypes.Connections) RiskService {
	readinterface := model.NewReadModel(log, database.NewDatabase(connections.DB[constants.AtRiskReadDBKey]))
	// writeinterface := model.NewWriteModel(log, database.NewDatabase(writeconn))

	return RiskService{
		log:            log,
		redis:          cache.NewRedis(connections.Redis[constants.AtRiskReadRedisKey], connections.Redis[constants.AtRiskWriteRedisKey], log, context.Background()),
		getAtRiskScore: readinterface.GetAtRiskScore,
	}
}

// CreateCache sets a key value pair in redis and returns total score for that email
func (s RiskService) CreateCache(key, value string) (datatypes.AtRiskResponse, error) {
	s.log.Info("setting cache", map[string]interface{}{"key": key, "value": value})
	//setting ttl as 60 days i.e. 5184000 secs
	err := s.redis.SetWithTTL(key, value, 5184000*time.Second)
	if err != nil {
		s.log.Error("error occured while setting cache value", map[string]interface{}{"error": err})
		return datatypes.AtRiskResponse{}, err
	}

	email := strings.Split(key, ":")[0]
	score, err := s.getTotalAtRiskScore(email)
	if err != nil {
		s.log.Error("error occured while getting total score", map[string]interface{}{"error": err})
		return datatypes.AtRiskResponse{}, err

	}

	return datatypes.AtRiskResponse{AtRiskScore: score}, nil
}

// DeleteCache returns the total score after removing the key value from redis
func (s RiskService) DeleteCache(key string) (datatypes.AtRiskResponse, error) {
	exists, err := s.redis.Exists(key)
	if err != nil {
		s.log.Error("error occured while checking if key exists in redis", map[string]interface{}{"key": key, "error": err})
		return datatypes.AtRiskResponse{}, err
	}

	if !exists {
		s.log.Error("AT_RISK_SCORE_NOT_FOUND. Cannnot unassign as score is not assigned to user at all. unassignAtRiskKey", map[string]interface{}{"key": key})
		return datatypes.AtRiskResponse{}, constants.ResourceNotFound
	}

	err = s.redis.Delete(key)
	if err != nil {
		s.log.Error("error occured while deleting key", map[string]interface{}{"key": key, "error": err})
		return datatypes.AtRiskResponse{}, err
	}

	email := strings.Split(key, ":")[0]
	score, err := s.getTotalAtRiskScore(email)
	if err != nil {
		s.log.Error("error occured while getting total score", map[string]interface{}{"key": key, "error": err})
		return datatypes.AtRiskResponse{}, err
	}

	return datatypes.AtRiskResponse{AtRiskScore: score}, err
}

// GetScore fetches risk score from the database based on email
func (s RiskService) GetScore(email string) ([]datatypes.RiskScore, error) {
	scores, err := s.getAtRiskScore(email)
	if err != nil {
		s.log.Error("unable to fetch scores from database", map[string]interface{}{"error": err})
		return nil, err
	}
	return scores, nil
}

// ExtendTTL updates the ttl value for all keys with email pattern
func (s RiskService) ExtendTTL(email string, ttl int) error {
	keys, err := s.redis.GetKeys(email + ":*")
	if err != nil {
		s.log.Error("unable to fetch all keys from redis: error", map[string]interface{}{"email": email, "error": err})
		return err
	}

	for _, key := range keys {
		err = s.redis.SetTTL(key, ttl)
		if err != nil {
			s.log.Error("unable to set expiry in redis", map[string]interface{}{"key": key, "error": err})
			return err
		}
	}

	return nil
}

// GetEventScore returns key, value & score for a specific event based on timestamp
func (s RiskService) GetEventScore(email, timestamp, mid string) (datatypes.EventScoreResponse, error) {
	atRiskKey := email + ":" + timestamp
	exists, err := s.redis.Exists(atRiskKey)
	if err != nil {
		s.log.Error("error checking if key exists in redis", map[string]interface{}{"key": atRiskKey, "error": err})
		return datatypes.EventScoreResponse{}, err
	}

	if !exists {
		timestp, err := strconv.Atoi(timestamp)
		if err != nil {
			s.log.Error("error converting timestamp into int", map[string]interface{}{"error": err})
			return datatypes.EventScoreResponse{}, constants.InvalidTimestampValue
		}
		newTime := strconv.Itoa(timestp / 1000)
		atRiskKey = email + ":" + newTime
	}

	atRiskValue, err := s.redis.GetValue(atRiskKey)
	if err != nil {
		s.log.Error("error fetching value in redis", map[string]interface{}{"error": err})
		return datatypes.EventScoreResponse{}, err
	}

	if atRiskValue == "" && mid != "" && mid[0] == '<' {
		allKeys, err := s.redis.GetKeys(email + ":*")
		if err != nil {
			s.log.Error("unable to fetch all keys from redis", map[string]interface{}{"email": email, "error": err})
			return datatypes.EventScoreResponse{}, err
		}

		for _, key := range allKeys {
			redisValue, err := s.redis.GetValue(key)
			if err != nil {
				s.log.Error("unable to fetch score from redis", map[string]interface{}{"key": key, "error": err})
				return datatypes.EventScoreResponse{}, err
			}
			if strings.Contains(redisValue, mid) {
				atRiskValue = redisValue
				atRiskKey = key
				s.log.Info("at risk score found by looking into value", map[string]interface{}{"atRiskKey": atRiskKey})
				break
			}
		}
	}

	atRiskScoreString := strings.Split(atRiskValue, ":")[0]
	atRiskScore, err := strconv.Atoi(atRiskScoreString)
	if err != nil {
		s.log.Error("unable to convert redis score value to int", map[string]interface{}{"error": err})
		return datatypes.EventScoreResponse{}, constants.InvalidScoreValue
	}

	return datatypes.EventScoreResponse{
		AtRiskKey:   atRiskKey,
		AtRiskValue: atRiskValue,
		AtRiskScore: atRiskScore,
	}, nil

}

// getTotalAtRiskScore return the total score of all events based on email
func (s RiskService) getTotalAtRiskScore(email string) (int, error) {
	keys, err := s.redis.GetKeys(email + ":*")
	if err != nil {
		s.log.Error("unable to fetch all keys from redis", map[string]interface{}{"email": email, "error": err})
		return 0, err
	}
	var totalScore int
	for _, key := range keys {
		redisValue, err := s.redis.GetValue(key)
		if err != nil {
			s.log.Error("unable to fetch score from redis ", map[string]interface{}{"key": key, "error": err})
			return 0, err
		}

		scoreString := strings.Split(redisValue, ":")
		score, err := strconv.Atoi(scoreString[0])
		if err != nil {
			s.log.Error("unable to convert redis score value to int", map[string]interface{}{"error": err})
			return 0, constants.InvalidKeyValue
		}
		totalScore += score
	}

	return totalScore, nil
}
