package model

import (
	"testing"
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/database/mocks"
	"www-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestGetStudentInfo(t *testing.T) {
	type tests struct {
		name    string
		db      func() *mocks.DatabaseOps
		want    datatypes.StudentInfo
		wantErr error
	}
	info := []datatypes.StudentInfo{}
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*[]datatypes.StudentInfo)
					*arg = append(*arg, datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"})
				}).Return(nil).Once()
				return moc
			},
			want: datatypes.StudentInfo{
				GivenName: "given_name", FamilyName: "family_name",
			},
			wantErr: nil,
		},
		{
			name: "invalid case, resource not found",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				return moc
			},
			want:    datatypes.StudentInfo{},
			wantErr: constants.ResourceNotFound,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			want:    datatypes.StudentInfo{},
			wantErr: test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			info, err := student.GetStudentInfo("email")
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected info %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetStudentInfoFid(t *testing.T) {
	type tests struct {
		name    string
		db      func() *mocks.DatabaseOps
		want    datatypes.StudentInfo
		wantErr error
	}
	info := []datatypes.StudentInfo{}
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*[]datatypes.StudentInfo)
					*arg = append(*arg, datatypes.StudentInfo{GivenName: "given_name", FamilyName: "family_name"})
				}).Return(nil).Once()
				return moc
			},
			want: datatypes.StudentInfo{
				GivenName: "given_name", FamilyName: "family_name",
			},
			wantErr: nil,
		},
		{
			name: "invalid case, resource not found",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
				return moc
			},
			want:    datatypes.StudentInfo{},
			wantErr: constants.ResourceNotFound,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &info, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			want:    datatypes.StudentInfo{},
			wantErr: test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			info, err := student.GetStudentInfoWithFid("email", "fid")
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected info %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}
