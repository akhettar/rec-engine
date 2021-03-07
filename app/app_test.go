package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/akhettar/rec-engine/model"
)

func TestApp_rate(t *testing.T) {

	// 1. Create recorder and request
	rw := httptest.NewRecorder()

	rate := model.Rate{User: "test1", Item: "I1", Score: 10}

	body, _ := json.Marshal(rate)

	req := httptest.NewRequest(http.MethodPost, "/api/rate", bytes.NewBuffer(body))

	// 2. Server the request
	app.router.ServeHTTP(rw, req)

	// 3. Assert the status code and the body
	if rw.Result().StatusCode != http.StatusCreated {
		t.Errorf("server responded with the wrong error code: %d", rw.Result().StatusCode)
	}

	// 4. Assert the payload
	var res model.Rate
	if err := json.NewDecoder(rw.Body).Decode(&res); err != nil {
		t.Logf("failed to decode the response")
	}

	if res.User != rate.User {
		t.Errorf("got user id different from the requet: %v", res.User)
	}
}

func TestApp_recommend(t *testing.T) {

	// 1. Create recorder and request
	rw := httptest.NewRecorder()

	rate1 := model.Rate{User: "test1", Item: "I1", Score: 10}
	rate2 := model.Rate{User: "test2", Item: "I1", Score: 8}
	rate4 := model.Rate{User: "test2", Item: "I3", Score: 8}
	rate3 := model.Rate{User: "test2", Item: "I2", Score: 5}
	rate5 := model.Rate{User: "test1", Item: "I4", Score: 6}

	rates := []model.Rate{rate1, rate2, rate3, rate4, rate5}

	// 2. Add all the rates
	for _, rate := range rates {

		body, _ := json.Marshal(rate)

		req := httptest.NewRequest(http.MethodPost, "/api/rate", bytes.NewBuffer(body))

		// 2. Server the request
		app.router.ServeHTTP(rw, req)

		// 3. Assert the status code and the body
		if rw.Result().StatusCode != http.StatusCreated {
			t.Errorf("server responded with the wrong error code: %d", rw.Result().StatusCode)
		}

		// 4. Assert the payload
		var res model.Rate
		if err := json.NewDecoder(rw.Body).Decode(&res); err != nil {
			t.Logf("failed to decode the response")
		}

		if res.User != rate.User {
			t.Errorf("got user id different from the requet: %v", res.User)
		}
	}

	// 3. Query recommendations for user1
	recReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/recommendation/user/%s", rate1.User), nil)
	app.router.ServeHTTP(rw, recReq)

	var res model.Recommendations
	if err := json.NewDecoder(rw.Body).Decode(&res); err != nil {
		t.Errorf("failed to deserialise the response with error %v", err)
	}

	rec1 := model.Recommendation{Item: "I3", Score: 8}
	rec2 := model.Recommendation{Item: "I2", Score: 5}

	if len(res.Data) > 2 {
		t.Errorf("We got more recommendations than expected, got %d", len(res.Data))
	}

	if !reflect.DeepEqual(rec1, res.Data[0]) {
		t.Errorf("got wrong recommendation [got]: %v but [expected]: %v", res.Data[0], rec1)
	}
	if !reflect.DeepEqual(rec2, res.Data[1]) {
		t.Errorf("got wrong recommendation [got]: %v but [expected]: %v", res.Data[1], rec2)
	}

	// 3. Query recommendations for user2
	recReq = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/recommendation/user/%s", rate2.User), nil)
	app.router.ServeHTTP(rw, recReq)

	if err := json.NewDecoder(rw.Body).Decode(&res); err != nil {
		t.Errorf("failed to deserialise the response with error %v", err)
	}

	rec := model.Recommendation{Item: "I4", Score: 6}

	if len(res.Data) > 1 {
		t.Errorf("We got more recommendations than expected, got %d", len(res.Data))
	}

	if !reflect.DeepEqual(rec, res.Data[0]) {
		t.Errorf("got wrong recommendation [got]: %v but [expected]: %v", res.Data[0], rec)
	}

}
