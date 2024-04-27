package queries

import "library-management-system/database"

type BookQueryResults struct {
	Count   int             `json:"count"`
	Results []database.Book `json:"results"`
}

type BorrowHistories struct {
	Count int               `json:"count"`
	Items []database.Borrow `json:"items"`
}
