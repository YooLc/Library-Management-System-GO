package queries

import (
	"fmt"
	"library-management-system/database"
)

type BorrowHistories struct {
	database.Book
	BorrowTime int64 `json:"borrowTime"`
	ReturnTime int64 `json:"returnTime"`
}

func (b BorrowHistories) String() string {
	return "BorrowHistories{ Book , BorrowTime: " + fmt.Sprint(b.BorrowTime) + ", ReturnTime: " + fmt.Sprint(b.ReturnTime) + "}"
}
