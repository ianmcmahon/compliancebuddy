package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ianmcmahon/compliancebuddy/faaservices"
	"github.com/ianmcmahon/compliancebuddy/model"
	"gopkg.in/redis.v5"
)

type Api struct {
	redisClient *redis.Client
	faaClient   *faaservices.Client
}

func New(redis *redis.Client, faa *faaservices.Client) *Api {
	return &Api{
		redisClient: redis,
		faaClient:   faa,
	}
}

func (a *Api) Router() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/registration/{tailnum}", a.RegistrationSearch).Methods("GET")

	return r
}

////////////// Handlers //////////////

func (a *Api) RegistrationSearch(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tailnum, ok := vars["tailnum"]
	if !ok {
		apiError(rw, http.StatusNotFound, fmt.Errorf("tailnum not present"))
		return
	}

	reg := model.NewRegistrationService(a.redisClient, a.faaClient)

	regData, err := reg.GetRegistration(tailnum)
	if err != nil {
		apiError(rw, http.StatusNotFound, err)
		return
	}

	if err := json.NewEncoder(rw).Encode(regData); err != nil {
		apiError(rw, http.StatusInternalServerError, err)
		return
	}
}
