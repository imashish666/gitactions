package model

import (
	"www-api/internal/datatypes"
)

// GetAtRiskScore fetches user_email, self_harm_score from AtRiskScore table based on user_email
func (m *ReadModel) GetAtRiskScore(email string) ([]datatypes.RiskScore, error) {
	scores := []datatypes.RiskScore{}
	err := m.db.Select(GetAtRiskQuery, &scores, email)
	if err != nil {
		m.log.Error("error fetching self_harm_scores from atRiskScore table", map[string]interface{}{"error": err, "email": email, "query": GetAtRiskQuery})
		return nil, err
	}
	return scores, nil
}
