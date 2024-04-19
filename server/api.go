package server

import (
	"library-management-system/database"
	"library-management-system/server/queries"
)

type APIResult struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func StoreBook(book database.Book) APIResult {
	if err := database.DB.Create(&book); err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to store book",
			Payload: nil,
		}
	}
	incrementId := book.Book_id + 1
	return APIResult{
		Ok:      true,
		Message: "Book stored successfully",
		Payload: incrementId,
	}
}

func IncBookStock(book database.Book, count int) APIResult {
	if err := database.DB.Model(&book).Update("stock", book.Stock+count); err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to increment book stock",
			Payload: nil,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book stock incremented successfully",
		Payload: nil,
	}
}

func StoreBookList(books queries.BookList) APIResult {
	for _, book := range books.Books {
		resp := StoreBook(book)
		if !resp.Ok {
			return resp
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Books stored successfully",
		Payload: nil,
	}
}

func RemoveBook(book database.Book) APIResult {
	if err := database.DB.Delete(&book); err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to remove book",
			Payload: nil,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book removed successfully",
		Payload: nil,
	}
}

func ModifyBookInfo(book database.Book) APIResult {
	return APIResult{}
}

func QueryBook(conditions queries.BookQueryConditions) APIResult {
	return APIResult{}
}

func BorrowBook(borrow database.Borrow) APIResult {
	return APIResult{}
}

func ReturnBook(borrow database.Borrow) APIResult {
	return APIResult{}
}

func ShowBorrowHistories(card_id int) APIResult {
	return APIResult{}
}

func RegisterCard(card database.Card) APIResult {
	return APIResult{}
}

func RemoveCard(card_id int) APIResult {
	return APIResult{}
}

func ShowCards() APIResult {
	return APIResult{}
}
