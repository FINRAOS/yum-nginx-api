package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/go-ozzo/ozzo-routing/content"
)

func newRouter() *routing.Router {
	rtr := routing.New()
	rtr.Group("/api", content.TypeNegotiator(content.JSON))
	return rtr
}
func testHealthHandler(t *testing.T, rtr *routing.Router) {
	rtr.Get("/health", healthRoute)
	req, err := http.NewRequest("GET", "/health", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	rtr.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `OK`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func testRepoHandler(t *testing.T, rtr *routing.Router) {
	rtr.Get("/repo", repoRoute)
	req, err := http.NewRequest("GET", "/repo", nil)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	rtr.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `[]`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestHealth(t *testing.T) {
	testHealthHandler(t, newRouter())
}

func TestRepo(t *testing.T) {
	testRepoHandler(t, newRouter())
}
