package infrastructure

import (
	"encoding/json"
	"flag"
	"go_shortener/src/domain"
	"go_shortener/src/interface/database"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid"
)

func Init() {
	address := flag.String("address", ":4000", "address to listen to")
	flag.Parse()
	// Initialize the shortener controller
	shortenerController := NewLinkStoreMySQL()

	// Start the server
	router := RunServer(shortenerController)
	errServ := http.ListenAndServe(*address, router)
	if errServ != nil {
		log.Fatalf("error: %s", errServ.Error())
	}
}

func RunServer(shortenerController database.LinkStore) http.Handler {
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
		r.Post("/", createShortener(shortenerController))
		r.Get("/", getAllShortener(shortenerController))
		r.Get("/redirect/{url}", getShortener(shortenerController))
		r.Put("/", updateShortener(shortenerController))
		r.Delete("/{url}", deleteShortener(shortenerController))
	})

	return router
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/index.html")
}

func getMainStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/styles/main.css")
}

func getResetStyle(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/styles/reset.css")
}

func getScript(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/scripts/main.js")
}

func getShortener(shortenerController database.LinkStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shortUrl := chi.URLParam(r, "url")
		url := &domain.Shortener{}
		url.ShortUrl = shortUrl
		originalURL, err := shortenerController.GetById(*url)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusSeeOther)
	}
}

func getAllShortener(shortenerController database.LinkStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		URLArray, err := shortenerController.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonResponse(w, http.StatusOK, URLArray)
	}
}

func createShortener(shortenerController database.LinkStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := &domain.Shortener{}
		err := json.NewDecoder(r.Body).Decode(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		url.ID = uuid.NewString()
		url.ExpiredAt = time.Now().Add(expired_time)

		if url.ShortUrl == "" {
			url.ShortUrl = shortuuid.New()
		}

		err = shortenerController.Create(*url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, http.StatusOK, url)
	}
}

func updateShortener(shortenerController database.LinkStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := &domain.Shortener{}
		err := json.NewDecoder(r.Body).Decode(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = shortenerController.Update(*url)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		jsonResponse(w, http.StatusOK, url)
	}
}

func deleteShortener(shortenerController database.LinkStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deleteUrl := &domain.Shortener{}
		deleteUrl.ShortUrl = chi.URLParam(r, "url")

		err := shortenerController.Delete(*deleteUrl)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonResponse(w, http.StatusOK, deleteUrl)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
