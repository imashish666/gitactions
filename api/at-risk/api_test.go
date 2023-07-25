package atRisk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"www-api/config"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/test"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewRisk(t *testing.T) {
	type tests struct {
		name        string
		log         logger.ZapLogger
		config      config.Config
		connections *datatypes.Connections
	}

	testCases := []tests{
		{
			name:   "valid case",
			log:    logger.ZapLogger{},
			config: config.Config{},
			connections: &datatypes.Connections{
				DB: map[string]*sqlx.DB{
					constants.SchoolsReadDBKey: &sqlx.DB{},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := NewRiskAPI(tc.config, tc.log, tc.connections)
			if riskService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, riskService.log)
			}
			if !assert.Equal(t, tc.config, riskService.config) {
				t.Errorf("expected config %v is different from actual config %v", tc.config, riskService.config)
			}
			if riskService.createCache == nil {
				t.Errorf("expected createCache but got nil")
			}
			if riskService.deleteCache == nil {
				t.Errorf("expected deleteCache but got nil")
			}
			if riskService.getScore == nil {
				t.Errorf("expected getScore but got nil")
			}
			if riskService.extentTTL == nil {
				t.Errorf("expected extentTTL but got nil")
			}
			if riskService.getEventScore == nil {
				t.Errorf("expected getEventScore but got nil")
			}
		})
	}
}

func TestCreateCache(t *testing.T) {
	type tests struct {
		name             string
		params           map[string]string
		body             map[string]interface{}
		createCache      func(key, value string) (datatypes.AtRiskResponse, error)
		expectedStatus   int
		expectedResponse string
	}

	// Convert the data to JSON

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"atRiskKey":   "some_key@securly.com:16546548465",
				"atRiskValue": "35:docs:4854184194",
			},
			createCache: func(key, value string) (datatypes.AtRiskResponse, error) {
				return datatypes.AtRiskResponse{10}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"totalAtRiskScore\":10}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"atRiskKey":   1,
				"atRiskValue": 2,
			},
			createCache:      nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CacheRequest.atRiskKey of type string\"}",
		},
		{
			name: "fail case, missing at_riskKey_param",
			body: map[string]interface{}{
				"atRiskValue": "35:docs:4854184194",
			},
			createCache:      nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"atRiskKey missing in request body\"}",
		},
		{
			name: "fail case, missing atRiskValue param",
			// params:           map[string]string{"atRiskKey": "some_key@securly.com:16546548465"},
			body: map[string]interface{}{
				"atRiskKey": "some_key@securly.com:16546548465",
			},
			createCache:      nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"atRiskValue missing in request body\"}",
		},
		{
			name: "fail case, error createCache func",
			// params: map[string]string{"atRiskKey": "some_key@securly.com:16546548465", "atRiskValue": "84:gmail:546515615"},
			body: map[string]interface{}{
				"atRiskKey":   "some_key@securly.com:16546548465",
				"atRiskValue": "35:docs:4854184194",
			},
			createCache: func(key, value string) (datatypes.AtRiskResponse, error) {
				return datatypes.AtRiskResponse{0}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := RiskAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, tc.createCache, nil, nil, nil, nil}
			u, err := url.Parse("")
			assert.NoError(t, err)

			q := u.Query()
			for name, value := range tc.params {
				q.Set(name, value)
			}

			u.RawQuery = q.Encode()
			jsonData, err := json.Marshal(tc.body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			riskService.CreateCache(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestDeleteCache(t *testing.T) {
	type tests struct {
		name             string
		params           map[string]string
		body             map[string]interface{}
		deleteCache      func(key string) (datatypes.AtRiskResponse, error)
		expectedStatus   int
		expectedResponse string
	}
	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"atRiskKey": "some_key@securly.com:16546548465",
			},
			deleteCache: func(key string) (datatypes.AtRiskResponse, error) {
				return datatypes.AtRiskResponse{10}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"totalAtRiskScore\":10}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"atRiskKey": 1,
			},
			deleteCache:      nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CacheRequest.atRiskKey of type string\"}",
		},
		{
			name:             "fail case, missing atRiskKey in request body",
			body:             map[string]interface{}{},
			deleteCache:      nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"atRiskKey missing in request body\"}",
		},
		{
			name: "fail case, error resource not found",
			body: map[string]interface{}{
				"atRiskKey": "some_key@securly.com:16546548465",
			},
			deleteCache: func(key string) (datatypes.AtRiskResponse, error) {
				return datatypes.AtRiskResponse{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"key doesn't exists\"}",
		},
		{
			name: "fail case, error deleteCache func",
			body: map[string]interface{}{
				"atRiskKey": "some_key@securly.com:16546548465",
			},
			deleteCache: func(key string) (datatypes.AtRiskResponse, error) {
				return datatypes.AtRiskResponse{0}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := RiskAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, tc.deleteCache, nil, nil, nil}
			u, err := url.Parse("")
			assert.NoError(t, err)

			q := u.Query()
			for name, value := range tc.params {
				q.Set(name, value)
			}

			u.RawQuery = q.Encode()
			jsonData, err := json.Marshal(tc.body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			req, err := http.NewRequest("DELETE", u.String(), bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			riskService.DeleteCache(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestScore(t *testing.T) {
	type tests struct {
		name             string
		params           map[string]string
		body             map[string]interface{}
		getScore         func(email string) ([]datatypes.RiskScore, error)
		expectedStatus   int
		expectedResponse string
	}
	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{"userEmail": "some1@email.com"},
			getScore: func(email string) ([]datatypes.RiskScore, error) {
				return []datatypes.RiskScore{
					{Email: "some1@email.com", SelfHarmScore: "65"},
					{Email: "some1@email.com", SelfHarmScore: "16"},
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"atRiskScores\":[{\"Email\":\"some1@email.com\",\"SelfHarmScore\":\"65\"},{\"Email\":\"some1@email.com\",\"SelfHarmScore\":\"16\"}]}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{"userEmail": 1},
			getScore: func(email string) ([]datatypes.RiskScore, error) {
				return []datatypes.RiskScore{
					{Email: "some1@email.com", SelfHarmScore: "65"},
					{Email: "some1@email.com", SelfHarmScore: "16"},
				}, nil
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field AtRiskRequest.userEmail of type string\"}",
		},
		{
			name:             "fail case, missing email in request body",
			body:             map[string]interface{}{},
			getScore:         nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"email missing in request body\"}",
		},
		{
			name: "fail case, error getScore func",
			body: map[string]interface{}{"userEmail": "some1@email.com"},
			getScore: func(email string) ([]datatypes.RiskScore, error) {
				return nil, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := RiskAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, tc.getScore, nil, nil}
			u, err := url.Parse("")
			assert.NoError(t, err)

			q := u.Query()
			for name, value := range tc.params {
				q.Set(name, value)
			}
			u.RawQuery = q.Encode()

			jsonData, err := json.Marshal(tc.body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			req, err := http.NewRequest("GET", u.String(), bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			riskService.Score(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestExtendTTL(t *testing.T) {
	type tests struct {
		name             string
		params           map[string]string
		body             map[string]interface{}
		extentTTL        func(email string, ttl int) error
		expectedStatus   int
		expectedResponse string
	}
	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{"userEmail": "some1@email.com"},
			extentTTL: func(email string, ttl int) error {
				return nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"message\":\"ttl extended\"}",
		},
		{
			name:             "invalid request body",
			body:             map[string]interface{}{"userEmail": 1},
			extentTTL:        nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field AtRiskRequest.userEmail of type string\"}",
		},
		{
			name:             "fail case, missing email in request body",
			body:             map[string]interface{}{},
			extentTTL:        nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"email missing in request body\"}",
		},
		{
			name:             "fail case, invalid ttl in request body",
			body:             map[string]interface{}{"userEmail": "some1@email.com", "ttl": "invalid"},
			extentTTL:        nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"invalid ttl in request body, should be numeric\"}",
		},
		{
			name: "fail case, error extentTTL func",
			body: map[string]interface{}{"userEmail": "some1@email.com"},
			extentTTL: func(email string, ttl int) error {
				return test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := RiskAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, nil, tc.extentTTL, nil}
			u, err := url.Parse("")
			assert.NoError(t, err)

			q := u.Query()
			for name, value := range tc.params {
				q.Set(name, value)
			}
			u.RawQuery = q.Encode()

			jsonData, err := json.Marshal(tc.body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			req, err := http.NewRequest("PUT", u.String(), bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			riskService.ExtendTTL(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestEventScore(t *testing.T) {
	type tests struct {
		name             string
		params           map[string]string
		body             map[string]interface{}
		getEventScore    func(email string, timestamp string, mid string) (datatypes.EventScoreResponse, error)
		expectedStatus   int
		expectedResponse string
	}
	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{"userEmail": "some_email@securly.com", "timestamp": "1684323604", "mid": "<<somemid"},
			getEventScore: func(email string, timestamp string, mid string) (datatypes.EventScoreResponse, error) {
				return datatypes.EventScoreResponse{"key", "value", 45}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"atRiskKey\":\"key\",\"atRiskScore\":45,\"atRiskValue\":\"value\"}",
		},
		{
			name:             "invalid request body",
			body:             map[string]interface{}{"userEmail": 1, "timestamp": "1684323604", "mid": "<<somemid"},
			getEventScore:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field AtRiskRequest.userEmail of type string\"}",
		},
		{
			name:             "fail case, missing userEmail",
			body:             map[string]interface{}{"timestamp": "1684323604"},
			getEventScore:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"email missing in request body\"}",
		},
		{
			name:             "fail case, missing timestamp",
			body:             map[string]interface{}{"userEmail": "some_email@securly.com"},
			getEventScore:    nil,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"timestamp missing in request body\"}",
		},
		{
			name: "fail case, error resource not found",
			body: map[string]interface{}{"userEmail": "some_email@securly.com", "timestamp": "1684323604", "mid": "<<somemid"},
			getEventScore: func(email string, timestamp string, mid string) (datatypes.EventScoreResponse, error) {
				return datatypes.EventScoreResponse{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"key doesn't exists\"}",
		},
		{
			name: "fail case, error extentTTL func",
			body: map[string]interface{}{"userEmail": "some_email@securly.com", "timestamp": "1684323604", "mid": "<<somemid"},
			getEventScore: func(email string, timestamp string, mid string) (datatypes.EventScoreResponse, error) {
				return datatypes.EventScoreResponse{}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskService := RiskAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, nil, nil, tc.getEventScore}
			u, err := url.Parse("")
			assert.NoError(t, err)

			q := u.Query()
			for name, value := range tc.params {
				q.Set(name, value)
			}
			u.RawQuery = q.Encode()

			jsonData, err := json.Marshal(tc.body)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			req, err := http.NewRequest("GET", u.String(), bytes.NewBuffer(jsonData))
			assert.NoError(t, err)

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			riskService.EventScore(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}
