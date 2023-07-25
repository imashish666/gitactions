package student

import (
	"context"
	"strconv"
	"strings"
	"time"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/cache"
	"www-api/pkg/database"
	"www-api/pkg/model"
	"www-api/utils"
)

const (
	SCHOOLTYPE_GOOGLE = 1
	SCHOOLTYPE_AZURE  = 2
)

type CustomerService struct {
	log                 logger.ZapLogger
	redis               cache.RedisOps
	getTimezoneFromUser func(fid string) (string, error)
	getNotification     func(fid string) (datatypes.Notification, error)
	getFilter           func(fid string) (datatypes.FilterType, error)
}

// NewCustomerService returns an instance of RiskService struct
func NewCustomerService(log logger.ZapLogger, connections *datatypes.Connections) CustomerService {
	// connections.Redis[constants.AtRiskReadRedisKey].Options().DB = constants.RedisDB15
	// connections.Redis[constants.AtRiskWriteRedisKey].Options().DB = constants.RedisDB15

	readinterface := model.NewReadModel(log, database.NewDatabase(connections.DB[constants.SchoolsReadDBKey]))
	//writeinterface := model.NewWriteModel(log, database.NewDatabase(writeconn))

	return CustomerService{
		log:                 log,
		redis:               cache.NewRedis(connections.Redis[constants.WWWReadRedisKey], connections.Redis[constants.WWWWriteRedisKey], log, context.Background()),
		getTimezoneFromUser: readinterface.GetUserTimezone,
		getNotification:     readinterface.GetAwareNotification,
		getFilter:           readinterface.GetFilterType,
	}
}

// ProuctPrivacyStatus gets info of a student based on email and fid (if available)
func (s CustomerService) ProuctPrivacyStatus(fid string) (map[string]int, error) {
	privacyMode := map[string]int{
		"Filter":    0,
		"Aware":     0,
		"24":        0,
		"Responder": 0,
		"suppBully": 0,
	}
	domainName := strings.Split(fid, "@")
	if domainName[1] == "" {
		s.log.Error("empty domain", nil)
		return privacyMode, constants.EmptyFid
	}
	s.redis.SetDB(constants.RedisDB21)
	privacyKey := domainName[1] + ":PF:ENHANCED_PRIVACY"
	value, err := s.redis.GetValue(privacyKey)
	if err != nil {
		if err == constants.ResourceNotFound {
			return privacyMode, nil
		}
		s.log.Error("error fetching value from redis", map[string]interface{}{"error": err})
		return privacyMode, err
	}

	privacyValue, err := strconv.Atoi(value)
	if err != nil {
		s.log.Error("error converting to int", map[string]interface{}{"error": err})
		return privacyMode, constants.InvalidCoversionToInt
	}
	s.log.Info("value", map[string]interface{}{"privacy": privacyValue})

	if utils.IsBitSet(privacyValue, constants.FILTER_PRIVACY) {
		privacyMode["Filter"] = 1
	} else if utils.IsBitSet(privacyValue, constants.AWARE_PRIVACY) {
		privacyMode["Aware"] = 1
	} else if utils.IsBitSet(privacyValue, constants.TWENTY_FOUR_PRIVACY) {
		privacyMode["24"] = 1
	} else if utils.IsBitSet(privacyValue, constants.RESPONDER_PRIVACY) {
		privacyMode["Responder"] = 1
	} else if utils.IsBitSet(privacyValue, constants.SUPPRESS_BULLY) {
		privacyMode["suppBully"] = 1
	}

	if privacyMode["Aware"] == 1 && (privacyMode["24"] != 1 || privacyMode["Responder"] != 1) {
		has24 := 0
		hasRespond := 0
		integrationEnabled := 0

		pnBitVector, err := s.redis.GetValue(fid + ":PN")
		if err != nil {
			if err == constants.ResourceNotFound {
				return privacyMode, nil
			}
			s.log.Error("error fetching value from redis", map[string]interface{}{"error": err})
			return privacyMode, err
		}
		pnBitVectorValue, err := strconv.Atoi(pnBitVector)
		if err != nil {
			s.log.Error("error converting to int", map[string]interface{}{"error": err})
			return privacyMode, constants.InvalidCoversionToInt
		}

		pfBitVector, err := s.redis.GetValue(fid + ":PF:3")
		if err != nil {
			if err == constants.ResourceNotFound {
				return privacyMode, nil
			}
			s.log.Error("error fetching value from redis", map[string]interface{}{"error": err})
			return privacyMode, err
		}
		pfBitVectorValue, err := strconv.Atoi(pfBitVector)
		if err != nil {
			s.log.Error("error converting to int", map[string]interface{}{"error": err})
			return privacyMode, constants.InvalidCoversionToInt
		}

		respondVector, err := s.redis.GetValue(fid + ":PN:RESPONDER")
		if err != nil {
			if err == constants.ResourceNotFound {
				return privacyMode, nil
			}
			s.log.Error("error fetching value from redis", map[string]interface{}{"error": err})
			return privacyMode, err
		}
		respondVectorValue, err := strconv.Atoi(respondVector)
		if err != nil {
			s.log.Error("error converting to int", map[string]interface{}{"error": err})
			return privacyMode, constants.InvalidCoversionToInt
		}
		if (utils.IsBitSet(pnBitVectorValue, 9) && utils.IsBitSet(pfBitVectorValue, 1) && utils.IsBitSet(pfBitVectorValue, 0)) ||
			(utils.IsBitSet(pnBitVectorValue, 10) && utils.IsBitSet(pfBitVectorValue, 2) && utils.IsBitSet(pfBitVectorValue, 0)) {
			has24 = 1
		}
		if utils.IsBitSet(respondVectorValue, 0) {
			hasRespond = 1
		}
		if utils.IsBitSet(pfBitVectorValue, 3) {
			integrationEnabled = 1
		}
		if has24 == 1 && hasRespond == 1 && integrationEnabled == 1 {
			privacyMode["24"] = 1
			privacyMode["Responder"] = 1
		}
	}

	return privacyMode, nil
}

// Timezone gets info of a student based on email and fid (if available)
func (s CustomerService) Timezone(fid string) (datatypes.TimezoneResponse, error) {
	location, err := s.getTimezoneFromUser(fid)
	if err != nil {
		s.log.Error("error occured while fetching timezone", map[string]interface{}{"error": err})
		return datatypes.TimezoneResponse{}, err
	}
	loc, err := time.LoadLocation(location)
	if err != nil {
		s.log.Error("location not found", map[string]interface{}{"error": err})
		return datatypes.TimezoneResponse{Tz: location, TzAbbr: ""}, nil
	}
	timezone, _ := time.Now().In(loc).Zone()
	return datatypes.TimezoneResponse{Tz: location, TzAbbr: timezone}, nil
}

// Notification gets notification based on fid
func (s CustomerService) Notification(fid string) (datatypes.Notification, error) {
	notification, err := s.getNotification(fid)
	if err != nil {
		s.log.Error("error occured while fetching timezone", map[string]interface{}{"error": err})
		return datatypes.Notification{}, err
	}

	return notification, nil
}

// Notification gets notification based on fid
func (s CustomerService) GetFilterType(fid string) (string, error) {
	filter, err := s.getFilter(fid)
	if err != nil {
		s.log.Error("error occured while fetching filter type", map[string]interface{}{"error": err})
		return "", err
	}

	if filter.SchoolType == SCHOOLTYPE_AZURE && len(filter.AdIntranet) != 0 {
		return "secGrp", nil
	}
	return "ou", nil
}
