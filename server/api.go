package server

import (
	"fmt"
	"library-management-system/database"
	"library-management-system/server/queries"
	"time"

	"gorm.io/gorm"
)

type APIResult struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

// StoreBook
// register a book to database.
//
//	Note that:
//	     (1) book_id should be stored to book after successfully
//	         completing this operation.
//	     (2) you should not register this book if the book already
//	         exists in the library system.
//
//	@param book all attributes of the book
func StoreBook(book database.Book) APIResult {
	if err := database.DB.Create(&book).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to store book, maybe the book already exists",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book stored successfully",
		Payload: book.Book_id,
	}
}

// IncBookStock
// increase the book's inventory by bookId & deltaStock.
//
// Note that:
//
//	(1) you need to check the correctness of book_id
//	(2) deltaStock can be negative, but make sure that
//	    the result of book.stock + deltaStock is not negative!
//
// @param bookId book's book_id
// @param deltaStock increase count to book's stock, must be greater
func IncBookStock(bookId int, count int) APIResult {
	// Check the correctness of book_id
	book := database.Book{}
	if err := database.DB.First(&book, bookId).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "Book not found",
			Payload: nil,
		}
	}

	// Check the result of book.stock + deltaStock is not negative
	if book.Stock+count < 0 {
		return APIResult{
			Ok:      false,
			Message: "Stock count becomes invalid after incrementing, please check the arguments",
			Payload: nil,
		}
	}

	// Performing the increment operation
	if err := database.DB.Model(&book).Update("stock", book.Stock+count).Error; err != nil {
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

// StoreBookList
// batch store books.
//
// Note that:
//
//	(1) you should not call the interface StoreBook()
//	    multiple times to achieve this function!!!
//	    ? hint: use {@link PreparedStatement#executeBatch()}
//	    ? and {@link PreparedStatement#addBatch()}
//	(2) if one of the books fails to import, all operations
//	    should be rolled back using rollback() function provided
//	    by JDBC!!!
//	(3) when binding params to SQL, you are required to avoid
//	    the risk of SQL injection attack!!!
//
// @param books list of books to be stored
func StoreBookList(books queries.BookList) APIResult {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		for _, book := range books.Books {
			if err := tx.Create(&book).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to store books, maybe one of them already exists",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Books stored successfully",
		Payload: nil,
	}
}

// RemoveBook
//
//	remove this book from library system.
//
//	Note that if someone has not returned this book,
//	the book should not be removed!
//
//	@param bookId the book to be removed
func RemoveBook(book database.Book) APIResult {
	// Check if someone has not returned this book

	// Remove the book
	if err := database.DB.Delete(&book).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to remove book",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book removed successfully",
		Payload: nil,
	}
}

// ModifyBookInfo
// modify a book's information by book_id.book_id.
//
// Note that you should not modify its book_id and stock!
//
// @param book the book to be modified
func ModifyBookInfo(book database.Book) APIResult {
	// Todo: Avoid modifying book_id and stock

	// Modify the book info
	if err := database.DB.Model(&book).Updates(book); err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to modify book info",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book info modified successfully",
		Payload: book,
	}
}

// QueryBook
// Note that:
//
//	(1) you should let the DBMS filter records
//	    that do not satisfy the conditions instead of
//	    filter records in your API.
//	(2) when binding params to SQL, you also need to avoid
//	    the risk of SQL injection attack.
//	(3) [*] if all else is equal, sort by book_id in
//	    ascending order!
//
// @param conditions query conditions
//
// @return query results should be returned by ApiResult.payload
//
//	and should be an instance of {@link queries.BookQueryResults}
func queryBook(conditions queries.BookQueryConditions) APIResult {
	return APIResult{}
}

/* Interface for borrow & return books */

// BorrowBook
//
// a user borrows one book with the specific card.
// the borrow operation will success iff there are
// enough books in stock & the user has not borrowed
// the book or has returned it.
//
// @param borrow information, include borrower &
// book's id & time
func BorrowBook(borrow database.Borrow) APIResult {
	borrow.Borrow_time = time.Now().UnixMilli()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&borrow).Error; err != nil {
			return err
		}

		stock := 0
		tx.Model(&database.Book{}).Select("stock").Where("book_id = ?", borrow.Book_id).Row().Scan(&stock)
		if stock <= 0 {
			return fmt.Errorf("book out of stock")
		}

		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.Book_id).Update("stock", gorm.Expr("stock - 1")).Error; err != nil {
			return err
		}

		// Return nil to commit the transaction
		return nil
	})

	if err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to borrow book",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book borrowed successfully",
		Payload: nil,
	}
}

// ReturnBook
//
// A user return one book with specific card.
//
// @param borrow
// borrow information, include borrower & book's id & return time
func ReturnBook(borrow database.Borrow) APIResult {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// return_time = 0 because a book can be borrowed multiple times by the same card
		if err := tx.Model(&database.Borrow{}).Where("card_id = ? and book_id = ? and return_time = 0", borrow.Card_id, borrow.Book_id).Update("return_time", time.Now().UnixMilli()).Error; err != nil {
			return err
		}

		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.Book_id).Update("stock", gorm.Expr("stock + 1")).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to return book",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Book returned successfully",
		Payload: nil,
	}
}

// ShowBorrowHistories
// list all borrow histories for a specific card.
// the returned records should be sorted by borrow_time DESC, book_id ASC
//
// @param cardId show which card's borrow history
// @return query results should be returned by ApiResult.payload
//
//	and should be an instance of {@link queries.BorrowHistories}
func ShowBorrowHistories(card_id int) APIResult {
	history := queries.BorrowHistories{}
	err := database.DB.Model(&database.Borrow{}).Joins("natural join book").Where("borrows.card_id = ?", card_id).Scan(&history.Items)
	if err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to fetch borrow histories",
			Payload: nil,
		}
	}
	history.Count = len(history.Items)
	return APIResult{
		Ok:      true,
		Message: "Borrow histories fetched successfully",
		Payload: history,
	}
}

// RegisterCard
// create a new borrow card. do nothing and return failed if
// the card already exists.
//
// Note that card_id should be stored to card after successfully
// completing this operation.
//
// @param card all attributes of the card
func RegisterCard(card database.Card) APIResult {
	return APIResult{}
}

// RemoveCard
// simply remove a card.
//
// Note that if there exists any un-returned books under this user,
// this card should not be removed.
//
// @param cardId card to be removed
func RemoveCard(card_id int) APIResult {
	return APIResult{}
}

// ShowCards
// list all cards order by card_id.
//
// @return query results should be returned by ApiResult.payload
//
//	and should be an instance of {@link queries.CardList}
func ShowCards() APIResult {
	cards := queries.CardList{}
	result := database.DB.Find(&cards.Cards)
	if result.Error != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to fetch cards",
			Payload: result.Error,
		}
	}
	cards.Count = int(result.RowsAffected)
	return APIResult{
		Ok:      true,
		Message: "Cards fetched successfully",
		Payload: cards,
	}
}
