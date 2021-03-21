package app

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/akhettar/rec-engine/model"
	m "github.com/akhettar/rec-engine/model"
)

// respondWithError return json error
func respondWithError(rw http.ResponseWriter, code int, msg string) {
	respondWithJSON(rw, 400, m.ErrResponse{Error: msg, Code: code})
}

// respondWithJSON returns json response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	json, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

// convertToRecommendations convert to an array of Suggestion
func convertToRecommendations(user string, results []string) m.Recommendations {
	var item string
	recs := []model.Recommendation{}
	for index, res := range results {
		if index%2 == 0 {
			item = res
		} else {
			score, _ := strconv.ParseFloat(res, 64)
			recs = append(recs, m.Recommendation{Item: item, Score: score})
		}
	}
	return m.Recommendations{User: user, Data: recs}
}

// convertToRecommendations convert to an array of Suggestion
func convertToIterms(user string, results []string) interface{} {
	var item string
	items := []model.Item{}
	for index, res := range results {
		if index%2 == 0 {
			item = res
		} else {
			score, _ := strconv.ParseFloat(res, 64)
			items = append(items, m.Item{Name: item, Score: score})
		}
	}
	if user == "" {
		return m.PopularItems{Data: items}
	}
	return m.Items{User: user, Data: items}
}

// convertToItemProbability convert an array of result to ItemProbability
func convertToItemProbability(results []string) m.ItemProbability {
	propability, _ := strconv.ParseFloat(results[2], 64)
	return m.ItemProbability{User: results[0], Item: results[1], Probability: propability}
}
