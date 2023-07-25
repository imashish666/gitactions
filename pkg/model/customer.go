package model

import (
	"www-api/internal/constants"
	"www-api/internal/datatypes"
)

// GetUserTimezone fetches givenName, familyName from usermap and azureUsers table based on userEmail
func (m ReadModel) GetUserTimezone(email string) (string, error) {
	var timezone []string
	err := m.db.Select(GetTimeZone, &timezone, email)
	if err != nil {
		m.log.Error("error fetching student info", map[string]interface{}{"error": err, "fid": email})
		return "", err
	}
	if len(timezone) == 0 {
		return "", constants.ResourceNotFound
	}
	return timezone[0], nil
}

// GetUserTimezone fetches givenName, familyName from usermap and azureUsers table based on userEmail
func (m ReadModel) GetAwareNotification(fid string) (datatypes.Notification, error) {
	var notification datatypes.Notification
	err := m.db.Get(GetAwareNotification, &notification, fid)
	if err != nil {
		m.log.Error("error fetching student info", map[string]interface{}{"error": err, "fid": fid})
		return notification, err
	}
	if notification.NotificationEmail == "" {
		return notification, constants.ResourceNotFound
	}
	return notification, nil
}

// GetUserTimezone fetches givenName, familyName from usermap and azureUsers table based on userEmail
func (m ReadModel) GetFilterType(fid string) (datatypes.FilterType, error) {
	var filter datatypes.FilterType
	err := m.db.Get(GetFilterType, &filter, fid)
	if err != nil {
		m.log.Error("error fetching student info", map[string]interface{}{"error": err, "fid": fid})
		return filter, err
	}
	if filter.UserID == 0 {
		return filter, constants.ResourceNotFound
	}
	return filter, nil
}
