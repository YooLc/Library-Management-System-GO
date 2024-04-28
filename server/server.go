package server

import (
	"net/http"
	"sync"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

var Mutex = &sync.Mutex{}

func InitServer(config Config) {
	// Configure logrus
	// initLogger()
	mux := http.NewServeMux()

	// Add CORS handler
	corsHandler := cors.AllowAll()
	handler := corsHandler.Handler(mux)

	// Add routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/api/book/add", storeBookHandler)
	mux.HandleFunc("/api/book/adds", storeBooksHandler)
	mux.HandleFunc("/api/book/remove", removeBookHandler)
	mux.HandleFunc("/api/book/query", queryBookHandler)
	mux.HandleFunc("/api/book/stock", incBookStockHandler)
	mux.HandleFunc("/api/book/modify", modifyBookHandler)

	mux.HandleFunc("/api/card/query", showCardsHandler)
	mux.HandleFunc("/api/card/add", registerCardHandler)
	mux.HandleFunc("/api/card/remove", removeCardHandler)

	mux.HandleFunc("/api/borrow/query", showBorrowsHandler)
	mux.HandleFunc("/api/borrow/add", borrowBookHandler)
	mux.HandleFunc("/api/borrow/return", returnBookHandler)

	host, port := config.Host, config.Port
	logrus.Info("Server will run on " + host + ":" + port)
	err := http.ListenAndServe(host+":"+port, handler)
	if err != nil {
		logrus.Panic("Failed to start server: ", err)
		return
	}
}
