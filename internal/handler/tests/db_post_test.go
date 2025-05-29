package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kachaje/goydb/internal/handler"
)

func TestDBPostWithId(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/sample", bytes.NewBuffer([]byte(`{"_id":"test","name":"Test Name"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")

	rr := httptest.NewRecorder()
	p := handler.DBPost{IBase: &MockBase{}}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusOK,
		)
	}
}

func TestDBPostWithoutId(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/sample", bytes.NewBuffer([]byte(`{"name":"Test Name"}`)))
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
