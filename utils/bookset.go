package utils

import "library-management-system/database"

type BookSet map[database.BookKey]*database.Book

func getKey(val database.Book) database.BookKey {
	return database.BookKey{
		Category:    val.Category,
		Title:       val.Title,
		Press:       val.Press,
		PublishYear: val.PublishYear,
	}
}

// Insert a book into the map
func (b BookSet) Insert(val database.Book) {
	key := getKey(val)
	if !b.Contains(val) {
		b[key] = &val
	}
}

// List return a list of books
func (b BookSet) List() []*database.Book {
	var list []*database.Book
	for _, v := range b {
		list = append(list, v)
	}
	return list
}

// Contains check if a book is in the set
func (b BookSet) Contains(val database.Book) bool {
	key := getKey(val)
	_, exists := b[key]
	return exists
}

// Remove a book from the set
func (b BookSet) Remove(val database.Book) bool {
	key := getKey(val)
	if !b.Contains(val) {
		return false
	}
	delete(b, key)
	return true
}

// Size return the size of the set
func (b BookSet) Size() int {
	return len(b)
}
