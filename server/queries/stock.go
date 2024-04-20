package queries

import "library-management-system/database"

type IncStockQuery struct {
	Book  database.Book `json:"book"`
	Count int           `json:"count"`
}
