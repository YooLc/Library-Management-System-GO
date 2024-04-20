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

func StoreBook(book database.Book) APIResult {
	if err := database.DB.Create(&book).Error; err != nil {
		return APIResult{
			Ok:      false,
			Message: "Failed to store book",
			Payload: err,
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
	if book.Stock+count < 0 {
		return APIResult{
			Ok:      false,
			Message: "Stock count becomes invalid after incrementing, please check the arguments",
			Payload: nil,
		}
	}
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

func storeBookList(books queries.BookList) APIResult {
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

func removeBook(book database.Book) APIResult {
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

func modifyBookInfo(book database.Book) APIResult {
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

func queryBook(conditions queries.BookQueryConditions) APIResult {
	return APIResult{}
}

func borrowBook(borrow database.Borrow) APIResult {
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

func returnBook(borrow database.Borrow) APIResult {
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

func showBorrowHistories(card_id int) APIResult {
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

func registerCard(card database.Card) APIResult {
	return APIResult{}
}

func removeCard(card_id int) APIResult {
	return APIResult{}
}

func showCards() APIResult {
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
