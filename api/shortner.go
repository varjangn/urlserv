package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/varjangn/urlserv/types"
)

type shortReqType struct {
	LongURL string `json:"long_url"`
}

func (s *APIServer) Shortner(w http.ResponseWriter, r *http.Request) {
	reqUser := r.Context().Value(authUserKey).(*types.User)
	reqBody := new(shortReqType)

	id, err := gonanoid.New(7)
	if err != nil {
		log.Println("ShortnerError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	}

	if err = json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		log.Println("ShortnerJSONError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "invalid request body"})
		return
	}

	url, err := s.store.GetURLbyLongURL(reqBody.LongURL)
	if err != nil {
		log.Println("ShortnerError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	} else if url != nil {
		WriteJSON(w, http.StatusCreated, url)
		return
	}

	longURL, err := s.store.GetLongURL(id)
	if err != nil {
		log.Println("ShortnerError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	} else if longURL != "" {
		log.Println("ShortnerDuplicateError: shortId=", id)
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	}

	url = types.NewURL(reqUser, id, reqBody.LongURL)
	if err = s.store.CreateURL(url); err != nil {
		log.Println("ShortnerError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	}
	WriteJSON(w, http.StatusCreated, url)
}

func (s *APIServer) GetUsersURLs(w http.ResponseWriter, r *http.Request) {
	reqUser := r.Context().Value(authUserKey).(*types.User)
	urls, err := s.store.GetURLs(reqUser)
	if err != nil {
		log.Println("ShortnerError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	}
	WriteJSON(w, http.StatusOK, urls)
}

func (s *APIServer) HandleAURL(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.GetUsersURL(w, r)
	} else if r.Method == "DELETE" {
		s.DeleteURL(w, r)
	} else if r.Method == "PATCH" || r.Method == "PUT" {
		s.UpdateLongURL(w, r)
	}
}

func (s *APIServer) GetUsersURL(w http.ResponseWriter, r *http.Request) {
	reqUser := r.Context().Value(authUserKey).(*types.User)
	vars := mux.Vars(r)
	urlId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest,
			Map{"error": "invalid url id"})
		return
	}
	url, err := s.store.GetURL(urlId, reqUser.Id)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "unknown error"})
		return
	}
	if url == nil {
		WriteJSON(w, http.StatusOK,
			Map{"error": "URL not found"})
		return
	}
	WriteJSON(w, http.StatusOK, url)
}

func (s *APIServer) DeleteURL(w http.ResponseWriter, r *http.Request) {
	reqUser := r.Context().Value(authUserKey).(*types.User)
	vars := mux.Vars(r)
	urlId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest,
			Map{"error": "invalid url id"})
		return
	}
	deleted, err := s.store.DeleteURL(urlId, reqUser.Id)
	if err != nil {
		WriteJSON(w, http.StatusOK,
			Map{"error": "unknown error"})
		return
	}
	if !deleted {
		WriteJSON(w, http.StatusNotFound,
			Map{"error": "unknown url id"})
		return
	}
	WriteJSON(w, http.StatusNoContent, Map{"msg": "deleted"})
}

func (s *APIServer) UpdateLongURL(w http.ResponseWriter, r *http.Request) {
	reqBody := new(shortReqType)
	reqUser := r.Context().Value(authUserKey).(*types.User)
	vars := mux.Vars(r)
	urlId, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest,
			Map{"error": "invalid url id"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		log.Println("ShortnerJSONError:", err.Error())
		WriteJSON(w, http.StatusInternalServerError,
			Map{"error": "invalid request body"})
		return
	}
	updated, err := s.store.UpdateLongURL(urlId, reqUser.Id, reqBody.LongURL)
	if err != nil {
		log.Println("ShortnerJSONError:", err.Error())
		WriteJSON(w, http.StatusNotFound,
			Map{"error": "url not found with given id"})
		return
	}
	if updated {
		WriteJSON(w, http.StatusOK,
			Map{"msg": "Updated"})
		return
	}
	WriteJSON(w, http.StatusNotFound,
		Map{"error": "url not found or already updated"})
}

func (s *APIServer) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortId := vars["id"]
	longUrl, err := s.store.GetLongURL(shortId)
	if err != nil || longUrl == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Println("Redirect:", longUrl)
	http.Redirect(w, r, longUrl, http.StatusSeeOther)
}
