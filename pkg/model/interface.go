package model

import (
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/database"
)

type ReadModel struct {
	log logger.ZapLogger
	db  database.DatabaseOps
}

type WriteModel struct {
	log logger.ZapLogger
	db  database.DatabaseOps
}

type DatabaseReadAction interface {
	GetAtRiskScore(email string) ([]datatypes.RiskScore, error)
	GetStudentInfo(email string) (datatypes.StudentInfo, error)
	GetStudentInfoWithFid(fid, email string) (datatypes.StudentInfo, error)
	GetAwareNotification(fid string) (datatypes.Notification, error)
	GetUserTimezone(email string) (string, error)
	GetFilterType(fid string) (datatypes.FilterType, error)
}

// NewReadModel returns an instance of ReadModel struct
func NewReadModel(log logger.ZapLogger, db database.Database) ReadModel {
	return ReadModel{
		log: log,
		db:  database.NewDatabase(db.DB),
	}
}

// NewWriteModel returns an instance of ReadModel struct
func NewWriteModel(log logger.ZapLogger, db database.Database) WriteModel {
	return WriteModel{
		log: log,
		db:  database.NewDatabase(db.DB),
	}
}
