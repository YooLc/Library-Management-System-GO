package queries

import (
	"fmt"
	"library-management-system/database"
)

type BorrowHistoryItem struct {
	database.Book
	Borrow_time int64 `json:"borrow_time"`
	Return_time int64 `json:"return_time"`
}

type BorrowHistories struct {
	Count int                 `json:"cnt"`
	Items []BorrowHistoryItem `json:"items"`
}

func (item BorrowHistoryItem) String() string {
	return fmt.Sprintf("BorrowHistory{Book: %s, Borrow Time: %v, Return Time: %v}", item.Book.Title, item.Borrow_time, item.Return_time)
}
