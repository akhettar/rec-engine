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
	a.router.HandleFunc("/api/items/user/{user}", a.userItems).Methods(http.MethodGet)
	a.router.HandleFunc("/api/items", a.popularItems).Methods(http.MethodGet)
	a.router.HandleFunc("/api/probability/user/{user}/item/{item}", a.itemProbability).Methods(http.MethodGet)
	a.router.PathPrefix("/swagger/").Handler(httpswag.WrapHandler)
}

// @Summary Create rating for a gien user with an item
// @ID post-rate
// @Description Adds rating for a given user with an item
// @Produce json
// @Param body body model.Rate true "body"
// @Success 201 {object} model.Rate "Rating created"
// @Failure 400 {object} model.ErrResponse "Invalid payload"
// @Failure 500 {object} model.ErrResponse "Internal server error"
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

	log.Infof("User %s ranked item %s with %f", req.User, req.Item, req.Score)
	respondWithJSON(rw, http.StatusCreated, req)
}

// @Summary Get recommendations
// @ID get-recommendations
// @Description Gets recommendations for a given user
// @Produce json
// @Param user path string true "user ID"
// @Success 200 {object} model.Recommendations "Recommendation returned"
// @Failure 400 {object} model.ErrResponse "Invalid payload"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /api/recommendation/user/{user} [get]
func (a *App) recommend(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.WithFields(log.Fields{"User": user}).Info("Received request to retrienve suggestion for user")

	// 1. batch upddate DB
	if err := a.eng.BatchUpdateSimilarUsers(-1); err != nil {
		log.Warnf("failed to update DB with erro %s", err)
	}

	// 2. Update suggested items
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
// @Success 200 {object} model.ItemProbability "ItemProbability returned"
// @Failure 400 {object} model.ErrResponse "Invalid payload"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /api/probability/user/{user}/item/{item} [get]
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

// @Summary Get User Items
// @ID get-user-item
// @Description Gets user items
// @Produce json
// @Param user path string true "user ID"
// @Success 200 {object} model.Items "Items returned"
// @Failure 400 {object} model.ErrResponse "Invalid payload"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /api/items/user/{user} [get]
func (a *App) userItems(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.WithFields(log.Fields{"User": user}).Info("Received request to calculate item probability of given user")

	result, err := a.eng.GetUserItems(user, 10000)

	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	log.Infof("Got results: %v", result)
	respondWithJSON(rw, http.StatusOK, convertToIterms(user, result))
}

// @Summary Get most popular items
// @ID get-popular-items
// @Description Gets the most popular items
// @Produce json
// @Param size query string false "number of results size"
// @Success 200 {object} model.Items "Items returned"
// @Failure 400 {object} model.ErrResponse "Invalid payload"
// @Failure 500 {object} model.ErrResponse "Internal server error"
// @Router /api/items [get]
func (a *App) popularItems(rw http.ResponseWriter, r *http.Request) {

	log.Info("Received request to retrieve the most popular items")
	resultSize := -1
	if size, err := strconv.Atoi(r.URL.Query().Get("size")); err == nil {
		resultSize = size - 1
	}

	log.Infof("size got: %d", resultSize)
	results, err := a.eng.GetPopularItems(resultSize)
	if err != nil {
		respondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	log.Infof("Got results: %v", results)
	respondWithJSON(rw, http.StatusOK, convertToIterms("", results))
}
