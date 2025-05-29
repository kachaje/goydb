package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kachaje/goydb/internal/handler"
)

func TestReplicateCode200(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/_replicate", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")

	rr := httptest.NewRecorder()
	p := handler.Replicate{IBase: &MockBase{}}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusOK,
		)
	}
}

func TestReplicateCode400(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/_replicate", bytes.NewBuffer([]byte("bad request")))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")

	rr := httptest.NewRecorder()
	p := handler.Replicate{IBase: &MockBase{}}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusBadRequest,
		)
	}
}

func TestReplicateCode202(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/_replicate", bytes.NewBuffer([]byte(`{"continuous":true}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")

	rr := httptest.NewRecorder()
	p := handler.Replicate{IBase: &MockBase{}}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusAccepted,
		)
	}
}

func TestReplicateCode401(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/_replicate", bytes.NewBuffer([]byte(`{"continuous":true}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secrets")

	rr := httptest.NewRecorder()
	p := handler.Replicate{IBase: &MockBase{}}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusUnauthorized,
		)
	}
}
