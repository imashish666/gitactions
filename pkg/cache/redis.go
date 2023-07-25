package cache

import (
	"context"
	"time"
	"www-api/internal/logger"

	"www-api/internal/constants"

	"github.com/redis/go-redis/v9"
)

type RedisOps interface {
	GetValue(key string) (string, error)
	SetWithTTL(key string, value interface{}, ttl time.Duration) error
	GetKeys(key string) ([]string, error)
	Exists(key string) (bool, error)
	Delete(key string) error
	SetTTL(key string, expiry int) error
	SetDB(db int)
}

type Redis struct {
	read  *redis.Client
	write *redis.Client
	log   logger.ZapLogger
	ctx   context.Context
}

// NewRedis returns a instance of NewRedis struct
func NewRedis(read *redis.Client, write *redis.Client, log logger.ZapLogger, ctx context.Context) Redis {
	return Redis{read, write, log, ctx}
}

// GetValue fetches a value mapped to a key
func (r Redis) GetValue(key string) (string, error) {
	value, err := r.read.Get(r.ctx, key).Result()

	if err != nil {
		if err.Error() == "redis: nil" {
			r.log.Error("key not found in redis", map[string]interface{}{"key": key})
			return "", constants.ResourceNotFound
		}
		r.log.Error("unable to fetch key from redis: error ", map[string]interface{}{"key": key, "err": err})
		return "", err
	}
	return value, nil
}

// SetWithTTL sets a key value pair along with a ttl
func (r Redis) SetWithTTL(key string, value interface{}, ttl time.Duration) error {
	_, err := r.write.Set(r.ctx, key, value, ttl).Result()
	if err != nil {
		r.log.Error("unable to set key value in redis", map[string]interface{}{"key": key, "value": value, "err": err})
		return err
	}
	return nil
}

// GetKeys fetches all the keys based on a pattern
func (r Redis) GetKeys(key string) ([]string, error) {
	var cursor uint64
	var result []string
	for {
		keys, cursor, err := r.read.Scan(r.ctx, cursor, key, 0).Result()
		if err != nil {
			r.log.Error("unable to scan redis with key", map[string]interface{}{"key": key, "err": err})
			return []string{}, err
		}

		result = append(result, keys...)

		if cursor == 0 {
			break
		}
	}

	return result, nil
}

// Exists returns whether a key exists in redis
func (r Redis) Exists(key string) (bool, error) {
	occurence, err := r.read.Exists(r.ctx, key).Result()
	if err != nil {
		r.log.Error("error while checking key exists in redis", map[string]interface{}{"key": key, "err": err})
		return false, err
	}

	return occurence != 0, nil
}

// Delete removes a key value pair based on key
func (r Redis) Delete(key string) error {
	_, err := r.write.Del(r.ctx, key).Result()
	if err != nil {
		r.log.Error("error while deleting key", map[string]interface{}{"key": key, "err": err})
		return err
	}

	return nil
}

// SetTTL updates a ttl for a key
func (r Redis) SetTTL(key string, expiry int) error {
	ttl := time.Second * time.Duration(expiry)
	_, err := r.write.Expire(r.ctx, key, ttl).Result()
	if err != nil {
		r.log.Error("error while setting ttl", map[string]interface{}{"key": key, "err": err})
		return err
	}

	return nil
}

func (r Redis) SetDB(db int) {
	r.write.Options().DB = db
	r.read.Options().DB = db
}
