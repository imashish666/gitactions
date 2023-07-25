package atrisk

import (
	"testing"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/cache/mocks"
	"www-api/test"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewRiskService(t *testing.T) {
	type tests struct {
		name        string
		log         logger.ZapLogger
		connections *datatypes.Connections
	}

	testCases := []tests{
		{
			name: "valid case",
			log:  logger.ZapLogger{},
			connections: &datatypes.Connections{
				DB: map[string]*sqlx.DB{
					constants.SchoolsReadDBKey: &sqlx.DB{},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := NewRiskService(tc.log, tc.connections)
			if riskService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, riskService.log)
			}
			if riskService.redis == nil {
				t.Errorf("expected redis interface but got nil")
			}
			if riskService.getAtRiskScore == nil {
				t.Errorf("expected getAtRiskScore but got nil")
			}
		})
	}
}

func TestCreateCache(t *testing.T) {

	type tests struct {
		name        string
		redisClient func() *mocks.RedisOps
		wantScore   datatypes.AtRiskResponse
		wantErr     error
	}

	testCases := []tests{
		{
			name: "valid case",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("SetWithTTL", mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1", "key2"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("45:scan:1dc13ds5c1651", nil).Once()
				moc.On("GetValue", mock.Anything).Return("27:docs:1c5sd16c51dd8", nil).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   nil,
		},
		{
			name: "fail case, error setting cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("SetWithTTL", mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(test.CacheSetErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheSetErr,
		},
		{
			name: "fail case, error getting cache keys",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("SetWithTTL", mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{}, test.CacheGetKeysErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheGetKeysErr,
		},
		{
			name: "fail case, error getting cache value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("SetWithTTL", mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", test.CacheGetValueErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheGetValueErr,
		},
		{
			name: "fail case, error invalid cache value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("SetWithTTL", mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("dede:scan:1dc13ds5c1651", nil).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   constants.InvalidKeyValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := RiskService{logger.ZapLogger{Logger: zap.NewExample()}, tc.redisClient(), nil}
			score, err := risk.CreateCache("key", "value")
			if tc.wantScore != score {
				t.Errorf("expected score %d got %d", tc.wantScore, score)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestDeleteCache(t *testing.T) {

	type tests struct {
		name        string
		log         logger.ZapLogger
		redisClient func() *mocks.RedisOps
		wantScore   datatypes.AtRiskResponse
		wantErr     error
	}

	testCases := []tests{
		{
			name: "valid case",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("Delete", mock.Anything).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1", "key2"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("60:scan:1dc13ds5c1651", nil).Once()
				moc.On("GetValue", mock.Anything).Return("32:docs:1c5sd16c51dd8", nil).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   nil,
		},
		{
			name: "fail case, error checking if key exists in cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, test.CacheKeyExistsErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheKeyExistsErr,
		},
		{
			name: "fail case, key doesn't exists in cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, nil).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   constants.ResourceNotFound,
		},
		{
			name: "fail case, error deleting key from cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("Delete", mock.Anything).Return(test.CacheDeleteKeyErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheDeleteKeyErr,
		},
		{
			name: "fail case, error getting cache keys",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("Delete", mock.Anything).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{}, test.CacheGetKeysErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheGetKeysErr,
		},
		{
			name: "fail case, error getting cache value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("Delete", mock.Anything).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", test.CacheGetValueErr).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   test.CacheGetValueErr,
		},
		{
			name: "fail case, error invalid cache value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("Delete", mock.Anything).Return(nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("dede:scan:1dc13ds5c1651", nil).Once()
				return moc
			},
			wantScore: datatypes.AtRiskResponse{},
			wantErr:   constants.InvalidKeyValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := RiskService{logger.ZapLogger{Logger: zap.NewExample()}, tc.redisClient(), nil}
			score, err := risk.DeleteCache("key")
			if tc.wantScore != score {
				t.Errorf("expected score %d got %d", tc.wantScore, score)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestExtendTTL(t *testing.T) {

	type tests struct {
		name        string
		log         logger.ZapLogger
		redisClient func() *mocks.RedisOps
		wantErr     error
	}

	testCases := []tests{
		{
			name: "valid case",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetKeys", mock.Anything).Return([]string{"key1", "key2"}, nil).Once()
				moc.On("SetTTL", mock.Anything, mock.AnythingOfType("int")).Return(nil).Twice()
				return moc
			},
			wantErr: nil,
		},
		{
			name: "fail case, error getting cache keys",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetKeys", mock.Anything).Return([]string{}, test.CacheGetKeysErr).Once()
				return moc
			},
			wantErr: test.CacheGetKeysErr,
		},
		{
			name: "fail case, error setting ttl for cache keys",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetKeys", mock.Anything).Return([]string{"key1"}, nil).Once()
				moc.On("SetTTL", mock.Anything, mock.AnythingOfType("int")).Return(test.CacheSetTTLErr).Once()
				return moc
			},
			wantErr: test.CacheSetTTLErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := RiskService{logger.ZapLogger{Logger: zap.NewExample()}, tc.redisClient(), nil}
			err := risk.ExtendTTL("email", 10)
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetEventScore(t *testing.T) {

	type tests struct {
		name        string
		log         logger.ZapLogger
		redisClient func() *mocks.RedisOps
		timestamp   string
		wantKey     string
		wantValue   string
		wantScore   int
		wantErr     error
	}

	testCases := []tests{
		{
			name: "valid case without mid",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("GetValue", mock.Anything).Return("78:email:5gf8d54ss45s8", nil).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "email:1684231487",
			wantValue: "78:email:5gf8d54ss45s8",
			wantScore: 78,
			wantErr:   nil,
		},
		{
			name: "valid case with mid",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1", "key2"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("46:scan:<<mid", nil).Once()

				return moc
			},
			timestamp: "1684231487",
			wantKey:   "key1",
			wantValue: "46:scan:<<mid",
			wantScore: 46,
			wantErr:   nil,
		},
		{
			name: "fail case, error checking if key exists in cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, test.CacheKeyExistsErr).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   test.CacheKeyExistsErr,
		},
		{
			name: "fail case, invalid timestamp value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, nil).Once()
				return moc
			},
			timestamp: "16fs54dfs5d46",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   constants.InvalidTimestampValue,
		},
		{
			name: "fail case, error getting cache value",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", test.CacheGetValueErr).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   test.CacheGetValueErr,
		},
		{
			name: "fail case, error getting keys from cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{}, test.CacheGetKeysErr).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   test.CacheGetKeysErr,
		},
		{
			name: "fail case, error getting value from cache",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(false, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", nil).Once()
				moc.On("GetKeys", mock.Anything).Return([]string{"key1", "key2"}, nil).Once()
				moc.On("GetValue", mock.Anything).Return("", test.CacheGetValueErr).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   test.CacheGetValueErr,
		},
		{
			name: "fail case, error converting cache value to integer",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("Exists", mock.Anything).Return(true, nil).Once()
				moc.On("GetValue", mock.Anything).Return("fd1b5df:email:5gf8d54ss45s8", nil).Once()
				return moc
			},
			timestamp: "1684231487",
			wantKey:   "",
			wantValue: "",
			wantScore: 0,
			wantErr:   constants.InvalidScoreValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := RiskService{logger.ZapLogger{Logger: zap.NewExample()}, tc.redisClient(), nil}
			score, err := risk.GetEventScore("email", tc.timestamp, "<<mid")
			if tc.wantKey != score.AtRiskKey {
				t.Errorf("expected key %s got %s", tc.wantKey, score.AtRiskKey)
			}
			if tc.wantValue != score.AtRiskValue {
				t.Errorf("expected value %s got %s", tc.wantValue, score.AtRiskValue)
			}
			if tc.wantScore != score.AtRiskScore {
				t.Errorf("expected score %d got %d", tc.wantScore, score.AtRiskScore)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetScore(t *testing.T) {

	type tests struct {
		name           string
		getAtRiskScore func(email string) ([]datatypes.RiskScore, error)
		wantResp       []datatypes.RiskScore
		wantErr        error
	}

	testCases := []tests{
		{
			name: "valid case",
			getAtRiskScore: func(email string) ([]datatypes.RiskScore, error) {
				return []datatypes.RiskScore{
					{Email: "email1", SelfHarmScore: "65"},
					{Email: "email2", SelfHarmScore: "47"},
				}, nil
			},
			wantResp: []datatypes.RiskScore{
				{Email: "email1", SelfHarmScore: "65"},
				{Email: "email2", SelfHarmScore: "47"},
			},
			wantErr: nil,
		},
		{
			name: "fail case",
			getAtRiskScore: func(email string) ([]datatypes.RiskScore, error) {
				return nil, test.DBSomethingWentWrongErr
			},
			wantResp: nil,
			wantErr:  test.DBSomethingWentWrongErr,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := RiskService{logger.ZapLogger{Logger: zap.NewExample()}, nil, tc.getAtRiskScore}
			resp, err := risk.GetScore("key")
			if !assert.Equal(t, tc.wantResp, resp) {
				t.Errorf("expected resp %+v got %+v", tc.wantResp, resp)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}
