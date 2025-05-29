package handler

import (
	"fmt"
	"net/http"
)

type Replicate struct {
	Base
}

func (s *Replicate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()

	if _, ok := (Authenticator{Base: s.Base, RequiresAdmin: true}.Do(w, r)); !ok {
		return
	}
	fmt.Println("Past")
}
