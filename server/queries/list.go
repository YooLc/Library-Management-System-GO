package queries

import "library-management-system/database"

type BookList struct {
	Count int             `json:"count"`
	Books []database.Book `json:"books"`
}

type CardList struct {
	Count int             `json:"count"`
	Cards []database.Card `json:"cards"`
}
