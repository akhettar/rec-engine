package model


// Rate request
type Rate struct {
	User  string  `json:"user"`
	Item  string  `json:"item`
	Score float64 `json:"score`
}

// Suggestion response
type Suggestion struct{
	Item string `json:"item"`
	Score float64 `json:"score"`
}

// ItemProbability item probability
type ItemProbability struct {
	User string `json:"user"`
	Item string `json:"item"`
	Probability float64 `json:"propability"`
}