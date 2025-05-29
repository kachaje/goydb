package handler

import (
	"encoding/json"
	"net/http"
)

type Replicate struct {
	IBase
}

func (s *Replicate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	if _, ok := s.Authenticate(w, r); !ok {
		return
	}

	var doc map[string]any
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
}
