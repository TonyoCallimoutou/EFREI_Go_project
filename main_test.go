package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateShortener(t *testing.T) {
	// Création d'une nouvelle requête HTTP simulée avec un corps JSON
	shortenerBody := `{"url": "http://example.com", "shortUrl": "test"}`
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(shortenerBody))
	if err != nil {
		t.Fatal(err)
	}
	// Définition de l'en-tête Content-Type pour indiquer que le corps est au format JSON
	req.Header.Set("Content-Type", "application/json")

	// Création d'un enregistrement de réponse HTTP simulé
	rr := httptest.NewRecorder()

	// Création d'une fonction mock SQL
	hh, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer hh.Close()

	db = hh

	// Définition des attentes de la requête SQL simulée
	mock.ExpectExec("INSERT INTO Shortener").WillReturnResult(sqlmock.NewResult(1, 1))

	// Appel de la fonction à tester
	createShortener(rr, req)

	type ResponseBody struct {
		Url      string `json:"Url"`
		ShortUrl string `json:"ShortUrl"`
	}

	// Vérification du code de statut HTTP retourné
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("error unmarshaling response body: %v", err)
	}
	expected := ResponseBody{
		Url:      "http://example.com",
		ShortUrl: "test",
	}

	if responseBody != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			responseBody, expected)
	}

	fmt.Println(rr.Body.String())

	// Vérification des attentes SQL simulées
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
