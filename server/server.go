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
		return
	})

	mux.HandleFunc("/book/add", storeBookHandler)
	mux.HandleFunc("/book/adds", storeBooksHandler)
	mux.HandleFunc("/book/remove", removeBookHandler)
	mux.HandleFunc("/book/query", queryBookHandler)
	mux.HandleFunc("/book/stock", incBookStockHandler)
	mux.HandleFunc("/book/modify", modifyBookHandler)

	mux.HandleFunc("/card/query", showCardsHandler)
	mux.HandleFunc("/card/add", registerCardHandler)
	mux.HandleFunc("/card/remove", removeCardHandler)

	mux.HandleFunc("/borrow/query", showBorrowsHandler)
	mux.HandleFunc("/borrow/add", borrowBookHandler)
	mux.HandleFunc("/borrow/return", returnBookHandler)

	host, port := config.Host, config.Port
	logrus.Info("Server will run on " + host + ":" + port)
	err := http.ListenAndServe(host+":"+port, handler)
	if err != nil {
		logrus.Panic("Failed to start server: ", err)
		return
	}
}
