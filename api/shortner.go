package api

import (
	"encoding/json"
	"log"
	"net/http"

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
