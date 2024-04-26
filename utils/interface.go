package utils

import (
	"library-management-system/database"
	"library-management-system/server/queries"
)

type ServerInterface interface {
	StoreBook(book *database.Book) database.APIResult
	IncBookStock(bookId int, deltaStock int) database.APIResult
	StoreBooks(books []*database.Book) database.APIResult
	RemoveBook(bookId int) database.APIResult
	ModifyBookInfo(book *database.Book) database.APIResult
	QueryBooks(conditions queries.BookQueryConditions) database.APIResult
	BorrowBook(borrow database.Borrow) database.APIResult
	ReturnBook(borrow database.Borrow) database.APIResult
	ShowBorrowHistories(userId int) database.APIResult
	RegisterCard(card *database.Card) database.APIResult
	RemoveCard(cardId int) database.APIResult
	ShowCards() database.APIResult
}
