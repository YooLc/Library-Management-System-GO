package server

import (
	"fmt"
	"library-management-system/database"
	"library-management-system/server/queries"
	"time"

	"gorm.io/gorm"
)

type Server struct{}

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
//	     (1) BookID should be stored to book after successfully
//	         completing this operation.
//	     (2) you should not register this book if the book already
//	         exists in the library system.
//
//	@param book all attributes of the book
func (s *Server) StoreBook(book *database.Book) database.APIResult {
	// Store the book
	// BookID is set via gorm
	// the database prevents duplicate book entries by primary key constraint
	if err := database.DB.Create(book).Error; err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to store book, maybe the book already exists",
			Payload: err,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Book stored successfully",
		Payload: book.BookId,
	}
}

// IncBookStock
// increase the book's inventory by bookId & deltaStock.
//
// Note that:
//
//	(1) you need to check the correctness of BookID
//	(2) deltaStock can be negative, but make sure that
//	    the result of book.stock + deltaStock is not negative!
//
// @param bookId book's BookID
// @param deltaStock increase count to book's stock, must be greater
func (s *Server) IncBookStock(bookId int, deltaStock int) database.APIResult {
	// Check the correctness of BookID
	book := database.Book{}
	if err := database.DB.First(&book, bookId).Error; err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "This book does not exist",
			Payload: nil,
		}
	}

	// Check the result of book.stock+deltaStock is not negative
	if book.Stock+deltaStock < 0 {
		return database.APIResult{
			Ok:      false,
			Message: "Stock deltaStock becomes invalid after incrementing, please check the arguments",
			Payload: nil,
		}
	}

	// Performing the increment operation
	// By default, gorm perform write (create/update/delete) operations
	// run inside a transaction to ensure data consistency
	if err := database.DB.Model(&book).Update("stock", book.Stock+deltaStock).Error; err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to increment book stock",
			Payload: nil,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Book stock incremented successfully",
		Payload: nil,
	}
}

// StoreBooks
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
func (s *Server) StoreBooks(books []*database.Book) database.APIResult {
	// Batch store books via transaction in gorm
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Add creation of each book to the transaction
		for _, book := range books {
			if err := tx.Create(book).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		// Restore BookId, which is assigned by gorm before inserting into the database
		// Even if the transaction fails, the BookId is still assigned
		for _, book := range books {
			book.BookId = 0
		}
		return database.APIResult{
			Ok:      false,
			Message: "Failed to store books, maybe one of them already exists",
			Payload: err,
		}
	}
	return database.APIResult{
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
func (s *Server) RemoveBook(bookId int) database.APIResult {
	// Check if someone has not returned this book
	var count int64
	database.DB.Model(&database.Borrow{}).Where("book_id = ? and return_time = 0", bookId).Count(&count)
	if count > 0 {
		return database.APIResult{
			Ok:      false,
			Message: "This book has some un-returned copies",
			Payload: nil,
		}
	}

	// Remove the book
	result := database.DB.Delete(&database.Book{}, bookId)
	if result.Error != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to remove book",
			Payload: result.Error,
		}
	} else if result.RowsAffected == 0 { // Delete will succeed even if the book does not exist
		return database.APIResult{
			Ok:      false,
			Message: "This book does not exist, maybe it was already removed",
			Payload: nil,
		}
	}

	return database.APIResult{
		Ok:      true,
		Message: "Book removed successfully",
		Payload: nil,
	}
}

// ModifyBookInfo
// modify a book's information by BookID.BookID.
//
// Note that you should not modify its BookID and stock!
//
// @param book the book to be modified
func (s *Server) ModifyBookInfo(book *database.Book) database.APIResult {
	// Avoid modifying BookID and stock
	origBook := database.Book{}
	//println(origBook.BookId, origBook.Title)
	if err := database.DB.First(&origBook, book.BookId).Error; err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "That book that does not exist, you cannot modify book_id",
			Payload: nil,
		}
	}

	// Modify the book info
	if err := database.DB.Model(book).Omit("book_id", "stock").Updates(book).Error; err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to modify book info",
			Payload: err,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Book info modified successfully",
		Payload: book,
	}
}

// QueryBooks
// Note that:
//
//	(1) you should let the DBMS filter records
//	    that do not satisfy the conditions instead of
//	    filter records in your API.
//	(2) when binding params to SQL, you also need to avoid
//	    the risk of SQL injection attack.
//	(3) [*] if all else is equal, sort by BookID in
//	    ascending order!
//
// @param conditions query conditions
//
// @return query results should be returned by database.APIResult.payload
//
//	and should be an instance of {@link queries.BookQueryResults}
func (s *Server) QueryBooks(conditions queries.BookQueryConditions) database.APIResult {
	// Query books
	books := queries.BookQueryResults{}

	query := database.DB.Model(&database.Book{})
	if conditions.Category != "" {
		query = query.Where("category like ?", "%"+conditions.Category+"%")
	}
	if conditions.Title != "" {
		query = query.Where("title like ?", "%"+conditions.Title+"%")
	}
	if conditions.Author != "" {
		query = query.Where("author like ?", "%"+conditions.Author+"%")
	}
	if conditions.MinPublishYear != 0 {
		query = query.Where("publish_year >= ?", conditions.MinPublishYear)
	}
	if conditions.MaxPublishYear != 0 {
		query = query.Where("publish_year <= ?", conditions.MaxPublishYear)
	}
	if conditions.MinPrice != 0 {
		query = query.Where("price >= ?", conditions.MinPrice)
	}
	if conditions.MaxPrice != 0 {
		query = query.Where("price <= ?", conditions.MaxPrice)
	}
	// TODO: handle order in conditions
	result := query.Order("book_id asc").Scan(&books.Results)

	if result.Error != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to query books",
			Payload: result.Error,
		}
	}
	books.Count = int(result.RowsAffected)
	return database.APIResult{
		Ok:      true,
		Message: "Books queried successfully",
		Payload: books,
	}
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
func (s *Server) BorrowBook(borrow database.Borrow) database.APIResult {
	borrow.BorrowTime = time.Now().UnixMilli()
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// Check if there are enough books in stock
		var stock int
		err := tx.Model(&database.Book{}).Select("stock").Where("book_id = ?", borrow.BookId).Row().Scan(&stock)
		if err != nil {
			return err
		}
		if stock <= 0 {
			return fmt.Errorf("book out of stock")
		}

		// Check if the user has not borrowed the book or has returned it
		var count int64
		tx.Model(&database.Borrow{}).Where("card_id = ? and book_id = ? and return_time = 0", borrow.CardId, borrow.BookId).Count(&count)
		if count > 0 { // There is a borrow record without return time
			return fmt.Errorf("user has not returned the book")
		}

		/* Check is OKay */
		// Borrow the book
		if err := tx.Create(&borrow).Error; err != nil {
			return err
		}
		// Update the stock of the book
		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.BookId).Update("stock", gorm.Expr("stock - 1")).Error; err != nil {
			return err
		}

		// Return nil to commit the transaction
		return nil
	})

	// If transaction failed, return error
	if err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to borrow book, maybe the user haven't returned the book or the book is out of stock",
			Payload: err,
		}
	}
	return database.APIResult{
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
func (s *Server) ReturnBook(borrow database.Borrow) database.APIResult {
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// return_time = 0 because a book can be borrowed
		// multiple times by the same card (but not the same time)
		if err := tx.Model(&database.Borrow{}).Where("card_id = ? and book_id = ? and return_time = 0", borrow.CardId, borrow.BookId).Update("return_time", time.Now().UnixMilli()).Error; err != nil {
			return err
		}

		// Update the stock of the book
		if err := tx.Model(&database.Book{}).Where("book_id = ?", borrow.BookId).Update("stock", gorm.Expr("stock + 1")).Error; err != nil {
			return err
		}
		return nil
	})

	// If transaction failed, return error
	if err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to return book",
			Payload: err,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Book returned successfully",
		Payload: nil,
	}
}

// ShowBorrowHistories
// list all borrow histories for a specific card.
// the returned records should be sorted by borrowTime DESC, bookId ASC
//
// @param cardId show which card's borrow history
// @return query results should be returned by database.APIResult.payload
//
//	and should be an instance of {@link queries.BorrowHistories}
func (s *Server) ShowBorrowHistories(cardId int) database.APIResult {
	history := queries.BorrowHistories{}
	err := database.DB.Model(&database.Borrow{}).Joins("natural join book").Where("borrows.CardId = ?", cardId).Scan(&history.Items)
	if err != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to fetch borrow histories",
			Payload: nil,
		}
	}
	history.Count = len(history.Items)
	return database.APIResult{
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
func (s *Server) RegisterCard(card *database.Card) database.APIResult {
	// Create a new borrow card
	if err := database.DB.Create(card).Error; err != nil {
		// Check if the card already exists by primary key
		return database.APIResult{
			Ok:      false,
			Message: "This card already exists",
			Payload: nil,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Card registered successfully",
		Payload: card.CardId,
	}
}

// RemoveCard
// simply remove a card.
//
// Note that if there exists any un-returned books under this user,
// this card should not be removed.
//
// @param cardId card to be removed
func (s *Server) RemoveCard(cardId int) database.APIResult {
	// Check if there exists any un-returned books under this user
	var count int64
	database.DB.Model(&database.Borrow{}).Where("card_id = ? and return_time = 0", cardId).Count(&count)
	if count > 0 {
		return database.APIResult{
			Ok:      false,
			Message: "This user has un-returned books",
			Payload: nil,
		}
	}

	// Remove the card
	result := database.DB.Delete(&database.Card{}, cardId)
	if result.Error != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to remove card",
			Payload: result.Error,
		}
	} else if result.RowsAffected == 0 {
		return database.APIResult{
			Ok:      false,
			Message: "This card does not exist, maybe it was already removed",
			Payload: nil,
		}
	}
	return database.APIResult{
		Ok:      true,
		Message: "Card removed successfully",
		Payload: nil,
	}
}

// ShowCards
// list all cards order by card_id.
//
// @return query results should be returned by database.APIResult.payload
//
//	and should be an instance of {@link queries.CardList}
func (s *Server) ShowCards() database.APIResult {
	cards := queries.CardList{}
	result := database.DB.Find(&cards.Cards)
	if result.Error != nil {
		return database.APIResult{
			Ok:      false,
			Message: "Failed to fetch cards",
			Payload: result.Error,
		}
	}
	cards.Count = int(result.RowsAffected)
	return database.APIResult{
		Ok:      true,
		Message: "Cards fetched successfully",
		Payload: cards,
	}
}
