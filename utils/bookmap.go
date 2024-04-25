package utils

import "library-management-system/database"

type BookMap map[database.BookKey]*database.Book

// Insert a book into the map
func (b BookMap) Insert(val database.Book) {
	key := database.BookKey{
		Category:    val.Category,
		Title:       val.Title,
		Press:       val.Press,
		PublishYear: val.PublishYear,
		Author:      val.Author,
	}
	_, exists := b[key]
	if !exists {
		b[key] = &val
	}
}

// List return a list of books
func (b BookMap) List() []*database.Book {
	var list []*database.Book
	for _, v := range b {
		list = append(list, v)
	}
	return list
}
