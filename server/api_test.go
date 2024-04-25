package server

import (
	"fmt"
	"library-management-system/database"
	"library-management-system/server/queries"
	"library-management-system/utils"
	"math/rand"
	"os"
	"sort"
	"testing"

	"github.com/go-playground/assert/v2"
	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	Server   Config                  `yaml:"server"`
	Database database.DatabaseConfig `yaml:"database"`
}

func TestMain(m *testing.M) {
	file, err := os.Open("../config.yaml")
	if err != nil {
		fmt.Println("Failed to open config file: ", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Failed to close config file: ", err)
		}
	}(file)

	var config AppConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		fmt.Println("Failed to parse config file: ", err)
		return
	}

	database.ConnectDatabase(config.Database)
	m.Run()
}

func TestBookRegister(t *testing.T) {
	database.ResetDatabase()
	b0 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", PublishYear: 2023,
		Author: "Mike", Price: 188.88, Stock: 10,
	}
	assert.Equal(t, StoreBook(&b0).Ok, true)

	/* Not allowed to create duplicated records */
	b1 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", PublishYear: 2023,
		Author: "Mike", Price: 188.88, Stock: 5,
	}
	b2 := database.Book{
		Category: "Computer Science", Title: "Database System Concepts",
		Press: "Machine Industry Press", PublishYear: 2023,
		Author: "Mike", Price: 99.99, Stock: 10,
	}
	assert.Equal(t, StoreBook(&b1).Ok, false)
	assert.Equal(t, StoreBook(&b2).Ok, false)
}

func TestIncBookStock(t *testing.T) {
	const numBooks = 50
	const numRandomTests = 1000
	/* simply insert some books to database */
	database.ResetDatabase()
	var books = make(map[database.Book]int) // Use map to avoid duplicated books
	var bookIds = make(map[int]int)
	for i := 0; i < numBooks; i++ {
		book := utils.RandomBook()
		books[book] = i
	}

	var bookList = make([]*database.Book, 0, len(books))
	for book, i := range books {
		result := StoreBook(&book)
		assert.Equal(t, result.Ok, true)
		bookIds[result.Payload.(int)] = i
		bookList = append(bookList, &book)
	}
	assert.Equal(t, len(books), len(bookIds))

	/* begin tests */
	type args struct {
		bookId     int
		deltaStock int
	}
	type test struct {
		name string
		args args
		want APIResult
	}
	var tests []test

	/* corner case: invalid book id */
	tests = append(tests, test{
		name: "Invalid book id - Negative",
		args: args{bookId: -1, deltaStock: 6},
		want: APIResult{Ok: false},
	})
	k := len(books) + 1
	_, ok := bookIds[k]
	for ok { // generate an invalid book id
		k++
		_, ok = bookIds[k]
	}
	tests = append(tests, test{
		name: "Invalid book id - Maximum",
		args: args{bookId: k, deltaStock: 10},
		want: APIResult{Ok: false},
	})

	/* corner case: invalid book stock */
	lastBook := bookList[len(bookList)-1]
	tests = append(tests, test{
		name: "Decrease book stock",
		args: args{bookId: lastBook.BookId, deltaStock: -lastBook.Stock},
		want: APIResult{Ok: true},
	})
	tests = append(tests, test{
		name: "Increase book stock",
		args: args{bookId: lastBook.BookId, deltaStock: 1},
		want: APIResult{Ok: true},
	})
	tests = append(tests, test{
		name: "Test for invalid book stock - Negative",
		args: args{bookId: lastBook.BookId, deltaStock: -2},
		want: APIResult{Ok: false},
	})

	/* randomly choose some books to do this operation */
	for i := 0; i < numRandomTests; i++ {
		book := bookList[rand.Intn(len(bookList)-1)]
		assert.NotEqual(t, book, nil)
		deltaStock := rand.Intn(24) - 8
		expected := book.Stock+deltaStock >= 0
		if expected {
			book.Stock = book.Stock + deltaStock
		}
		tests = append(tests, test{
			name: fmt.Sprintf("Random test %d", i),
			args: args{bookId: book.BookId, deltaStock: deltaStock},
			want: APIResult{Ok: expected},
		})
	}

	/* run tests */
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IncBookStock(tt.args.bookId, tt.args.deltaStock); got.Ok != tt.want.Ok {
				t.Errorf("IncBookStock() = %v, want %v", got, tt.want)
			}
		})
	}

	/* use query interface to check correctness */
	// TODO: Implement query
}

func TestBulkRegisterBook(t *testing.T) {
	const numBulkBooks = 1000
	const numDuplicateBooks = 3
	database.ResetDatabase()

	/* simply insert some books to database */
	var books = make(utils.BookMap)
	for i := 0; i < numBulkBooks; i++ {
		book := utils.RandomBook()
		books.Insert(book)
	}

	///* provide some duplicate records */
	bookList1 := books.List()
	for i := 0; i < numDuplicateBooks; i++ {
		newBook := *bookList1[rand.Intn(len(bookList1))]
		cb := &newBook
		// randomly change some attributes
		if rand.Intn(2) == 0 {
			cb.Stock = utils.RandomStock()
			cb.Price = utils.RandomPrice()
		}
		bookList1 = append(bookList1, cb)
	}
	// shuffle the book list
	rand.Shuffle(len(bookList1), func(i, j int) {
		bookList1[i], bookList1[j] = bookList1[j], bookList1[i]
	})
	assert.Equal(t, StoreBooks(bookList1).Ok, false)

	/* make sure that none of the books are inserted */
	queryResult1 := QueryBooks(queries.BookQueryConditions{})
	assert.Equal(t, queryResult1.Ok, true)
	selectedResults1 := queryResult1.Payload.(queries.BookQueryResults)
	assert.Equal(t, 0, selectedResults1.Count)

	/* normal batch insert */
	bookList2 := books.List()
	assert.Equal(t, StoreBooks(bookList2).Ok, true)
	queryResult2 := QueryBooks(queries.BookQueryConditions{})
	assert.Equal(t, queryResult2.Ok, true)
	selectedResults2 := queryResult2.Payload.(queries.BookQueryResults)
	assert.Equal(t, len(bookList2), selectedResults2.Count)
	sort.Slice(bookList2, func(i, j int) bool {
		return bookList2[i].BookId < bookList2[j].BookId
	})
	for i := 0; i < len(bookList2); i++ {
		assert.Equal(t, bookList2[i], selectedResults2.Results[i])
	}
}
