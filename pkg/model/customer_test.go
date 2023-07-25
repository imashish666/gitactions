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

func TestGetUserTimezone(t *testing.T) {
	type tests struct {
		name    string
		db      func() *mocks.DatabaseOps
		want    string
		wantErr error
	}
	var timezone []string
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &timezone, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*[]string)
					*arg = append(*arg, "Asia/Kolkata")
				}).Return(nil).Once()
				return moc
			},
			want:    "Asia/Kolkata",
			wantErr: nil,
		},
		{
			name: "invalid case, resource not found",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &timezone, mock.Anything, mock.Anything).Return(nil).Once()
				return moc
			},
			want:    "",
			wantErr: constants.ResourceNotFound,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &timezone, mock.Anything, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			want:    "",
			wantErr: test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			info, err := student.GetUserTimezone("email")
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected info %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetAwareNotification(t *testing.T) {
	type tests struct {
		name    string
		db      func() *mocks.DatabaseOps
		want    datatypes.Notification
		wantErr error
	}
	var noti datatypes.Notification
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*datatypes.Notification)
					arg.ID = 1
					arg.Fid = "fid"
					arg.NotificationEmail = "notification_email"
					arg.Basegen = 431
				}).Return(nil).Once()
				return moc
			},
			want: datatypes.Notification{
				ID: 1, Fid: "fid", NotificationEmail: "notification_email", Basegen: 431,
			},
			wantErr: nil,
		},
		{
			name: "invalid case, resource not found",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything).Return(nil).Once()
				return moc
			},
			want:    datatypes.Notification{},
			wantErr: constants.ResourceNotFound,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			want:    datatypes.Notification{},
			wantErr: test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			info, err := student.GetAwareNotification("email")
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected info %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}

func TestGetFilterType(t *testing.T) {
	type tests struct {
		name    string
		db      func() *mocks.DatabaseOps
		want    datatypes.FilterType
		wantErr error
	}
	var noti datatypes.FilterType
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*datatypes.FilterType)
					arg.ID = 1
					arg.BlockPageMsg = []byte("page")
					arg.UserID = 11
					arg.AdIntranet = []byte("intranet")
					arg.SchoolType = 1
					arg.ParentSetting = 0
					arg.ShowPnp = 1
					arg.ShowPause = 0
					arg.ShowEns = 1
					arg.AzureGrpImportPref = "azure"
				}).Return(nil).Once()
				return moc
			},
			want: datatypes.FilterType{
				ID: 1, BlockPageMsg: []byte("page"), UserID: 11, AdIntranet: []byte("intranet"), SchoolType: 1,
				ParentSetting: 0, ShowPnp: 1, ShowPause: 0, ShowEns: 1, AzureGrpImportPref: "azure",
			},
			wantErr: nil,
		},
		{
			name: "invalid case, resource not found",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything).Return(nil).Once()
				return moc
			},
			want:    datatypes.FilterType{},
			wantErr: constants.ResourceNotFound,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Get", mock.Anything, &noti, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			want:    datatypes.FilterType{},
			wantErr: test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			student := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			info, err := student.GetFilterType("email")
			if !assert.Equal(t, tc.want, info) {
				t.Errorf("expected info %v got %v", tc.want, info)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}
}
