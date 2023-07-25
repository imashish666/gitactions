package datatypes

// type StudentInfo struct {
// 	GivenName  string `db:"givenName"`
// 	FamilyName string `db:"familyName"`
// }

type CustomerRequest struct {
	Email string `json:"email"`
	Fid   string `json:"fid"`
}

type Notification struct {
	ID                int    `db:"id"`
	Fid               string `db:"fid"`
	NotificationEmail string `db:"notifEmail"`
	Basegen           int    `db:"basegen"`
}

type FilterType struct {
	ID                 int    `db:"id"`
	BlockPageMsg       []byte `db:"block_page_msg"`
	UserID             int    `db:"user_id"`
	AdIntranet         []byte `db:"ad_intranet"`
	SchoolType         int    `db:"schoolType"`
	ParentSetting      int    `db:"parentSetting"`
	LockValue          int    `db:"lockValue"`
	ShowPnp            int    `db:"showPnp"`
	ShowPause          int    `db:"showPause"`
	ShowEns            int    `db:"showENS"`
	AzureGrpImportPref string `db:"azureGrpImportPref"`
}

type TimezoneResponse struct {
	Tz     string
	TzAbbr string
}
