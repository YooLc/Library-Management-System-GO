package utils

import (
	"library-management-system/database"
	"math/rand"
	"time"
)

var categories = []string{
	"Computer Science", "Nature", "Philosophy", "History", "Autobiography",
	"Magazine", "Dictionary", "Novel", "Horror", "Others",
}
var press = []string{
	"Press-A", "Press-B", "Press-C", "Press-D",
	"Press-E", "Press-F", "Press-G", "Press-H",
}
var authors = []string{
	"Nonehyo", "DouDou", "Coco", "Yuuku", "SoonWhy",
	"Fubuki", "Authentic", "Immortal", "ColaOtaku", "Erica", "ZaiZai", "DaDa", "Hgs",
}
var titles = []string{
	"Database System Concepts", "Computer Networking",
	"Algorithms", "Database System Designs", "Compiler Designs", "C++ Primer", "Operating System",
	"The Old Man and the Sea", "How steel is made", "Le Petit Prince", "The Metamorphosis",
	"Miserable World", "Gone with the wind", "Eugenie Grandet", "Analysis of Dreams",
}
var departments = []string{
	"Computer Science", "Law",
	"Management", "Civil Engineering", "Architecture", "Environmental Science",
	"English Language", "General Education", "Ideological & Political",
}

// time stamp of 2023-1-1, use package to convert to time
var timeStart, _ = time.Parse(time.RFC3339, "2017-01-01T00:00:00Z")
var timeEnd, _ = time.Parse(time.RFC3339, "2024-12-31T23:59:59Z")

func RandomBook() database.Book {
	return database.Book{
		Category:    randomCategory(),
		Title:       randomTitle(),
		Press:       randomPress(),
		PublishYear: randomPublishYear(),
		Author:      randomAuthor(),
		Price:       RandomPrice(),
		Stock:       RandomStock(),
	}
}

func randomCategory() string {
	return categories[rand.Intn(len(categories))]
}

func randomPress() string {
	return press[rand.Intn(len(press))]
}

func randomAuthor() string {
	return authors[rand.Intn(len(authors))]
}

func randomTitle() string {
	return titles[rand.Intn(len(titles))]
}

func randomPublishYear() int {
	return rand.Intn(25) + 2000
}

func RandomPrice() float64 {
	return float64(rand.Intn(10000)) / 100
}

func RandomStock() int {
	return rand.Intn(100) + 1
}

func randomDepartment() string {
	return departments[rand.Intn(len(departments))]
}

func randomTime() int64 {
	return rand.Int63n(timeEnd.Unix()-timeStart.Unix()) + timeStart.Unix()
}
