package model

import (
	"testing"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/database"
	"www-api/pkg/database/mocks"
	"www-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestNewRiskModel(t *testing.T) {
	type tests struct {
		name string
		log  logger.ZapLogger
		db   database.Database
	}

	testCases := []tests{
		{
			name: "valid case",
			log:  logger.ZapLogger{},
			db:   database.Database{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			riskModel := NewReadModel(tc.log, tc.db)
			if riskModel.log != tc.log {
				t.Errorf("expected logger %v is different from actual logger %v", tc.log, riskModel.log)
			}
			if riskModel.db != tc.db {
				t.Errorf("expected db %v is different from actual db %v", tc.db, riskModel.db)
			}
		})
	}
}

func TestGetAtRiskScore(t *testing.T) {
	type tests struct {
		name      string
		db        func() *mocks.DatabaseOps
		wantScore []datatypes.RiskScore
		wantErr   error
	}
	scores := []datatypes.RiskScore{}
	testCases := []tests{
		{
			name: "valid case",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &scores, mock.Anything).Run(func(args mock.Arguments) {
					arg := args.Get(1).(*[]datatypes.RiskScore)
					*arg = append(*arg, datatypes.RiskScore{Email: "email1@securly.com", SelfHarmScore: "51"}, datatypes.RiskScore{Email: "email2@securly.com", SelfHarmScore: "56"})
				}).Return(nil).Once()
				return moc
			},
			wantScore: []datatypes.RiskScore{
				{Email: "email1@securly.com", SelfHarmScore: "51"},
				{Email: "email2@securly.com", SelfHarmScore: "56"},
			},
			wantErr: nil,
		},
		{
			name: "fail case, error select func",
			db: func() *mocks.DatabaseOps {
				moc := mocks.NewDatabaseOps(t)
				moc.On("Select", mock.Anything, &scores, mock.Anything).Return(test.DBSomethingWentWrongErr).Once()
				return moc
			},
			wantScore: nil,
			wantErr:   test.DBSomethingWentWrongErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			risk := ReadModel{logger.ZapLogger{Logger: zap.NewExample()}, tc.db()}
			score, err := risk.GetAtRiskScore("email")
			if !assert.Equal(t, tc.wantScore, score) {
				t.Errorf("expected score %v got %v", tc.wantScore, score)
			}
			if tc.wantErr != err {
				t.Errorf("expected error %v got %v", tc.wantErr, err)
			}
		})
	}

}
