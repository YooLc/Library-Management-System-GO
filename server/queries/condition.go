package queries

type Order int
type SortColumn int

const (
	Asc Order = iota
	Desc
)
const (
	Book_id SortColumn = iota
	Category
	Title
	Press
	Publish_year
	Author
	Price
	Stock
)

type BookQueryConditions struct {
	Category       string     `json:"category"`
	Title          string     `json:"title"`
	Press          string     `json:"press"`
	MinPublishYear int        `json:"minPublishYear"`
	MaxPublishYear int        `json:"maxPublishYear"`
	Author         string     `json:"author"`
	MinPrice       float64    `json:"minPrice"`
	MaxPrice       float64    `json:"maxPrice"`
	SortBy         SortColumn `json:"sortBy"`
	SortOrder      Order      `json:"sortOrder"`
}
