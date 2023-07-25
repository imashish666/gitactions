package customer

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

func TestNewCustomer(t *testing.T) {
	type tests struct {
		name        string
		log         logger.ZapLogger
		config      config.Config
		db          *sqlx.DB
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
			customerService := NewCustomerAPI(tc.config, tc.log, tc.connections)
			if customerService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, customerService.log)
			}
			if !assert.Equal(t, tc.config, customerService.config) {
				t.Errorf("expected config %v is different from actual config %v", tc.config, customerService.config)
			}
			if customerService.getPrivacyStatus == nil {
				t.Errorf("expected getPrivacyStatus but got nil")
			}
			if customerService.getTimezone == nil {
				t.Errorf("expected getTimezone but got nil")
			}
			if customerService.getNotificationEmail == nil {
				t.Errorf("expected getNotificationEmail but got nil")
			}
			if customerService.getFilterType == nil {
				t.Errorf("expected getFilterType but got nil")
			}
		})
	}
}

func TestPrivacyStatus(t *testing.T) {
	type tests struct {
		name             string
		header           map[string]string
		params           map[string]string
		body             map[string]interface{}
		getPrivacyStatus func(fid string) (map[string]int, error)
		expectedStatus   int
		expectedResponse string
	}

	// Convert the data to JSON

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getPrivacyStatus: func(fid string) (map[string]int, error) {
				return map[string]int{
					"24":        1,
					"Filter":    0,
					"Aware":     0,
					"Responder": 0,
					"suppBully": 0,
				}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"24\":1,\"Aware\":0,\"Filter\":0,\"Responder\":0,\"suppBully\":0}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"fid": 1,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CustomerRequest.fid of type string\"}",
		},
		{
			name:             "fail case, missing fid in request body",
			body:             map[string]interface{}{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"fid missing in request body\"}",
		},
		{
			name: "fail case, resource not found",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getPrivacyStatus: func(fid string) (map[string]int, error) {
				return map[string]int{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"user some_key@securly.com doesn't exists\"}",
		},
		{
			name: "fail case, error getPrivacyStatus func",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getPrivacyStatus: func(fid string) (map[string]int, error) {
				return map[string]int{}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custService := CustomerAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, tc.getPrivacyStatus, nil, nil, nil}
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
			for name, value := range tc.header {
				req.Header.Add(name, value)
			}

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			custService.PrivacyStatus(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestTimezone(t *testing.T) {
	type tests struct {
		name             string
		header           map[string]string
		params           map[string]string
		body             map[string]interface{}
		getTimezone      func(fid string) (datatypes.TimezoneResponse, error)
		expectedStatus   int
		expectedResponse string
	}

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getTimezone: func(fid string) (datatypes.TimezoneResponse, error) {
				return datatypes.TimezoneResponse{"Asia/Kolkata", "IST"}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"tz\":\"Asia/Kolkata\",\"tzAbbr\":\"IST\"}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"fid": 1,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CustomerRequest.fid of type string\"}",
		},
		{
			name:             "fail case, missing fid in request body",
			body:             map[string]interface{}{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"fid missing in request body\"}",
		},
		{
			name: "fail case, resource not found",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getTimezone: func(fid string) (datatypes.TimezoneResponse, error) {
				return datatypes.TimezoneResponse{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"user some_key@securly.com doesn't exists\"}",
		},
		{
			name: "fail case, error getTimezone func",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getTimezone: func(fid string) (datatypes.TimezoneResponse, error) {
				return datatypes.TimezoneResponse{}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custService := CustomerAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, tc.getTimezone, nil, nil}
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

			for name, value := range tc.header {
				req.Header.Add(name, value)
			}

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			custService.Timezone(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestNotification(t *testing.T) {
	type tests struct {
		name                 string
		header               map[string]string
		params               map[string]string
		body                 map[string]interface{}
		getNotificationEmail func(fid string) (datatypes.Notification, error)
		expectedStatus       int
		expectedResponse     string
	}

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getNotificationEmail: func(fid string) (datatypes.Notification, error) {
				return datatypes.Notification{ID: 1, Fid: "some_key@securly.com", NotificationEmail: "email", Basegen: 343}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"ID\":1,\"Fid\":\"some_key@securly.com\",\"NotificationEmail\":\"email\",\"Basegen\":343}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"fid": 1,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CustomerRequest.fid of type string\"}",
		},
		{
			name:             "fail case, missing fid in request body",
			body:             map[string]interface{}{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"fid missing in request body\"}",
		},
		{
			name: "fail case, resource not found",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getNotificationEmail: func(fid string) (datatypes.Notification, error) {
				return datatypes.Notification{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"user some_key@securly.com doesn't exists\"}",
		},
		{
			name: "fail case, error getNotificationEmail func",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getNotificationEmail: func(fid string) (datatypes.Notification, error) {
				return datatypes.Notification{}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custService := CustomerAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, tc.getNotificationEmail, nil}
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

			for name, value := range tc.header {
				req.Header.Add(name, value)
			}

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			custService.Notification(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}

func TestFilterType(t *testing.T) {
	type tests struct {
		name             string
		header           map[string]string
		params           map[string]string
		body             map[string]interface{}
		getFilterType    func(fid string) (string, error)
		expectedStatus   int
		expectedResponse string
	}

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getFilterType: func(fid string) (string, error) {
				return "", nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"ID\":1,\"BlockPageMsg\":\"c29tZV9rZXlAc2VjdXJseS5jb20=\",\"UserID\":3,\"AdIntranet\":\"\",\"SchoolType\":0,\"ParentSetting\":1,\"LockValue\":1,\"ShowPnp\":0,\"ShowPause\":1,\"ShowEns\":0,\"AzureGrpImportPref\":\"azure_grp_import_pref\"}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"fid": 1,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field CustomerRequest.fid of type string\"}",
		},
		{
			name:             "fail case, missing fid in request body",
			body:             map[string]interface{}{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"fid missing in request body\"}",
		},
		{
			name: "fail case, resource not found",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getFilterType: func(fid string) (string, error) {
				return "", constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"user some_key@securly.com doesn't exists\"}",
		},
		{
			name: "fail case, error getNotificationEmail func",
			body: map[string]interface{}{
				"fid": "some_key@securly.com",
			},
			getFilterType: func(fid string) (string, error) {
				return "", test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custService := CustomerAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, nil, nil, nil, tc.getFilterType}
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

			for name, value := range tc.header {
				req.Header.Add(name, value)
			}

			// Create a new recorder to capture the response
			recorder := httptest.NewRecorder()

			// Create a mock Gin context using the recorder and request
			c, _ := gin.CreateTestContext(recorder)
			c.Request = req

			// Call your handler function, passing in the mock context
			custService.FilterType(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}
