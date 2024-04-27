package server

import (
	"encoding/json"
	"library-management-system/database"
	"library-management-system/server/queries"
	"net/http"
	"strconv"
	"sync"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Frontend string `yaml:"frontend"`
}

var mutex = &sync.Mutex{}

func response(w http.ResponseWriter, resp database.APIResult) {
	w.Header().Set("Content-Type", "application/json")
	bytes, _ := json.Marshal(resp)
	_, err := w.Write(bytes)
	if err != nil {
		logrus.Fatal("unable to response")
		return
	}
}

func storeBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock mutex
	mutex.Lock()
	defer mutex.Unlock()
	server := Server{}

	// Parse request body
	var book database.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Store book
	result := server.StoreBook(&book)
	response(w, result)
}

func showCardsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}
	result := server.ShowCards()
	response(w, result)
}

func registerCardHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}
	var cardData database.Card
	if err := json.NewDecoder(r.Body).Decode(&cardData); err != nil {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	res := server.RegisterCard(&cardData)
	response(w, res)
}

func removeCardHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}

	params := r.URL.Query()
	cardIdStr := params.Get("card_id")
	var err error
	var cardId int
	if cardId, err = strconv.Atoi(cardIdStr); err != nil || cardId <= 0 {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request parameter, expect positive integer",
			Payload: nil,
		})
		return
	}
	res := server.RemoveCard(cardId)
	response(w, res)
}

func showBorrowsHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}
	params := r.URL.Query()
	cardIdStr := params.Get("card_id")
	var err error
	var cardId int
	if cardId, err = strconv.Atoi(cardIdStr); err != nil || cardId <= 0 {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request parameter, expect positive integer",
			Payload: nil,
		})
		return
	}
	result := server.ShowBorrowHistories(cardId)
	response(w, result)
}

func removeBookHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}
	params := r.URL.Query()
	bookIdStr := params.Get("book_id")
	var err error
	var bookId int
	if bookId, err = strconv.Atoi(bookIdStr); err != nil || bookId <= 0 {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request parameter, expect positive integer",
			Payload: nil,
		})
		return
	}
	result := server.RemoveBook(bookId)
	response(w, result)
}

func queryBookHandler(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	server := Server{}
	params := r.URL.Query()

	var err error
	var minPublishYear int
	var maxPublishYear int
	var minPrice float64
	var maxPrice float64
	if minPublishYear, err = strconv.Atoi(params.Get("min_publish_year")); err != nil {
		minPublishYear = 0
	}
	if maxPublishYear, err = strconv.Atoi(params.Get("max_publish_year")); err != nil {
		maxPublishYear = 0
	}
	if minPrice, err = strconv.ParseFloat(params.Get("min_price"), 64); err != nil {
		minPrice = 0
	}
	if maxPrice, err = strconv.ParseFloat(params.Get("max_price"), 64); err != nil {
		maxPrice = 0
	}
	condition := queries.BookQueryConditions{
		Category:       params.Get("category"),
		Title:          params.Get("title"),
		Press:          params.Get("press"),
		MinPublishYear: minPublishYear,
		MaxPublishYear: maxPublishYear,
		Author:         params.Get("author"),
		MinPrice:       minPrice,
		MaxPrice:       maxPrice,
		SortBy:         queries.SortColumn(params.Get("sort_by")),
		SortOrder:      queries.Order(params.Get("sort_order")),
	}
	result := server.QueryBooks(condition)
	response(w, result)
}

func incBookStockHandler(w http.ResponseWriter, r *http.Request) {
	// Lock mutex
	mutex.Lock()
	defer mutex.Unlock()

	// Parse request body
	server := Server{}
	type IncStockQuery struct {
		BookId     int `json:"book_id"`
		DeltaStock int `json:"delta_stock"`
	}
	var query IncStockQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Increment book stock
	result := server.IncBookStock(query.BookId, query.DeltaStock)
	response(w, result)
}

func InitServer(config Config) {
	// Configure logrus
	// initLogger()
	mux := http.NewServeMux()

	// Add CORS handler
	corsHandler := cors.AllowAll()
	handler := corsHandler.Handler(mux)

	//fs := http.FileServer(http.Dir(config.Frontend))

	// Set X-Content-Type-Options to nosniff to enhance security
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			//fs.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	})

	mux.HandleFunc("/store", storeBookHandler)
	mux.HandleFunc("/card/query", showCardsHandler)
	mux.HandleFunc("/card/add", registerCardHandler)
	mux.HandleFunc("/card/remove", removeCardHandler)
	mux.HandleFunc("/borrow/query", showBorrowsHandler)
	mux.HandleFunc("/book/remove", removeBookHandler)
	mux.HandleFunc("/book/query", queryBookHandler)
	mux.HandleFunc("/book/stock", incBookStockHandler)

	host, port := config.Host, config.Port
	logrus.Info("Server will run on " + host + ":" + port)
	err := http.ListenAndServe(host+":"+port, handler)
	if err != nil {
		logrus.Panic("Failed to start server: ", err)
		return
	}
}
