package server

import (
	"encoding/json"
	"library-management-system/database"
	"net/http"
	"strconv"
)

func showCardsHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	server := Server{}
	result := server.ShowCards()
	server.Response(w, result)
}

func registerCardHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	server := Server{}
	var cardData database.Card
	if err := json.NewDecoder(r.Body).Decode(&cardData); err != nil {
		server.Response(w, database.APIResult{
			Ok:      false,
			Message: "Invalid Arguments: failed to parse request body",
			Payload: nil,
		})
		return
	}

	res := server.RegisterCard(&cardData)
	server.Response(w, res)
}

func removeCardHandler(w http.ResponseWriter, r *http.Request) {
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
	res := server.RemoveCard(cardId)
	server.Response(w, res)
}
