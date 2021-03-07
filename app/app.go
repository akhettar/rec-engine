package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/akhettar/rec-engine/docs"
	m "github.com/akhettar/rec-engine/model"
	"github.com/akhettar/rec-engine/redrec"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	httpswag "github.com/swaggo/http-swagger"
)

// App server instance type
type App struct {
	router *mux.Router
	eng    *redrec.Redrec
}

// InitialiseApp create New APP server
func InitialiseApp(redisURL string) *App {

	// 1. Create redis connection
	engine, err := redrec.New(redisURL)
	if err != nil {
		log.Fatalf("failed to intialise recommendation engine %s", err)
	}
	app := &App{router: mux.NewRouter(), eng: engine}
	app.initialiseRoutes()
	return app
}

// Run start the server
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.router))
}

func (a *App) initialiseRoutes() {
	a.router.HandleFunc("/api/rate", a.rate).Methods(http.MethodPost)
	a.router.HandleFunc("/api/recommendation/user/{user}", a.recommend).Methods(http.MethodGet)
	a.router.HandleFunc("/api/probability/user/{user}/item/{item}", a.itemProbability).Methods(http.MethodGet)
	a.router.PathPrefix("/swagger/").Handler(httpswag.WrapHandler)
}

// @Summary Create rating for a gien user with an item
// @ID post-rate
// @Description Adds rating for a given user with an item
// @Produce json
// @Param body body model.Rate true "body"
// @Success 201 {object} model.Rate "Rating created"
// @Failure 400 {object} model.ErrorMessage "Invalid payload"
// @Failure 500 {object} model.ErrorMessage "Internal server error"
// @Router /api/rate [post]
func (a *App) rate(rw http.ResponseWriter, r *http.Request) {
	var req m.Rate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(rw, http.StatusBadRequest, fmt.Sprintf("Failed to desrialise the payload: %s", err))
		return
	}
	if err := a.eng.Rate(req.Item, req.User, req.Score); err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// update DB
	if err := a.eng.BatchUpdateSimilarUsers(10000); err != nil {
		log.Warnf("failed to update DB with erro %s", err)
	}

	log.Infof("User %s ranked item %s with %f", req.User, req.Item, req.Score)
	respondWithJSON(rw, http.StatusCreated, req)
}

// @Summary Get suggestions
// @ID get-suggestions
// @Description Gets suggestions for a given user
// @Produce json
// @Param user path string true "user ID"
// @Success 200 {object} model.Suggestion "Suggestion returned"
// @Failure 400 {object} model.ErrorMessage "Invalid payload"
// @Failure 500 {object} model.ErrorMessage "Internal server error"
// @Router /api/recommendation/{user} [get]
func (a *App) recommend(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.WithFields(log.Fields{"User": user}).Info("Received request to retrienve suggestion for user")

	// 1. Update suggested items
	if err := a.eng.UpdateSuggestedItems(user, 10000); err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Get suggestions for a given user
	results, err := a.eng.GetUserSuggestions(user, 10000)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("Got results: %v", results)
	respondWithJSON(rw, http.StatusOK, convertToRecommendations(user, results))
}

// @Summary Get probability
// @ID get-probability
// @Description Gets probability for a given user and item
// @Produce json
// @Param user path string true "user ID"
// @Param item path string true "item ID"
// @Success 200 {object} model.Suggestion "Suggestion returned"
// @Failure 400 {object} model.ErrorMessage "Invalid payload"
// @Failure 500 {object} model.ErrorMessage "Internal server error"
// @Router /api/probability/{user}/{item} [get]
func (a *App) itemProbability(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	item := vars["item"]
	log.WithFields(log.Fields{"User": user, "Item": item}).Info("Received request to calculate item probability of given user")

	result, err := a.eng.CalcItemProbability(user, item)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	log.Infof("Got results: %v", m.ItemProbability{User: user, Item: item, Probability: result})
	respondWithJSON(rw, http.StatusOK, m.ItemProbability{User: user, Item: item, Probability: result})
}

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
func convertToRecommendations(user string, results []string) m.Recommendations{
	var item string
	recs := make([]m.Recommendation, len(results), len(results))
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

// convertToItemProbability convert an array of result to ItemProbability
func convertToItemProbability(results []string) m.ItemProbability {
	propability, _ := strconv.ParseFloat(results[2], 64)
	return m.ItemProbability{User: results[0], Item: results[1], Probability: propability}
}
