package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	m "github.com/akhettar/rec-engine/model"
	"github.com/akhettar/rec-engine/redrec"
	"github.com/akhettar/rec-engine/utils"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
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
	a.router.HandleFunc("/api/suggestion/{user}", a.suggest).Methods(http.MethodGet)
	a.router.HandleFunc("/api/probability/{user}/{item}", a.itemProbability).Methods(http.MethodGet)
}

func (a *App) rate(rw http.ResponseWriter, r *http.Request) {
	var req m.Rate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(rw, http.StatusBadRequest, fmt.Sprintf("Failed to desrialise the payload: %s", err))
		return
	}
	if err := a.eng.Rate(req.Item, req.User, req.Score); err != nil {
		utils.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// update DB
	if err := a.eng.BatchUpdateSimilarUsers(10000); err != nil {
		log.Warnf("failed to update DB with erro %s", err)
	}

	msg := fmt.Sprintf("User %s ranked item %s with %f", req.User, req.Item, req.Score)
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(msg))
}

func (a *App) suggest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	log.WithFields(log.Fields{"User": user}).Info("Received request to retrienve suggestion for user")

	// 1. Update suggested items
	if err := a.eng.UpdateSuggestedItems(user, 10000); err != nil {
		utils.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	// 2. Get suggestions for a given user
	results, err := a.eng.GetUserSuggestions(user, 10000)

	if err != nil {
		utils.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}
	log.Infof("Got results: %v", results)
	utils.RespondWithJSON(rw, http.StatusOK, utils.ConvertToSuggestion(results))
}

func (a *App) itemProbability(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["user"]
	item := vars["item"]
	log.WithFields(log.Fields{"User": user, "Item": item}).Info("Received request to calculate item probability of given user")

	result, err := a.eng.CalcItemProbability(user, item)

	if err != nil {
		utils.RespondWithError(rw, http.StatusInternalServerError, err.Error())
		return
	}

	log.Infof("Got results: %v", m.ItemProbability{user, item, result})
	utils.RespondWithJSON(rw, http.StatusOK, m.ItemProbability{user, item, result})
}
