package utils

import "library-management-system/database"

type ServerInterface interface {
	StoreBook(book *database.Book) database.APIResult
	StoreBooks(books []*database.Book) database.APIResult
	RegisterCard(card *database.Card) database.APIResult
	BorrowBook(borrow database.Borrow) database.APIResult
	ReturnBook(borrow database.Borrow) database.APIResult
}
