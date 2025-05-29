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

	var doc map[string]any
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
}
