package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateShortener(t *testing.T) {
	// Création d'une nouvelle requête HTTP simulée avec un corps JSON
	shortenerBody := `{"url": "http://example.com", "shortUrl": "caca"}`
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(shortenerBody))
	if err != nil {
		t.Fatal(err)
	}
	// Définition de l'en-tête Content-Type pour indiquer que le corps est au format JSON
	req.Header.Set("Content-Type", "application/json")

	// Création d'un enregistrement de réponse HTTP simulé
	rr := httptest.NewRecorder()

	// Création d'une fonction mock SQL
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()


	// Définition des attentes de la requête SQL simulée
	mock.ExpectExec("INSERT INTO Shortener").WillReturnResult(sqlmock.NewResult(1, 1))

	
	fmt.Println(shortenerBody)

	// Appel de la fonction à tester
	createShortener(rr, req)

	
	fmt.Println(shortenerBody)

	// Vérification du code de statut HTTP retourné
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Vérification du corps de la réponse HTTP
	expected := `{"ID":"testID","Url":"http://example.com","ShortUrl":"test","ExpiredAt":"0001-01-01T00:00:00Z"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Vérification des attentes SQL simulées
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
