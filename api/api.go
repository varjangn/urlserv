package api

import (
	"log"
	"net/http"

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
	router := mux.NewRouter()

	router.Use(LoggingMiddleware)

	router.HandleFunc(v1Prefix,
		Method(func(w http.ResponseWriter, r *http.Request) {
			WriteJSON(w, http.StatusOK, map[string]string{
				"status": "API is running",
				"v":      "v1"})
		}, "GET"))

	router.HandleFunc(v1Prefix+"users/register/",
		Method(s.Register, "POST"))

	router.HandleFunc(v1Prefix+"users/login/",
		Method(s.Login, "POST"))

	router.HandleFunc(v1Prefix+"users/profile/",
		Method(JWTAuth(s.Profile, s.store), "GET"))

	router.HandleFunc(v1Prefix+"users/shortner/",
		Method(JWTAuth(s.Shortner, s.store), "POST"))

	router.HandleFunc(v1Prefix+"users/urls/",
		Method(JWTAuth(s.GetUsersURLs, s.store), "GET"))

	log.Println("APIServer running on", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}
