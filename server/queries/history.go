package queries

import (
	"library-management-system/database"
)

type BorrowHistories struct {
	Count int               `json:"cnt"`
	Items []database.Borrow `json:"items"`
}
