package model

// Rate request
type Rate struct {
	User  string  `json:"user"`
	Item  string  `json:"item"`
	Score float64 `json:"score"`
}

// Recommendation response
type Recommendation struct {
	Item  string  `json:"item"`
	Score float64 `json:"score"`
}

// Recommendations type
type Recommendations struct {
	User string           `json:"user"`
	Data []Recommendation `json:"data"`
}

// Item user item
type Item struct {
	Name  string  `json:"item"`
	Score float64 `json:"score"`
}

// Items list of item
type Items struct {
	User string `json:"user"`
	Data []Item `json:"data"`
}

// ItemProbability item probability
type ItemProbability struct {
	User        string  `json:"user"`
	Item        string  `json:"item"`
	Probability float64 `json:"propability"`
}

// ErrResponse error response
type ErrResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
