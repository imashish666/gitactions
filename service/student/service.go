package student

import (
	"www-api/internal/constants"
	"www-api/internal/datatypes"
	"www-api/internal/logger"
	"www-api/pkg/cache"
	"www-api/pkg/database"
	"www-api/pkg/model"
)

type StudentService struct {
	log                   logger.ZapLogger
	redis                 cache.RedisOps
	getStudentInfo        func(email string) (datatypes.StudentInfo, error)
	getStudentInfoWithFid func(fid, email string) (datatypes.StudentInfo, error)
}

// NewStudentService returns an instance of StudentService struct
func NewStudentService(log logger.ZapLogger, connections *datatypes.Connections) StudentService {
	readinterface := model.NewReadModel(log, database.NewDatabase(connections.DB[constants.SchoolsReadDBKey]))
	// writeinterface := model.NewWriteModel(log, database.NewDatabase(writeconn))

	return StudentService{
		log:                   log,
		getStudentInfo:        readinterface.GetStudentInfo,
		getStudentInfoWithFid: readinterface.GetStudentInfoWithFid,
	}
}

// StudentInfo gets info of a student based on email and fid (if available)
func (s StudentService) StudentInfo(fid, email string) (datatypes.StudentInfo, error) {
	s.log.Info("fetching student info", map[string]interface{}{"email": email, "fid": fid})

	if fid == "" {
		info, err := s.getStudentInfo(email)
		if err != nil {
			s.log.Error("error occured while fetching student info", map[string]interface{}{"error": err, "email": email})
			return info, err
		}
		return info, nil
	}

	info, err := s.getStudentInfoWithFid(fid, email)
	if err != nil {
		s.log.Error("error occured while fetching student info", map[string]interface{}{"error": err, "fid": fid, "email": email})
		return info, err
	}

	return info, nil
}
