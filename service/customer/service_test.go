package student

import (
	"testing"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/cache/mocks"
	"www-api/test"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewCustomerService(t *testing.T) {
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
				Redis: map[string]*redis.Client{
					constants.WWWReadRedisKey:  &redis.Client{},
					constants.WWWWriteRedisKey: &redis.Client{},
				},
				Elastic: nil,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := NewCustomerService(tc.log, tc.connections)
			if riskService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, riskService.log)
			}
			if riskService.redis == nil {
				t.Errorf("expected redis interface but got nil")
			}
			if riskService.getTimezoneFromUser == nil {
				t.Errorf("expected getTimezoneFromUser but got nil")
			}
			if riskService.getNotification == nil {
				t.Errorf("expected getNotification but got nil")
			}
			if riskService.getFilter == nil {
				t.Errorf("expected getFilter but got nil")
			}
		})
	}
}

func TestProuctPrivacyStatus(t *testing.T) {

	type tests struct {
		name        string
		fid         string
		redisClient func() *mocks.RedisOps
		want        map[string]int
		wantErr     error
	}

	testCases := []tests{
		{
			name: "valid case",
			fid:  "checkemail@rtqa1securly.com",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetValue", mock.Anything).Return("1", nil).Once()
				return moc
			},
			want:    map[string]int{"24": 0, "Aware": 0, "Filter": 1, "Responder": 0, "suppBully": 0},
			wantErr: nil,
		},
		{
			name: "invalid case, invalid fid",
			fid:  "checkemail@",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				// moc.On("GetValue", mock.Anything).Return("1", nil).Once()
				return moc
			},
			want:    map[string]int{"24": 0, "Aware": 0, "Filter": 0, "Responder": 0, "suppBully": 0},
			wantErr: constants.EmptyFid,
		},
		{
			name: "invalid case, redis call error out",
			fid:  "checkemail@rtqa1securly.com",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetValue", mock.Anything).Return("", test.CacheGetValueErr).Once()
				return moc
			},
			want:    map[string]int{"24": 0, "Aware": 0, "Filter": 0, "Responder": 0, "suppBully": 0},
			wantErr: test.CacheGetValueErr,
		},
		{
			name: "invalid case, string to int conversion failed",
			fid:  "checkemail@rtqa1securly.com",
			redisClient: func() *mocks.RedisOps {
				moc := mocks.NewRedisOps(t)
				moc.On("GetValue", mock.Anything).Return("invalid", nil).Once()
				return moc
			},
			want:    map[string]int{"24": 0, "Aware": 0, "Filter": 0, "Responder": 0, "suppBully": 0},
			wantErr: constants.InvalidCoversionToInt,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cust := CustomerService{logger.ZapLogger{Logger: zap.NewExample()}, tc.redisClient(), nil, nil, nil}
			status, err := cust.ProuctPrivacyStatus(tc.fid)
			if !assert.Equal(t, tc.want, status) {
				t.Errorf("expected status %v got %v", tc.want, status)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestTimezone(t *testing.T) {

	type tests struct {
		name                string
		getTimezoneFromUser func(fid string) (string, error)
		wantlocation        datatypes.TimezoneResponse
		wantErr             error
	}

	testCases := []tests{
		{
			name: "valid case",
			getTimezoneFromUser: func(fid string) (string, error) {
				return "Asia/Kolkata", nil
			},
			wantlocation: datatypes.TimezoneResponse{
				Tz:     "",
				TzAbbr: "Asia/Kolkata",
			},
			wantErr: nil,
		},
		{
			name: "invalid case, getTimezoneFromUser error out",
			getTimezoneFromUser: func(fid string) (string, error) {
				return "", test.InternalServerErr
			},
			wantErr: test.InternalServerErr,
		},
		{
			name: "valid case, invalid location",
			getTimezoneFromUser: func(fid string) (string, error) {
				return "Asia/invalid", nil
			},
			wantlocation: datatypes.TimezoneResponse{
				Tz:     "",
				TzAbbr: "Asia/Kolkata",
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cust := CustomerService{logger.ZapLogger{Logger: zap.NewExample()}, nil, tc.getTimezoneFromUser, nil, nil}
			location, err := cust.Timezone("")
			if tc.wantlocation != location {
				t.Errorf("expected response %v got %v", tc.wantlocation, location)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestNotification(t *testing.T) {
	type tests struct {
		name            string
		getNotification func(fid string) (datatypes.Notification, error)
		want            datatypes.Notification
		wantErr         error
	}

	testCases := []tests{
		{
			name: "valid case",
			getNotification: func(fid string) (datatypes.Notification, error) {
				return datatypes.Notification{ID: 1, Fid: "fid", NotificationEmail: "notification_email", Basegen: 345}, nil
			},
			want:    datatypes.Notification{ID: 1, Fid: "fid", NotificationEmail: "notification_email", Basegen: 345},
			wantErr: nil,
		},
		{
			name: "invalid case, getNotification error out",
			getNotification: func(fid string) (datatypes.Notification, error) {
				return datatypes.Notification{}, test.InternalServerErr
			},
			wantErr: test.InternalServerErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cust := CustomerService{logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, tc.getNotification, nil}
			notification, err := cust.Notification("")
			if !assert.Equal(t, tc.want, notification) {
				t.Errorf("expected notification %v got %v", tc.want, notification)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetFilterType(t *testing.T) {
	type tests struct {
		name      string
		getFilter func(fid string) (datatypes.FilterType, error)
		want      datatypes.FilterType
		wantErr   error
	}

	testCases := []tests{
		{
			name: "valid case",
			getFilter: func(fid string) (datatypes.FilterType, error) {
				return datatypes.FilterType{
					ID:                 1,
					BlockPageMsg:       []byte("some_key@securly.com"),
					UserID:             3,
					AdIntranet:         []byte(""),
					SchoolType:         0,
					ParentSetting:      1,
					LockValue:          1,
					ShowPnp:            0,
					ShowPause:          1,
					ShowEns:            0,
					AzureGrpImportPref: "azure_grp_import_pref",
				}, nil
			},
			want: datatypes.FilterType{
				ID:                 1,
				BlockPageMsg:       []byte("some_key@securly.com"),
				UserID:             3,
				AdIntranet:         []byte(""),
				SchoolType:         0,
				ParentSetting:      1,
				LockValue:          1,
				ShowPnp:            0,
				ShowPause:          1,
				ShowEns:            0,
				AzureGrpImportPref: "azure_grp_import_pref",
			},
			wantErr: nil,
		},
		{
			name: "invalid case, getFilter error out",
			getFilter: func(fid string) (datatypes.FilterType, error) {
				return datatypes.FilterType{}, test.InternalServerErr
			},
			wantErr: test.InternalServerErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cust := CustomerService{logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, nil, tc.getFilter}
			notification, err := cust.GetFilterType("")
			if !assert.Equal(t, tc.want, notification) {
				t.Errorf("expected notification %v got %v", tc.want, notification)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}
