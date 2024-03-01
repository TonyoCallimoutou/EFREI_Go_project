package infrastructure

import (
	"database/sql"
	"errors"
	"fmt"
	"go_shortener/src/domain"
	"go_shortener/src/interface/database"
	"log"
	"os"
	"time"

	"github.com/blockloop/scan/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var expired_time = 10 * time.Second

type LinkStore struct {
	db *sql.DB
}

func NewLinkStoreMySQL() database.LinkStore {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	shortenerSQL := os.Getenv("MYSQL_USER")
	passwordSQL := os.Getenv("MYSQL_PASSWORD")
	dbSQL := os.Getenv("MYSQL_DATABASE")
	urlSQL := os.Getenv("MYSQL_URL")
	portSQL := os.Getenv("MYSQL_PORT")

	connection := shortenerSQL + ":" + passwordSQL + "@tcp(" + urlSQL + ":" + portSQL + ")/" + dbSQL + "?parseTime=true"

	fmt.Println(connection)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		panic(err.Error)
	}
	linkStore := new(LinkStore)
	linkStore.db = db
	return linkStore
}

func (handler *LinkStore) Create(url domain.Shortener) error {

	fmt.Println("coucou4")

	url.ID = uuid.NewString()
	url.ExpiredAt = time.Now().Add(expired_time)

	_, errSQL := handler.db.Exec("INSERT INTO Shortener (id, url, shortUrl, expiredAt) VALUES (?,?,?,?)",
		url.ID, url.Url, url.ShortUrl, url.ExpiredAt)
	return errSQL
}

func (handler *LinkStore) GetById(url domain.Shortener) (string, error) {

	var originalURL string
	err := handler.db.QueryRow("SELECT url FROM Shortener WHERE shortUrl=?", url.ShortUrl).Scan(&originalURL)
	return originalURL, err
}

func (handler *LinkStore) GetAll() ([]domain.Shortener, error) {
	var URLArray []domain.Shortener
	rows, err := handler.db.Query("SELECT * FROM Shortener")
	if err != nil {
		return URLArray, err
	}

	err = scan.Rows(&URLArray, rows)
	if err != nil {
		return URLArray, err
	}
	return URLArray, nil
}

func (handler *LinkStore) Update(url domain.Shortener) error {

	url.ExpiredAt = time.Now().Add(expired_time)
	res, errSQL := handler.db.Exec("Update Shortener set url = ?,expiredAt = ? where shortUrl = ?",
		url.Url, url.ExpiredAt, url.ShortUrl)
	if errSQL != nil {
		return errSQL
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("short URL not found")
	}
	return nil
}

func (handler *LinkStore) Delete(url domain.Shortener) error {

	res, errSQL := handler.db.Exec("Delete from Shortener where shortUrl = ?", url.ShortUrl)
	if errSQL != nil {
		return errSQL
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("short URL not found")
	}
	return nil
}
