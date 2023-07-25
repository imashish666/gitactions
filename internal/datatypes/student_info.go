package datatypes

type StudentInfo struct {
	GivenName  string `db:"givenName"`
	FamilyName string `db:"familyName"`
}

type StudentInfoRequest struct {
	Email string `json:"email"`
	Fid   string `json:"fid"`
}
