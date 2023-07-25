package student

import (
	"testing"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/test"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewStudentService(t *testing.T) {
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
			riskService := NewStudentService(tc.log, tc.connections)
			if riskService.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, riskService.log)
			}
			if riskService.getStudentInfo == nil {
				t.Errorf("expected getStudentInfo but got nil")
			}
			if riskService.getStudentInfoWithFid == nil {
				t.Errorf("expected getStudentInfoWithFid but got nil")
			}
		})
	}
}

func TestStudentInfo(t *testing.T) {

	type tests struct {
		name                  string
		fid                   string
		email                 string
		getStudentInfo        func(email string) (datatypes.StudentInfo, error)
		getStudentInfoWithFid func(fid, email string) (datatypes.StudentInfo, error)
		want                  datatypes.StudentInfo
		wantErr               error
	}

	testCases := []tests{
		{
			name:  "valid case, fetch with email",
			email: "email",
			fid:   "",
			getStudentInfo: func(email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"}, nil
			},
			want:    datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"},
			wantErr: nil,
		},
		{
			name:  "valid case, fetch with fid and email",
			email: "email",
			fid:   "fid",
			getStudentInfoWithFid: func(fid, email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"}, nil
			},
			want:    datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"},
			wantErr: nil,
		},
		{
			name:  "invalid case, getStudentInfo error out",
			email: "email",
			fid:   "",
			getStudentInfo: func(email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{}, test.InternalServerErr
			},
			want:    datatypes.StudentInfo{},
			wantErr: test.InternalServerErr,
		},
		{
			name:  "invalid case, getStudentInfo error out",
			email: "email",
			fid:   "fid",
			getStudentInfoWithFid: func(fid, email string) (datatypes.StudentInfo, error) {
				return datatypes.StudentInfo{}, test.InternalServerErr
			},
			want:    datatypes.StudentInfo{},
			wantErr: test.InternalServerErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := StudentService{logger.ZapLogger{Logger: zap.NewExample()}, nil, tc.getStudentInfo, tc.getStudentInfoWithFid}
			info, err := risk.StudentInfo(tc.fid, tc.email)
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected %v got %v", tc.wantErr, err)
			}
		})
	}
}
