package api

import (
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/varjangn/urlserv/storage"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() error {
	v1Prefix := "/api/v1/"
	r := mux.NewRouter()

	r.Use(LoggingMiddleware)

	r.HandleFunc(v1Prefix,
		Method(func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]string{
				"status": "API is running",
				"v":      "v1"})
		}, "GET"))

	r.HandleFunc(v1Prefix+"users/register/",
		Method(s.Register, "POST"))

	r.HandleFunc(v1Prefix+"users/login/",
		Method(s.Login, "POST"))

	r.HandleFunc(v1Prefix+"users/profile/",
		Method(JWTAuth(s.Profile, s.store), "GET"))

	r.HandleFunc(v1Prefix+"users/shortner/",
		Method(JWTAuth(s.Shortner, s.store), "POST"))

	r.HandleFunc(v1Prefix+"users/urls/",
		Method(JWTAuth(s.GetUsersURLs, s.store), "GET"))

	r.HandleFunc(v1Prefix+"users/urls/{id}/",
		Method(JWTAuth(s.HandleAURL, s.store), "GET", "DELETE", "PUT", "PATCH"))

	r.HandleFunc("/{id:["+regexp.QuoteMeta(`A-Za-z0-9_-`)+"]{7}}",
		Method(s.Redirect, "GET"))

	log.Println("APIServer running on", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, r)
}
