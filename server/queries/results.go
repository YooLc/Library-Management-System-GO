package queries

import "library-management-system/database"

type BookQueryResults struct {
	Count   int             `json:"count"`
	Results []database.Book `json:"results"`
}
