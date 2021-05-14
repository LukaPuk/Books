package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/LukaPuk/Books/internal/driver"
	"github.com/LukaPuk/Books/internal/models"
	"github.com/LukaPuk/Books/internal/render"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

var users []models.User

var books []models.Book

var userBooks []models.UserBooks

func Homepage(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, "admin-homepage.page.tmpl")

	if err != nil {
		log.Println(err)
	}

}

func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

//GetBook Gets single Book info(title, Availability) based on Book ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)

	rows := driver.DB.QueryRow("select * from books where id= $1", params["id"])

	err := rows.Scan(&book.ID, &book.Title, &book.Availability)
	if err != nil {
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(book); err != nil {
		log.Println(err)
		return
	}

}

//GetUserBooksApi Gets Json with all borrowed books by User
func GetUserBooksApi(w http.ResponseWriter, r *http.Request) {
	var userBook models.UserBooks
	userBooks = []models.UserBooks{}

	params := mux.Vars(r)

	rows, err := driver.DB.Query(`select u.*, b.book_id as bookid, bo.title from borrowedbooks b 
						left join users u on u.id = b.user_id 
						left join books bo on bo.id = b.book_id 
						where b.user_id = $1`, params["id"]) // parami v db so $1
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		err := rows.Scan(&userBook.ID, &userBook.FirstName, &userBook.LastName, &userBook.BookId, &userBook.Book)
		LogFatal(err)

		userBooks = append(userBooks, userBook)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		log.Println(err)
		return
	}

}

//GetUserBooks Displays webpage with all borrowed books by User
func GetUserBooks(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, "admin-user-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}

}

//GetBooksApi Gets JSON with All books
func GetBooksApi(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	books = []models.Book{}

	rows, err := driver.DB.Query(`select * from books
							ORDER BY ID ASC`)
	LogFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Availability)
		LogFatal(err)

		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(books); err != nil {
		LogFatal(err)
	}

}

//GetBooks Displays All books on a page
func GetBooks(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-Books.page.tmpl")

	if err != nil {
		log.Println(err)
	}

}



//GetAllBorrowedBooksApi Gets JSON with All borrowed books
func GetAllBorrowedBooksApi(w http.ResponseWriter, r *http.Request) {
	var userBook models.UserBooks
	userBooks = []models.UserBooks{}

	rows, err := driver.DB.Query(`select b.book_id as bookid, bo.title, u.* from borrowedbooks b 
						left join users u on u.id = b.user_id 
						left join books bo on bo.id = b.book_id
						ORDER BY b.book_id asc;
						`)

	if err != nil {
		log.Print(err)
	}
	LogFatal(err)

	for rows.Next() {
		err := rows.Scan(&userBook.BookId, &userBook.Book, &userBook.ID, &userBook.FirstName, &userBook.LastName)
		LogFatal(err)

		userBooks = append(userBooks, userBook)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		LogFatal(err)
	}

}

//GetAllBorrowedBooks Displays page with all borrowed books
func GetAllBorrowedBooks(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-all-borrowed-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

//GetUsersApi Gets JSON with All users
func GetUsersApi(w http.ResponseWriter, r *http.Request) {
	var user models.User
	users = []models.User{}

	rows, err := driver.DB.Query(`select u.id, u.first_name , u.last_name, coalesce(b.count, 0) as borrowed_count 
						from users u left join (select user_id, count(user_id) 
						from borrowedbooks group by user_id) as b on b.user_id = u.id order by u.id asc`)
	LogFatal(err)

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BorrowedCount)
		LogFatal(err)

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(users); err != nil {
		LogFatal(err)
	}

}

//GetUsers Displays page with All users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-users.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

//AddUserPage Displays Add user page
func AddUserPage(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-add-user.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

//AddUser adds user to database
func AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := r.PostForm
	if len(form) < 1 {
		json.NewDecoder(r.Body).Decode(&user)
	} else {
		user.FirstName = form["first_name"][0]
		user.LastName = form["last_name"][0]
	}

	//

	row := driver.DB.QueryRow("select * from adduser($1,$2)", user.FirstName, user.LastName)

	if row.Err() != nil {
		_, err = w.Write([]byte(row.Err().Error()))

		if err != nil {
			log.Println(err)
		}
		return
	}

	if len(form) > 0 {
		err := render.Template(w, "admin-users.page.tmpl")

		if err != nil {
			log.Println(err)
		}
	} else {
		_, err = w.Write([]byte("Success!"))

		if err != nil {
			log.Println(err)

		}
	}

}


//DeleteUser Deletes User
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var UserID int
	json.NewDecoder(r.Body).Decode(&user)

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	form := r.PostForm
	fmt.Println(form)

	json.NewDecoder(r.Body).Decode(&user)
	err = driver.DB.QueryRow("DELETE from users where id =  $1 RETURNING id",
		user.ID).Scan(&UserID)

	if err != nil {
		log.Println(err)
		return
	}



	json.NewEncoder(w).Encode(UserID)
}

//BorrowBook Borrows Book by User
func BorrowBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook models.BorrowBook

	json.NewDecoder(r.Body).Decode(&borrowBook)



	row := driver.DB.QueryRow("select * from borrowbook($1,$2)", borrowBook.UserID, borrowBook.BookID)

	if row.Err() != nil {
		if strings.Contains(row.Err().Error(), "userid") {
			_, err := w.Write([]byte("User doesn't exist"))

			if err != nil {
				log.Println(err)
			}
		} else if strings.Contains(row.Err().Error(), "book_id") {
			_, err := w.Write([]byte("User already borrowed this book"))

			if err != nil {
				log.Println(err)
			}
		} else {
			_, err := w.Write([]byte(row.Err().Error()))

			if err != nil {
				log.Println(err)
			}

		}
	}

}

//ReturnBook Returns Book by User
func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook models.BorrowBook
	var Bookid int

	json.NewDecoder(r.Body).Decode(&borrowBook)
	err := driver.DB.QueryRow("DELETE FROM borrowedbooks WHERE user_id = $1 AND book_id = $2 RETURNING id",
		borrowBook.UserID, borrowBook.BookID).Scan(&Bookid)

	if err != nil {
		_, err = w.Write([]byte("No such user/book relationship"))
		LogFatal(err)
		return
	} else {
		json.NewEncoder(w).Encode("Success!")
	}

}
