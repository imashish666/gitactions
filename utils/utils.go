package utils

import (
	"context"
	"encoding/json"
	"net/mail"
	"strconv"
	"strings"
	"www-api/internal/constants"
	"www-api/internal/logger"

	awsec2metadata "github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func FetchAWSSecrets(region, secretName string, log logger.ZapLogger) map[string]string {

	config, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatal("unable to load default config", map[string]interface{}{"error": err})
	}
	// Create Secrets Manager client
	svc := secretsmanager.NewFromConfig(config)

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     aws.String(secretName),
		VersionStage: aws.String("AWSCURRENT"), // VersionStage defaults to AWSCURRENT if unspecified
	}

	result, err := svc.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal("unable to fetch secrets from secret manager", map[string]interface{}{"error": err})
	}

	var creds map[string]string
	err = json.Unmarshal([]byte(*result.SecretString), &creds)
	if err != nil {
		log.Fatal("unable to unmarshal fetched secret", map[string]interface{}{"error": err})
	}

	return creds
}

func ValidateAtRiskKey(key string, log logger.ZapLogger) error {
	if key == "" {
		log.Error("blank atRiskKey query param", map[string]interface{}{"atRiskKey": key})
		return constants.BlankAtRiskKey
	}

	split := strings.Split(key, ":")
	if len(split) != 2 {
		log.Error("atRiskKey should be seperated by single colon", map[string]interface{}{"atRiskKey": key})
		return constants.InvalidKeyLength
	}

	_, err := mail.ParseAddress(split[0])
	if err != nil {
		log.Error("first part of atRiskKey should have valid email address", map[string]interface{}{"atRiskKey": key})
		return constants.InvalidKeyEmail
	}

	_, err = strconv.Atoi(split[1])
	if err != nil {
		log.Error("second part of atRiskKey should be numeric (timestamp)", map[string]interface{}{"error": err, "atRiskKey": key})
		return constants.InvalidKeyTimestamp
	}

	return nil
}

func ValidateAtRiskValue(value string, log logger.ZapLogger) error {
	if value == "" {
		log.Error("blank atRiskValue query param", map[string]interface{}{"atRiskValue": value})
		return constants.BlankAtRiskValue
	}

	split := strings.Split(value, ":")
	if len(split) != 3 {
		log.Error("atRiskValue should be seperated by two colons", map[string]interface{}{"atRiskValue": value})
		return constants.InvalidKeyLength
	}

	_, err := strconv.Atoi(split[0])
	if err != nil {
		log.Error("first part of atRiskValue should be numeric", map[string]interface{}{"error": err, "atRiskValue": value})
		return constants.InvalidValueScore
	}

	return nil
}

func ValidateTimestamp(timestamp string, log logger.ZapLogger) error {
	if timestamp == "" {
		log.Error("blank timestamp query param", map[string]interface{}{"timestamp": timestamp})
		return constants.BlankTimestamp
	}

	_, err := strconv.Atoi(timestamp)
	if err != nil {
		log.Error("invalid timestamp query param", map[string]interface{}{"timestamp": timestamp})
		return constants.InvalidTimestampParam
	}

	return nil
}

func ValidateEmail(email string, log logger.ZapLogger) error {
	if email == "" {
		log.Error("blank userEmail query param", map[string]interface{}{"email": email})
		return constants.BlankEmail
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		log.Error("invalid email address passed in query param", map[string]interface{}{"error": err, "email": email})
		return constants.InvalidEmailParam
	}

	return nil
}

func ValidateFid(fid string, log logger.ZapLogger) error {
	if fid == "" {
		log.Error("blank fid in request body", map[string]interface{}{"fid": fid})
		return constants.BlankFid
	}

	_, err := mail.ParseAddress(fid)
	if err != nil {
		log.Error("invalid email address in request body", map[string]interface{}{"error": err, "fid": fid})
		return constants.InvalidFidParam
	}

	return nil
}

func IsBitSet(val int, idx int) bool {
	if idx < 0 {
		return false
	}
	return val&(1<<(idx)) > 0
}

func GetRegion() (region string, err error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", err
	}
	c := awsec2metadata.New(sess)
	identity, err := c.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}
	return identity.Region, nil
}
