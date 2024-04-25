package queries

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
