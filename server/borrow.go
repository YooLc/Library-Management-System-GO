package server

import (
	"encoding/json"
	"library-management-system/database"
	"net/http"
	"strconv"
)

func showBorrowsHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	server := Server{}
	params := r.URL.Query()
	cardIdStr := params.Get("card_id")
	var err error
	var cardId int
	if cardId, err = strconv.Atoi(cardIdStr); err != nil || cardId <= 0 {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request parameter, expect positive integer",
			Payload: nil,
		})
		return
	}
	result := server.ShowBorrowHistories(cardId)
	server.Response(w, result)
}

func borrowBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()

	// Parse request body
	server := Server{}
	var borrow database.Borrow
	err := json.NewDecoder(r.Body).Decode(&borrow)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Borrow book
	borrow.ReturnTime = 0 // make sure ReturnTime is 0
	result := server.BorrowBook(borrow)
	server.Response(w, result)
}

func returnBookHandler(w http.ResponseWriter, r *http.Request) {
	// Lock Mutex
	Mutex.Lock()
	defer Mutex.Unlock()

	// Parse request body
	server := Server{}
	var borrow database.Borrow
	err := json.NewDecoder(r.Body).Decode(&borrow)
	if err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	// Return book
	if borrow.ReturnTime == 0 {
		borrow.ResetReturnTime()
	}
	result := server.ReturnBook(borrow)
	server.Response(w, result)
}
