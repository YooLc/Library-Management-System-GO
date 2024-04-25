package utils

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"library-management-system/database"
	"math/rand"
	"sort"
)

type Library struct {
	Books   []*database.Book
	Cards   []*database.Card
	Borrows []database.Borrow
}

func CreateLibrary(nBooks int, nCards int, nBorrows int, server ServerInterface) Library {
	/* create hooks */
	bookSet := make(BookSet)
	for bookSet.Size() < nBooks {
		bookSet.Insert(RandomBook())
	}
	books := bookSet.List()
	assert.Equal(nil, server.StoreBooks(books).Ok, true)

	/* create cards */
	cards := make([]*database.Card, 0)
	for i := 1; i <= nCards; i++ {
		card := database.Card{}
		card.Name = fmt.Sprintf("User%05d", i)
		card.Department = RandomDepartment()
		card.Type = RandomCardType()
		cards = append(cards, &card)
		assert.Equal(nil, server.RegisterCard(&card).Ok, true)
	}

	/* create histories */
	borrows := make([]database.Borrow, 0)
	timeStamps := make([]int64, 0)
	for i := 0; i < nBorrows*2; i++ {
		timeStamps = append(timeStamps, RandomTime())
	}
	sort.Slice(timeStamps, func(i, j int) bool {
		return timeStamps[i] < timeStamps[j]
	})

	count := 0
	for count < nBorrows {
		book := books[rand.Intn(len(books))]
		if book.Stock == 0 {
			continue
		}
		card := cards[rand.Intn(len(cards))]
		borrow := database.Borrow{
			BookId:     book.BookId,
			CardId:     card.CardId,
			BorrowTime: timeStamps[count*2],
			ReturnTime: timeStamps[count*2+1],
		}
		assert.Equal(nil, server.BorrowBook(borrow).Ok, true)
		assert.Equal(nil, server.ReturnBook(borrow).Ok, true)
		borrows = append(borrows, borrow)
		count++
	}
	return Library{
		Books:   books,
		Cards:   cards,
		Borrows: borrows,
	}
}

func (lib Library) NumBooks() int {
	return len(lib.Books)
}

func (lib Library) NumCards() int {
	return len(lib.Cards)
}

func (lib Library) NumBorrows() int {
	return len(lib.Borrows)
}
