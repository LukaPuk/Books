package models


type Book struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Availability int    `json:"availability"`
}

type User struct {
	ID            int    `json:"id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	BorrowedCount int    `json:"borrowed_count"`
}

// All books by user
type UserBooks struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	BookId    int    `json:"bookid"`
	Book      string `json:"book"`
}


type BorrowBook struct {
	ID     int `json:"id"`
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}
