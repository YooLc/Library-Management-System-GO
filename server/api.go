package server

type APIResult struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

type Order int

const (
	Asc Order = iota
	Desc
)

type BookQueryConditions struct {
	Category       string  `json:"category"`
	Title          string  `json:"title"`
	Press          string  `json:"press"`
	MinPublishYear int     `json:"minPublishYear"`
	MaxPublishYear int     `json:"maxPublishYear"`
	Author         string  `json:"author"`
	MinPrice       float64 `json:"minPrice"`
	MaxPrice       float64 `json:"maxPrice"`
	SortBy         string  `json:"sortBy"`
	SortOrder      Order   `json:"sortOrder"`
}
