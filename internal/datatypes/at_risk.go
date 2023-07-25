package datatypes

type RiskScore struct {
	Email         string `db:"user_email"`
	SelfHarmScore string `db:"self_harm_score"`
}

type CacheRequest struct {
	AtRiskKey   string `json:"atRiskKey"`
	AtRiskValue string `json:"atRiskValue"`
}

type AtRiskRequest struct {
	UserEmail string `json:"userEmail"`
	TTL       string `json:"ttl"`
	Timestamp string `json:"timestamp"`
	Mid       string `json:"mid"`
}

type AtRiskResponse struct {
	AtRiskScore int
}

type EventScoreResponse struct {
	AtRiskKey   string
	AtRiskValue string
	AtRiskScore int
}
