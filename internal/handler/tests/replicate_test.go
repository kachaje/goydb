package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/sessions"
	"github.com/kachaje/goydb/internal/adapter/storage"
	"github.com/kachaje/goydb/internal/handler"
)

func TestReplicate(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/_replicate", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "secret")

	rr := httptest.NewRecorder()
	p := handler.Replicate{}

	store := sessions.NewCookieStore([]byte("admin:secret"))
	p.SessionStore = store
	p.Storage = &storage.Storage{}

	hnd := http.HandlerFunc(p.ServeHTTP)

	hnd.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v; want %v",
			status, http.StatusOK,
		)
	}
}
