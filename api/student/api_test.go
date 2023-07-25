package student

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

func TestNewInfoAPI(t *testing.T) {
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
			customerService := NewInfoAPI(tc.config, tc.log, tc.connections)
			if customerService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, customerService.log)
			}
			if !assert.Equal(t, tc.config, customerService.config) {
				t.Errorf("expected config %v is different from actual config %v", tc.config, customerService.config)
			}
			if customerService.getStudentInfo == nil {
				t.Errorf("expected getStudentInfo but got nil")
			}
		})
	}
}

func TestGetInfo(t *testing.T) {
	type tests struct {
		name             string
		header           map[string]string
		params           map[string]string
		body             map[string]interface{}
		getStudentInfo   func(fid, email string) (datatypes.StudentInfo, error)
		expectedStatus   int
		expectedResponse string
	}

	testCases := []tests{
		{
			name: "valid case",
			body: map[string]interface{}{
				"email": "some_key@securly.com",
			},
			getStudentInfo: func(fid, email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"}, nil
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: "{\"GivenName\":\"given_name\",\"FamilyName\":\"family_name\"}",
		},
		{
			name: "invalid request body",
			body: map[string]interface{}{
				"email": 1,
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"json: cannot unmarshal number into Go struct field StudentInfoRequest.email of type string\"}",
		},
		{
			name:             "fail case, missing email in request body",
			body:             map[string]interface{}{},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"email missing in request body\"}",
		},
		{
			name: "fail case, resource not found",
			body: map[string]interface{}{
				"email": "some_key@securly.com",
			},
			getStudentInfo: func(fid, email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{}, constants.ResourceNotFound
			},
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "{\"message\":\"user some_key@securly.com doesn't exists\"}",
		},
		{
			name: "fail case, error getStudentInfo func",
			body: map[string]interface{}{
				"email": "some_key@securly.com",
			},
			getStudentInfo: func(fid, email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{}, test.InternalServerErr
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "{\"message\":\"internal server error\"}",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			custService := InfoAPI{config.Config{}, logger.ZapLogger{Logger: zap.NewExample()}, tc.getStudentInfo}
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
			custService.GetInfo(c)

			// Assert the expected response
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}
