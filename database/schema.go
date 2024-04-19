package database

type Book struct {
	Book_id      int     `json:"book_id" gorm:"primaryKey;autoIncrement"`
	Category     string  `json:"category" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Title        string  `json:"title" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Press        string  `json:"press" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Publish_year int     `json:"publish_year" gorm:"not null;uniqueIndex:idx_book"`
	Author       string  `json:"author" gorm:"size:63;not null;uniqueIndex:idx_book"`
	Price        float64 `json:"price" gorm:"not null;type:decimal(7,2);default:0.00"`
	Stock        int     `json:"stock" gorm:"not null;default:0"`
	Borrow       Borrow  `gorm:"foreignKey:Book_id;references:Book_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Card struct {
	Card_id    int    `json:"card_id" gorm:"primaryKey;autoIncrement"`
	Name       string `json:"name" gorm:"size:63;not null;uniqueIndex:idx_card"`
	Department string `json:"department" gorm:"size:63;not null;uniqueIndex:idx_card"`
	Type       string `json:"type" gorm:"type:char(1);not null;check:type in ('T', 'S');uniqueIndex:idx_card"`
	Borow      Borrow `gorm:"foreignKey:Card_id;references:Card_id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Borrow struct {
	Card_id     int   `gorm:"primaryKey"`
	Book_id     int   `gorm:"primaryKey"`
	Borrow_time int64 `gorm:"primaryKey;not null"`
	Return_time int64 `gorm:"default:0"`
}
