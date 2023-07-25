package model

var GetAtRiskQuery = "SELECT user_email, self_harm_score FROM AtRiskScore WHERE user_email IN (?)"
var GetStudentInfoQuery = "SELECT givenName, familyName FROM usermap WHERE userEmail = ? UNION SELECT givenName, familyName FROM azureUsers WHERE userEmail = ?"
var GetStudentInfoWithFidQuery = "SELECT givenName, familyName FROM usermap WHERE email = ? AND userEmail = ? UNION SELECT givenName, familyName FROM azureUsers WHERE fid = ? AND userEmail = ?"
var GetTimeZone = "SELECT timezone FROM user WHERE email = ?"
var GetAwareNotification = "SELECT * FROM awareEmailNotification WHERE fid = ?"
var GetFilterType = "select s.* from setting as s left join user as u on s.user_id = u.userId where u.email = ? limit 1"
