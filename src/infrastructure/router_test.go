package infrastructure

import (
	"bytes"
	"encoding/json"
	"go_shortener/src/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockLinkStore struct{}

func (m *mockLinkStore) Create(url domain.Shortener) error {
	return nil
}

func (m *mockLinkStore) GetById(url domain.Shortener) (string, error) {
	return "http://example.com", nil
}

func (m *mockLinkStore) GetAll() ([]domain.Shortener, error) {
	return []domain.Shortener{{ID: "1", Url: "http://example.com", ShortUrl: "abc"}}, nil
}

func (m *mockLinkStore) Update(url domain.Shortener) error {
	return nil
}

func (m *mockLinkStore) Delete(url domain.Shortener) error {
	return nil
}

func TestGetAllShortener(t *testing.T) {
	// Mocking the LinkStore
	linkStore := &mockLinkStore{}

	// Creating a new request to get all shorteners
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Creating a response recorder
	rr := httptest.NewRecorder()

	// Creating a handler function with the mock LinkStore
	handler := http.HandlerFunc(getAllShortener(linkStore))

	// Serving the request to the handler
	handler.ServeHTTP(rr, req)

	// Checking the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parsing the response body
	var response []domain.Shortener
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Checking the response data
	if len(response) != 1 {
		t.Errorf("unexpected number of shorteners: got %d want %d", len(response), 1)
	}
}

func TestCreateShortener(t *testing.T) {
	// Mocking the LinkStore
	linkStore := &mockLinkStore{}

	// Creating a sample request body
	url := domain.Shortener{
		Url:       "http://example.com",
		ShortUrl:  "abc",
	}
	body, err := json.Marshal(url)
	if err != nil {
		t.Fatal(err)
	}

	// Creating a new request to create a shortener
	req, err := http.NewRequest("POST", "/api", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Creating a response recorder
	rr := httptest.NewRecorder()

	// Creating a handler function with the mock LinkStore
	handler := http.HandlerFunc(createShortener(linkStore))

	// Serving the request to the handler
	handler.ServeHTTP(rr, req)

	// Checking the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parsing the response body
	var response domain.Shortener
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Checking the response data
	if response.Url != url.Url || response.ShortUrl != url.ShortUrl {
		t.Errorf("unexpected response: got %v want %v", response, url)
	}
}

func TestUpdateShortener(t *testing.T) {
	// Mocking the LinkStore
	linkStore := &mockLinkStore{}

	// Creating a sample request body
	url := domain.Shortener{
		ShortUrl: "abc",
	}
	body, err := json.Marshal(url)
	if err != nil {
		t.Fatal(err)
	}

	// Creating a new request to update a shortener
	req, err := http.NewRequest("PUT", "/api", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Creating a response recorder
	rr := httptest.NewRecorder()

	// Creating a handler function with the mock LinkStore
	handler := http.HandlerFunc(updateShortener(linkStore))

	// Serving the request to the handler
	handler.ServeHTTP(rr, req)

	// Checking the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parsing the response body
	var response domain.Shortener
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Checking the response data
	if response.ShortUrl != url.ShortUrl {
		t.Errorf("unexpected response: got %v want %v", response, url)
	}
}

func TestDeleteShortener(t *testing.T) {
	// Mocking the LinkStore
	linkStore := &mockLinkStore{}

	// Creating a new request to delete a shortener
	req, err := http.NewRequest("DELETE", "/api/abc", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Creating a response recorder
	rr := httptest.NewRecorder()

	// Creating a handler function with the mock LinkStore
	handler := http.HandlerFunc(deleteShortener(linkStore))

	// Serving the request to the handler
	handler.ServeHTTP(rr, req)

	// Checking the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Parsing the response body
	var response domain.Shortener
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("error unmarshalling response: %v", err)
	}

	// Checking the response data
	if response.ShortUrl != "abc" {
		t.Errorf("unexpected response: got %v want abc", response.ShortUrl)
	}
}
