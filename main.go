package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Shortener struct {
	ID        string    `json:"id"`
	Url       string    `json:"url"`
	ShortUrl  string    `json:"shortUrl"`
	ExpiredAt time.Time `json:"expired_at"`
}

var expired_time = 10 * time.Second

var db *sql.DB

func main() {
	address := flag.String("address", ":4000", "address to listen to")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	userSQL := os.Getenv("MYSQL_USER")
	passwordSQL := os.Getenv("MYSQL_PASSWORD")
	dbSQL := os.Getenv("MYSQL_DATABASE")
	urlSQL := os.Getenv("MYSQL_URL")
	portSQL := os.Getenv("MYSQL_PORT")

	connection := userSQL + ":" + passwordSQL + "@tcp(" + urlSQL + ":" + portSQL + ")/" + dbSQL

	db, err = sql.Open("mysql", connection)
	if err != nil {
		panic(err)
	}

	router := RunServer()
	errServ := http.ListenAndServe(*address, router)
	if errServ != nil {
		log.Fatalf("error: %s", errServ.Error())
	}
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/index.html")
}

func getMainStyle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/styles/main.css")
}

func getResetStyle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/styles/reset.css")
}

func getScript(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	http.ServeFile(w, r, "./public/scripts/main.js")
}

func RunServer() http.Handler {
	router := chi.NewRouter()
	router.Get("/", getRoot)
	router.Route("/styles", func(r chi.Router) {
		r.Get("/main", getMainStyle)
		r.Get("/reset", getResetStyle)
	})
	router.Route("/scripts", func(r chi.Router) {
		r.Get("/main", getScript)
	})

	router.Route("/api", func(r chi.Router) {
		r.Post("/", createShortener)
		r.Get("/{url}", getShortener)
		r.Put("/", updateShortener)
		r.Delete("/", deleteShortener)
	})

	return router
}

func createShortener(w http.ResponseWriter, r *http.Request) {
	url := &Shortener{}
	err := json.NewDecoder(r.Body).Decode(url)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	url.ID = uuid.NewString()
	url.ExpiredAt = time.Now().Add(expired_time)

	_, errSQL := db.Exec("INSERT INTO Shortener (id, url, shortUrl, expiredAt) VALUES (?,?,?,?)",
		url.ID, url.Url, url.ShortUrl, url.ExpiredAt)
	if errSQL != nil {
		http.Error(w, errSQL.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(url)

}

func getShortener(w http.ResponseWriter, r *http.Request) {
	shortUrl := chi.URLParam(r, "url")

	var originalURL string
	err := db.QueryRow("SELECT url FROM Shortener WHERE shortUrl=?", shortUrl).Scan(&originalURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func updateShortener(w http.ResponseWriter, r *http.Request) {
	updateUrl := &Shortener{}
	err := json.NewDecoder(r.Body).Decode(updateUrl)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	updateUrl.ExpiredAt = time.Now().Add(expired_time)
	res, errSQL := db.Exec("Update Shortener set url = ?,expiredAt = ? where shortUrl = ?",
		updateUrl.Url, updateUrl.ExpiredAt, updateUrl.ShortUrl)
	if errSQL != nil {
		http.Error(w, errSQL.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "short URL not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updateUrl)
	return

}

func deleteShortener(w http.ResponseWriter, r *http.Request) {
	deleteUrl := &Shortener{}
	err := json.NewDecoder(r.Body).Decode(deleteUrl)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	res, errSQL := db.Exec("Delete from Shortener where shortUrl = ?", deleteUrl.ShortUrl)
	if errSQL != nil {
		http.Error(w, errSQL.Error(), http.StatusInternalServerError)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "short URL not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deleteUrl)
	return
}
