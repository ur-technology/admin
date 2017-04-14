package user

import (
	"encoding/json"
	"net/http"
	"ur/api"
	"ur/api/middleware"

	"goji.io/pat"

	goji "goji.io"
)

type apiV1 struct {
	parentMux    *goji.Mux
	modelFactory ModelFactory
}

func NewAPIV1(parentMux *goji.Mux, mf ModelFactory) api.API {
	return &apiV1{
		parentMux:    parentMux,
		modelFactory: mf,
	}
}

func (a *apiV1) Register() error {

	a.parentMux.Use(middleware.ContentTypeJSON)

	// Get single record
	a.parentMux.HandleFunc(pat.Get("/:id"), a.fetchUser)

	return nil
}

func (a *apiV1) fetchUser(w http.ResponseWriter, r *http.Request) {

	name := pat.Param(r, "id")
	user, err := a.modelFactory.FetchOne(name)

	if err != nil {
		api.JSONError(w, err)
		return
	}

	jsn, err := json.Marshal(user)

	if err != nil {
		api.JSONError(w, err)
		return
	}

	w.Write(jsn)
}
