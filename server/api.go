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

/**
 * Note:
 *      (1) all functions in this interface will be regarded as a
 *          transaction. this means that after successfully completing
 *          all operations in a function, you need to call commit(),
 *          or call rollback() if one of the operations in a function fails.
 *          as an example, you can see {@link LibraryManagementSystemImpl#resetDatabase}
 *          to find how to use commit() and rollback().
 *      (2) for each function, you need to briefly introduce how to
 *          achieve this function and how to solve challenges in your
 *          lab report.
 *      (3) if you don't know what the function means, or what it is
 *          supposed to do, looking to the test code might help.
 */

/* Interface for books */

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
	// Store the book
	// book_id is set via gorm
	// the database prevents duplicate book entries by primary key constraint
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
func IncBookStock(bookId int, deltaStock int) APIResult {
	// Check the correctness of book_id
	book := database.Book{}
	if err := database.DB.First(&book, bookId).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "This book does not exist",
			Payload: nil,
		}
	}

	// Check the result of book.stock+deltaStock is not negative
	if book.Stock+deltaStock < 0 {
		return APIResult{
			Ok:      false,
			Message: "Stock deltaStock becomes invalid after incrementing, please check the arguments",
			Payload: nil,
		}
	}

	// Performing the increment operation
	// By default, gorm perform write (create/update/delete) operations
	// run inside a transaction to ensure data consistency
	if err := database.DB.Model(&book).Update("stock", book.Stock+deltaStock).Error; err != nil {
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
func StoreBookList(books []database.Book) APIResult {
	// Batch store books via transaction in gorm
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Add creation of each book to the transaction
		for _, book := range books {
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
		// Check if there are enough books in stock
		var stock int
		tx.Model(&database.Book{}).Select("stock").Where("book_id = ?", borrow.Book_id).Row().Scan(&stock)
		if stock <= 0 {
			return fmt.Errorf("book out of stock")
		}

		// Check if the user has not borrowed the book or has returned it
		var count int64
		tx.Model(&database.Borrow{}).Where("card_id = ? and book_id = ? and return_time = 0", borrow.Card_id, borrow.Book_id).Count(&count)
		if count > 0 { // There is a borrow record without return time
			return fmt.Errorf("user has not returned the book")
		}

		/* Check is OKay */
		// Borrow the book
		if err := tx.Create(&borrow).Error; err != nil {
			return err
		}
		// Update the stock of the book
		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.Book_id).Update("stock", gorm.Expr("stock - 1")).Error; err != nil {
			return err
		}

		// Return nil to commit the transaction
		return nil
	})

	// If transaction failed, return error
	if err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to borrow book, maybe the user haven't returned the book or the book is out of stock",
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
		// return_time = 0 because a book can be borrowed
		// multiple times by the same card (but not the same time)
		if err := tx.Model(&database.Borrow{}).Where("card_id = ? and book_id = ? and return_time = 0", borrow.Card_id, borrow.Book_id).Update("return_time", time.Now().UnixMilli()).Error; err != nil {
			return err
		}

		// Update the stock of the book
		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.Book_id).Update("stock", gorm.Expr("stock + 1")).Error; err != nil {
			return err
		}
		return nil
	})

	// If transaction failed, return error
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
func RemoveCard(cardId int) APIResult {
	// Check if there exists any un-returned books under this user
	var count int64
	database.DB.Model(&database.Borrow{}).Where("card_id = ? and return_time = 0", cardId).Count(&count)
	if count > 0 {
		return APIResult{
			Ok:      false,
			Message: "This user has un-returned books",
			Payload: nil,
		}
	}

	// Remove the card
	if err := database.DB.Delete(&database.Card{}, cardId).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to remove card",
			Payload: err,
		}
	}
	return APIResult{
		Ok:      true,
		Message: "Card removed successfully",
		Payload: nil,
	}
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
