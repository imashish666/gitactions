package test

import "errors"

var CacheSetErr = errors.New("error setting cache")
var CacheGetKeysErr = errors.New("error getting keys from cache")
var CacheGetValueErr = errors.New("error getting values from cache")
var CacheKeyExistsErr = errors.New("error checking if key exists in cache")
var CacheDeleteKeyErr = errors.New("error checking if key exists in cache")
var CacheSetTTLErr = errors.New("error setting ttl for cache keys")
var DBSomethingWentWrongErr = errors.New("error something went wrong")
var InternalServerErr = errors.New("internal server error")
