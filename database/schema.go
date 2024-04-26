package database

import (
	"fmt"
	"time"
)

type BookKey struct {
	Category    string // `json:"category" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Title       string // `json:"title" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Press       string // `json:"press" gorm:"size:63;not null;uniqueIndex:idx_book"`
	PublishYear int    // `json:"publish_year" gorm:"not null;uniqueIndex:idx_book"`
	Author      string // `json:"author" gorm:"size:63;not null;uniqueIndex:idx_book"`
}

type Book struct {
	BookId      int     `json:"book_id" gorm:"primaryKey;autoIncrement"`
	Category    string  `json:"category" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Title       string  `json:"title" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Press       string  `json:"press" gorm:"size:63;not null;uniqueIndex:idx_book"`
	PublishYear int     `json:"publish_year" gorm:"not null;uniqueIndex:idx_book"`
	Author      string  `json:"author" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Price       float64 `json:"price" gorm:"not null;type:decimal(7,2);default:0.00"`
	Stock       int     `json:"stock" gorm:"not null;default:0"`
	Borrow      Borrow  `gorm:"foreignKey:BookId;references:BookId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Card struct {
	CardId     int    `json:"card_id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"size:63;not null;uniqueIndex:idx_card"`
	Department string `json:"department" gorm:"size:63;not null;uniqueIndex:idx_card"`
	Type       string `json:"type" gorm:"type:char(1);not null;check:type in ('T', 'S');uniqueIndex:idx_card"`
	Borrow     Borrow `gorm:"foreignKey:CardId;references:CardId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Borrow struct {
	CardId     int   `gorm:"primaryKey"`
	BookId     int   `gorm:"primaryKey"`
	BorrowTime int64 `gorm:"primaryKey;not null"`
	ReturnTime int64 `gorm:"default:0"`
}

func (b *Borrow) ResetBorrowTime() {
	b.BorrowTime = time.Now().UnixMilli()
}
func (b *Borrow) ResetReturnTime() {
	b.ReturnTime = time.Now().UnixMilli()
}
func CreateBorrow(cardId, bookId int) Borrow {
	return Borrow{
		CardId:     cardId,
		BookId:     bookId,
		BorrowTime: time.Now().UnixMilli(),
		ReturnTime: 0,
	}
}

func (b *Book) String() string {
	return fmt.Sprintf("Book{BookId: %v, Category: %v, Title: %v, Press: %v, PublishYear: %v, Author: %v, Price: %v, Stock: %v}",
		b.BookId, b.Category, b.Title, b.Press, b.PublishYear, b.Author, b.Price, b.Stock)
}
func (c *Card) String() string {
	return fmt.Sprintf("Card{CardId: %v, Name: %v, Department: %v, Type: %v}",
		c.CardId, c.Name, c.Department, c.Type)
}
func (b *Borrow) String() string {
	return fmt.Sprintf("Borrow{CardId: %v, BookId: %v, BorrowTime: %v, ReturnTime: %v}",
		b.CardId, b.BookId, b.BorrowTime, b.ReturnTime)
}
