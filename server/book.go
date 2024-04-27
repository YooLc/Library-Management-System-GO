package server

import (
	"encoding/json"
	"library-management-system/database"
	"library-management-system/server/queries"
	"net/http"
	"strconv"
)

func storeBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()
	server := Server{}

	// Parse request body
	var book database.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Store book
	result := server.StoreBook(&book)
	server.Response(w, result)
}

func storeBooksHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()
	server := Server{}

	// Parse request body
	var list queries.BookList
	err := json.NewDecoder(r.Body).Decode(&list)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Store books
	insList := []*database.Book{}
	for _, book := range list.Books {
		insList = append(insList, &book)
	}
	result := server.StoreBooks(insList)
	server.Response(w, result)
}

func incBookStockHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()

	// Parse request body
	server := Server{}
	type IncStockQuery struct {
		BookId     int `json:"book_id"`
		DeltaStock int `json:"delta_stock"`
	}
	var query IncStockQuery
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Increment book stock
	result := server.IncBookStock(query.BookId, query.DeltaStock)
	server.Response(w, result)
}

func modifyBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()

	// Parse request body
	server := Server{}
	var book database.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Modify book
	result := server.ModifyBookInfo(&book)
	server.Response(w, result)
}

func removeBookHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	server := Server{}
	params := r.URL.Query()
	bookIdStr := params.Get("book_id")
	var err error
	var bookId int
	if bookId, err = strconv.Atoi(bookIdStr); err != nil || bookId <= 0 {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request parameter, expect positive integer",
			Payload: nil,
		})
		return
	}
	result := server.RemoveBook(bookId)
	server.Response(w, result)
}

func queryBookHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

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
	server.Response(w, result)
}
