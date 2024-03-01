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

	"github.com/blockloop/scan/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v3"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Shortener struct {
	ID        string    `json:"id" db:"id"`
	Url       string    `json:"url" db:"url"`
	ShortUrl  string    `json:"shortUrl" db:"shortUrl"`
	ExpiredAt time.Time `json:"expired_at" db:"expiredAt"`
	Count     int       `json:"count" db:"count"`
}

type ErrorResponse struct {
	Message string `json:"message"`
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

	connection := userSQL + ":" + passwordSQL + "@tcp(" + urlSQL + ":" + portSQL + ")/" + dbSQL + "?parseTime=true"

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
		r.Get("/", getAllUrls)
		r.Post("/", createShortener)
		r.Get("/redirect/{url}", getShortener)
		r.Put("/", updateShortener)
		r.Delete("/{url}", deleteShortener)
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

	if url.ShortUrl == "" {
		url.ShortUrl = shortuuid.New()
	}

	_, errSQL := db.Exec("INSERT INTO Shortener (id, url, shortUrl, expiredAt) VALUES (?,?,?,?)",
		url.ID, url.Url, url.ShortUrl, url.ExpiredAt)
	if errSQL != nil {
		handleSQLError(w, errSQL)
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
	_, errSQL := db.Exec("UPDATE Shortener set count = count+1 where shortUrl = ?", shortUrl)
	if errSQL != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusSeeOther)
}

func getAllUrls(w http.ResponseWriter, r *http.Request) {
	var URLArray []Shortener
	rows, err := db.Query("SELECT * FROM Shortener")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = scan.Rows(&URLArray, rows)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	json.NewEncoder(w).Encode(URLArray)
}

func updateShortener(w http.ResponseWriter, r *http.Request) {
	updateUrl := &Shortener{}

	updateUrl.ExpiredAt = time.Now().Add(expired_time)
	res, errSQL := db.Exec("Update Shortener set url = ?,expiredAt = ? where shortUrl = ?",
		updateUrl.Url, updateUrl.ExpiredAt, updateUrl.ShortUrl)
	if errSQL != nil {
		handleSQLError(w, errSQL)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		handleSQLError(w, err)
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
	shortUrl := chi.URLParam(r, "url")

	res, errSQL := db.Exec("Delete from Shortener where shortUrl = ?", shortUrl)
	if errSQL != nil {
		handleSQLError(w, errSQL)
		return
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		handleSQLError(w, err)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "short URL not found", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(deleteUrl)
	return
}

func handleSQLError(w http.ResponseWriter, errSQL error) {
	errorResponse := ErrorResponse{Message: errSQL.Error()}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonResponse)
}
