package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/ianmcmahon/compliancebuddy/faaservices"
	"github.com/ianmcmahon/compliancebuddy/model"
	"gopkg.in/redis.v5"
)

var api *Api

func setup(t *testing.T) func() {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatal(err)
	}

	redisClient := redis.NewClient(&redis.Options{Addr: s.Addr()})
	faaClient := faaservices.NewClient()
	faaClient.RegistryFetcher = faaservices.LocalRegistryFetcher{"../faaservices/testdata/abbrev"}

	regService := model.NewRegistrationService(redisClient, faaClient)

	if err := regService.UpdateRegistrationData(); err != nil {
		t.Fatal(err)
	}

	api = New(redisClient, faaClient)

	return func() {
		s.Close()
	}
}

func assertAndUnpack(t *testing.T, rr *httptest.ResponseRecorder, status int, v interface{}) {
	if got := rr.Code; got != status {
		t.Errorf("incorrect status: expected %d, got %d -- body: %s", status, got, rr.Body.String())
	}

	if err := json.NewDecoder(rr.Body).Decode(v); err != nil {
		t.Errorf("error unmarshaling response: %v", err)
	}
}

func TestRegistrationService(t *testing.T) {
	defer setup(t)()

	rr := httptest.NewRecorder()
	router := api.Router()

	req, err := http.NewRequest("GET", "/registration/N7706Y", nil)
	if err != nil {
		t.Fatal(err)
	}

	router.ServeHTTP(rr, req)

	regData := faaservices.AircraftRegistration{}
	assertAndUnpack(t, rr, http.StatusOK, &regData)

	fmt.Printf("%v\n", regData)
}
