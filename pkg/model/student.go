package model

import (
	"www-api/internal/constants"
	"www-api/internal/datatypes"
)

// GetStudentInfo fetches givenName, familyName from usermap and azureUsers table based on userEmail
func (m ReadModel) GetStudentInfo(email string) (datatypes.StudentInfo, error) {
	info := []datatypes.StudentInfo{}
	err := m.db.Select(GetStudentInfoQuery, &info, email, email)
	if err != nil {
		m.log.Error("error fetching student info", map[string]interface{}{"error": err, "email": email, "query": GetStudentInfoQuery})
		return datatypes.StudentInfo{}, err
	}

	if len(info) == 0 {
		return datatypes.StudentInfo{}, constants.ResourceNotFound
	}
	return info[0], nil
}

// GetStudentInfoWithFid fetches givenName, familyName from usermap and azureUsers table based on userEmail and fid
func (m ReadModel) GetStudentInfoWithFid(fid, email string) (datatypes.StudentInfo, error) {
	info := []datatypes.StudentInfo{}
	err := m.db.Select(GetStudentInfoWithFidQuery, &info, fid, email, fid, email)
	if err != nil {
		m.log.Error("error fetching student info", map[string]interface{}{"error": err, "fid": fid, "email": email, "query": GetStudentInfoWithFidQuery})
		return datatypes.StudentInfo{}, err
	}

	if len(info) == 0 {
		return datatypes.StudentInfo{}, constants.ResourceNotFound
	}
	return info[0], nil
}
