package queries

import (
	"library-management-system/database"
)

type BorrowHistories struct {
	Count int               `json:"count"`
	Items []database.Borrow `json:"items"`
}
