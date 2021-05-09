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

func GetBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	params := mux.Vars(r)

	rows := driver.DB.QueryRow("select * from books where id= $1", params["id"]) // parami v db so $1

	err := rows.Scan(&book.ID, &book.Title, &book.Availability) // damo sam tega specificnega v row pa izpisemo
	if err != nil {
		log.Println(err)
		return
	}

	if err = json.NewEncoder(w).Encode(book); err != nil {
		log.Println(err)
		return
	}

}

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
		err := rows.Scan(&userBook.ID, &userBook.FirstName, &userBook.LastName, &userBook.BookId, &userBook.Book) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		LogFatal(err)

		userBooks = append(userBooks, userBook)
	}

	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		log.Println(err)
		return
	}

}

func GetUserBooks(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, "admin-user-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}

}

func GetBooksApi(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	books = []models.Book{}

	rows, err := driver.DB.Query(`select * from books
							ORDER BY ID ASC`) // parami v db so $1
	LogFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Availability) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		LogFatal(err)

		books = append(books, book)
	}

	if err = json.NewEncoder(w).Encode(books); err != nil {
		LogFatal(err)
	}

}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-Books.page.tmpl")

	if err != nil {
		log.Println(err)
	}

}

func GetUsersApi(w http.ResponseWriter, r *http.Request) {
	var user models.User
	users = []models.User{}

	rows, err := driver.DB.Query(`select u.id, u.first_name , u.last_name, coalesce(b.count, 0) as borrowed_count 
						from users u left join (select user_id, count(user_id) 
						from borrowedbooks group by user_id) as b on b.user_id = u.id order by u.id asc`)
	LogFatal(err)

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BorrowedCount) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		LogFatal(err)

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(users); err != nil {
		LogFatal(err)
	}

}

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
		err := rows.Scan(&userBook.BookId, &userBook.Book, &userBook.ID, &userBook.FirstName, &userBook.LastName) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		LogFatal(err)

		userBooks = append(userBooks, userBook)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		LogFatal(err)
	}

}

func GetAllBorrowedBooks(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-all-borrowed-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-main.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

func AddUserPage(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-add-user.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	//var UserID int

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
		err := render.Template(w, "admin-main.page.tmpl")

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

// vrnemo se v sql RETURNING id, ki samo pove kateri ID nam pripada, to si zapisemo skozi scan, in pol ga napisemo kot response POSTU,
// tok da vemo kaj smo vstavil
// pa idje itak sam vstavlja

//json.NewEncoder(w).Encode(UserID)

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
	err = driver.DB.QueryRow("DELETE from users where id =  $1 RETURNING id", // placeholderje za valueje k jih postamo,
		user.ID).Scan(&UserID)

	if err != nil {
		log.Println(err)
	}

	// vrnemo se v sql RETURNING id, ki samo pove kateri ID nam pripada, to si zapisemo skozi scan, in pol ga napisemo kot response POSTU,
	// tok da vemo kaj smo vstavil
	// pa idje itak sam vstavlja

	json.NewEncoder(w).Encode(UserID)
}

func BorrowBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook models.BorrowBook
	//var Bookid int

	json.NewDecoder(r.Body).Decode(&borrowBook)

	//	_, err := db.Exec(`
	//
	//
	//`)

	//if err != nil {
	//	log.Println(err)
	//	return
	//}

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

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook models.BorrowBook
	var Bookid int

	json.NewDecoder(r.Body).Decode(&borrowBook)
	err := driver.DB.QueryRow("DELETE FROM borrowedbooks WHERE user_id = $1 AND book_id = $2 RETURNING id", // placeholderje za valueje k jih postamo,
		borrowBook.UserID, borrowBook.BookID).Scan(&Bookid)

	if err != nil {
		_, err = w.Write([]byte("No such user/book relationship"))
		LogFatal(err)
		return
	} else {
		json.NewEncoder(w).Encode("Success!")
	}

}
