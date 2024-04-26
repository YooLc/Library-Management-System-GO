package queries

import (
	"cmp"
	"fmt"
	"library-management-system/database"
)

type Order string
type SortColumn string

const (
	Asc  Order = "asc"
	Desc       = "desc"
)
const (
	BookId      SortColumn = "book_id"
	Category               = "category"
	Title                  = "title"
	Press                  = "press"
	PublishYear            = "publish_year"
	Author                 = "author"
	Price                  = "price"
	Stock                  = "stock"
)

var SortColumns = []SortColumn{BookId, Category, Title, Press, PublishYear, Author, Price, Stock}
var SortOrders = []Order{Asc, Desc}

// BookQueryConditions
//
// Note: (1) all non-null attributes should be used as query
//
//	conditions and connected by "AND" operations.
//	(2) for range query of an attribute, the maximum and
//	minimum values use closed intervals.
//	eg: minA=x, maxA=y ==> x <= A <= y
//	    minA=null, maxA=y ==> A <= y
//	    minA=x, maxA=null ==> A >= x
type BookQueryConditions struct {
	Category       string     `json:"category"` /* Note: use fuzzy matching */
	Title          string     `json:"title"`    /* Note: use fuzzy matching */
	Press          string     `json:"press"`    /* Note: use fuzzy matching */
	MinPublishYear int        `json:"minPublishYear"`
	MaxPublishYear int        `json:"maxPublishYear"`
	Author         string     `json:"author"` /* Note: use fuzzy matching */
	MinPrice       float64    `json:"minPrice"`
	MaxPrice       float64    `json:"maxPrice"`
	SortBy         SortColumn `json:"sortBy"`    /* sort by which field */
	SortOrder      Order      `json:"sortOrder"` /* default sort by Primary Key */
}

func (c BookQueryConditions) String() string {
	return fmt.Sprintf("BookQueryConditions{Category: `%s`, Title: `%s`, Press: `%s`,"+
		"MinPublishYear: `%d`, MaxPublishYear: `%d`,"+
		"Author: `%s`, MinPrice: `%f`, MaxPrice: `%f`, SortBy: `%s`, SortOrder: `%s`}",
		c.Category, c.Title, c.Press, c.MinPublishYear, c.MaxPublishYear, c.Author, c.MinPrice, c.MaxPrice, c.SortBy, c.SortOrder)
}

func BookIdCmp(a, b *database.Book) int {
	return a.BookId - b.BookId
}
func CategoryCmp(a, b *database.Book) int {
	return cmp.Compare(a.Category, b.Category)
}
func TitleCmp(a, b *database.Book) int {
	return cmp.Compare(a.Title, b.Title)
}
func PressCmp(a, b *database.Book) int {
	return cmp.Compare(a.Press, b.Press)
}
func PublishYearCmp(a, b *database.Book) int {
	return a.PublishYear - b.PublishYear
}
func AuthorCmp(a, b *database.Book) int {
	return cmp.Compare(a.Author, b.Author)
}
func PriceCmp(a, b *database.Book) int {
	if a.Price < b.Price {
		return -1
	} else if a.Price > b.Price {
		return 1
	}
	return 0
}
func StockCmp(a, b *database.Book) int {
	return a.Stock - b.Stock
}

type BookComparator func(a, b *database.Book) int

func (Cmp BookComparator) Reverse() BookComparator {
	return func(a, b *database.Book) int {
		return -Cmp(a, b)
	}
}
func (Cmp BookComparator) ThenByIdAsc() BookComparator {
	return func(a, b *database.Book) int {
		res := Cmp(a, b)
		if res == 0 {
			return BookIdCmp(a, b)
		}
		return res
	}
}

func GetComparator(sortBy SortColumn) func(a, b *database.Book) int {
	switch sortBy {
	case BookId:
		return BookIdCmp
	case Category:
		return CategoryCmp
	case Title:
		return TitleCmp
	case Press:
		return PressCmp
	case PublishYear:
		return PublishYearCmp
	case Author:
		return AuthorCmp
	case Price:
		return PriceCmp
	case Stock:
		return StockCmp
	}
	return nil
}
