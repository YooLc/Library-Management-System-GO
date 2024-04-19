package server

import (
	"encoding/json"
	"library-management-system/database"
	"net/http"
	"sync"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var mutex = &sync.Mutex{}

func response(w http.ResponseWriter, resp APIResult) {
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(resp)
	w.Write(bytes)
}

func storeBook(w http.ResponseWriter, r *http.Request) {
	// Lock mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Parse request body
	var book database.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		response(w, APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Store book
	result := StoreBook(book)
	response(w, result)
}

func InitServer(config ServerConfig) {
	// Configure logrus
	// initLogger()
	mux := http.NewServeMux()

	// Add CORS handler
	corsHandler := cors.Default()
	handler := corsHandler.Handler(mux)

	// fs := http.FileServer(http.Dir(frontend))

	// Set X-Content-Type-Options to nosniff to enhance security
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			// fs.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})

	mux.HandleFunc("/store", storeBook)

	host, port := config.Host, config.Port
	logrus.Info("Server is running on " + host + ":" + port)
	http.ListenAndServe(host+":"+port, handler)
}
