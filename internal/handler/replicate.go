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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var options map[string]any
	err := json.NewDecoder(r.Body).Decode(&options)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	continuous := false

	if options["continuous"] != nil {
		val, ok := options["continuous"].(bool)
		if ok {
			continuous = val
		}
	}

	if continuous {
		w.WriteHeader(http.StatusAccepted)
	}
}
