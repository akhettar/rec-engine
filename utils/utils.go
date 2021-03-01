package utils

import (
	"encoding/json"
	"net/http"
	"strconv"

	m "github.com/akhettar/rec-engine/model"
)

// RespondWithError return json error
func RespondWithError(rw http.ResponseWriter, code int, msg string) {
	RespondWithJSON(rw, 400, map[string]string{"error": msg})
}

// RespondWithJSON returns json response
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	json, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

// ConvertToSuggestion convert to an array of Suggestion
func ConvertToSuggestion(results []string) []m.Suggestion {
	var item string
	var suggestions []m.Suggestion
	for index, res := range results {
		if index%2 == 0 {
			item = res
		} else {
			score, _ := strconv.ParseFloat(res, 64)
			suggestions = append(suggestions, m.Suggestion{item, score})
		}
	}
	return suggestions
}

// ConvertToItemProbability convert an array of result to ItemProbability
func ConvertToItemProbability(results []string) m.ItemProbability {
	propability, _ := strconv.ParseFloat(results[2], 64)
	return m.ItemProbability{User: results[0], Item: results[1], Probability: propability}
}
