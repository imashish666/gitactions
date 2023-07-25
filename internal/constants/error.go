package constants

import "errors"

var ResourceNotFound = errors.New("key does not exists")

var BlankAtRiskKey = errors.New("atRiskKey missing in request body")
var BlankAtRiskValue = errors.New("atRiskValue missing in request body")
var BlankEmail = errors.New("email missing in request body")
var BlankFid = errors.New("fid missing in request body")
var BlankTimestamp = errors.New("timestamp missing in request body")
var EmptyFid = errors.New("fid not received")

var InvalidKeyValue = errors.New("invalid key value")
var InvalidTimestampValue = errors.New("invalid timestamp value")
var InvalidScoreValue = errors.New("invalid score value retrieved from cache")
var InvalidKeyLength = errors.New("invalid key format, atRiskKey should be seperated by single colon")
var InvalidKeyEmail = errors.New("invalid email address in atRiskKey")
var InvalidKeyTimestamp = errors.New("invalid timestamp in atRiskKey, second part should be numeric")
var InvalidValueScore = errors.New("invalid score value in atRiskValue, first part should be numeric")
var InvalidValueLength = errors.New("invalid value format, atRiskValue should be seperated by two colon")
var InvalidTimestampParam = errors.New("invalid timestamp, should be numeric")
var InvalidEmailParam = errors.New("invalid userEmail")
var InvalidFidParam = errors.New("invalid fid")
var InvalidCoversionToInt = errors.New("invalid value, cannot be converted to int")

var EmptyString = ""
