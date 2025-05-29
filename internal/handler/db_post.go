package handler

import (
	"encoding/json"
	"net/http"
)

type DBPost struct {
	IBase
}

func (s *DBPost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if _, ok := s.Authenticate(w, r); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	doc, code, err := s.PutDoc(w, r)
	if err != nil {
		WriteError(w, code, err.Error())
		return
	}

	if doc["_id"] == nil {
		WriteError(w, http.StatusInternalServerError, "no id attached")
		return
	} else if doc["_rev"] == nil {
		WriteError(w, http.StatusInternalServerError, "no rev attached")
		return
	}

	id, ok := doc["_id"].(string)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "failed to load id")
		return
	}
	rev, ok := doc["_rev"].(string)
	if !ok {
		WriteError(w, http.StatusInternalServerError, "failed to load rev")
		return
	}

	response := PostResponse{
		ID:  id,
		OK:  true,
		Rev: rev,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

type PostResponse struct {
	ID  string `json:"id"`
	OK  bool   `json:"ok"`
	Rev string `json:"rev"`
}
