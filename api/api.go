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

	router.HandleFunc(v1Prefix, func(w http.ResponseWriter, r *http.Request) {
		WriteJSON(w, http.StatusOK, map[string]string{"msg": "API is running"})
	})

	router.HandleFunc(v1Prefix+"register/", s.Register)
	router.HandleFunc(v1Prefix+"login/", s.Login)
	router.HandleFunc(v1Prefix+"profile/", JWTAuth(s.Profile, s.store))

	log.Println("APIServer running on", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, router)
}
