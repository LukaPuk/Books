package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/LukaPuk/Books/render"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Book struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Availability int `json:"availability"`
}


type User struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	BorrowedCount int `json:"borrowed_count"`
}


type UserBooks struct {
	ID int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	BookId int `json:"bookid"`
	Book string `json:"book"`
}

type BorrowBook struct {
	ID int `json:"id"`
	UserID int `json:"user_id"`
	BookID int `json:"book_id"`
}

var users []User

var books []Book

var userBooks []UserBooks

var db *sql.DB


func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	gotenv.Load() // loads env file

}

// add neccessary json format or error

func main() {
	ElephantUrl := os.Getenv("ELEPHANTSQL_URL")
	pgUrl, err := pq.ParseURL(ElephantUrl)
	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	err = db.Ping()
	logFatal(err)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(0)
	db.SetConnMaxLifetime(time.Millisecond)
	defer db.Close()

	//_,  err = render.CreateTemplateCache()
	//if err != nil {
	//	log.Fatal("cannot create template cache")
	//}



	router := mux.NewRouter()
	router.HandleFunc("/", Homepage).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/booksapi", getBooksApi).Methods("GET")
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/usersapi", getUsersApi).Methods("GET")
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/add-user", addUser).Methods("POST")
	router.HandleFunc("/add-user", addUserPage).Methods("GET")
	router.HandleFunc("/borrowbook", borrowBook).Methods("POST")
	router.HandleFunc("/returnbook", returnBook).Methods("POST")
	router.HandleFunc("/delete", deleteUser).Methods("POST")
	router.HandleFunc("/userapi/{id}", getUserBooksApi).Methods("GET")
	router.HandleFunc("/user/{id}", getUserBooks).Methods("GET")
	router.HandleFunc("/all-borrowed-booksapi", getAllBorrowedBooksApi).Methods("GET")
	router.HandleFunc("/all-borrowed-books", getAllBorrowedBooks).Methods("GET")


	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}

}

func Homepage(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, "admin-homepage.page.tmpl")

	if err != nil {
		log.Println(err)
	}


}



func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.Vars(r)

	rows := db.QueryRow("select * from books where id= $1", params["id"]) // parami v db so $1

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

func getUserBooksApi(w http.ResponseWriter, r *http.Request) {
	var userBook UserBooks
	userBooks = []UserBooks{}

	params := mux.Vars(r)

	rows, err := db.Query(`select u.*, b.book_id as bookid, bo.title from borrowedbooks b 
						left join users u on u.id = b.user_id 
						left join books bo on bo.id = b.book_id 
						where b.user_id = $1`, params["id"]) // parami v db so $1
	if err != nil {
		log.Print(err)
	}



	for rows.Next() {
		err := rows.Scan(&userBook.ID, &userBook.FirstName, &userBook.LastName, &userBook.BookId, &userBook.Book) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		logFatal(err)

		userBooks = append(userBooks, userBook)
	}


	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		log.Println(err)
		return
	}



}

func getUserBooks(w http.ResponseWriter, r *http.Request) {

	err := render.Template(w, "admin-user-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}


}

func getBooksApi(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, err := db.Query(`select * from books
							ORDER BY ID ASC`) // parami v db so $1
	logFatal(err)

	defer rows.Close()


	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Availability) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		logFatal(err)

		books = append(books, book)
	}

	if err = json.NewEncoder(w).Encode(books); err != nil {
		logFatal(err)
	}



}

func getBooks(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-Books.page.tmpl")

	if err != nil {
		log.Println(err)
	}


}




func getUsersApi(w http.ResponseWriter, r *http.Request) {
	var user User
	users = []User{}

	rows, err := db.Query(`select u.id, u.first_name , u.last_name, coalesce(b.count, 0) as borrowed_count 
						from users u left join (select user_id, count(user_id) 
						from borrowedbooks group by user_id) as b on b.user_id = u.id order by u.id asc`)
	logFatal(err)



	for rows.Next() {
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.BorrowedCount) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		logFatal(err)

		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(users); err != nil {
		logFatal(err)
	}

}


func getAllBorrowedBooksApi(w http.ResponseWriter, r *http.Request) {
	var userBook UserBooks
	userBooks = []UserBooks{}

	rows, err := db.Query(`select b.book_id as bookid, bo.title, u.* from borrowedbooks b 
						left join users u on u.id = b.user_id 
						left join books bo on bo.id = b.book_id
						ORDER BY b.book_id asc;
						`)

	if err != nil {
		log.Print(err)
	}
	logFatal(err)



	for rows.Next() {
		err := rows.Scan(&userBook.BookId, &userBook.Book, &userBook.ID, &userBook.FirstName, &userBook.LastName  ) // kam podatke vnesemo, v prazem book, vsakic spremenimo
		logFatal(err)

		userBooks = append(userBooks, userBook)
	}

	w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(userBooks); err != nil {
		logFatal(err)
	}

}

func getAllBorrowedBooks (w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-all-borrowed-books.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

func getUsers (w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-main.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}

func addUserPage (w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, "admin-add-user.page.tmpl")

	if err != nil {
		log.Println(err)
	}
}


func addUser(w http.ResponseWriter, r *http.Request) {
	var user User
	//var UserID int

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := r.PostForm
	if len(form) < 1{
		json.NewDecoder(r.Body).Decode(&user)
	} else {
		user.FirstName = form["first_name"][0]
		user.LastName = form["last_name"][0]
	}


	//


	row := db.QueryRow("select * from adduser($1,$2)", user.FirstName, user.LastName)

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


func deleteUser(w http.ResponseWriter, r *http.Request) {
	var user User
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
	err = db.QueryRow("DELETE from users where id =  $1 RETURNING id", // placeholderje za valueje k jih postamo,
		user.ID).Scan(&UserID)

	if err != nil {
		log.Println(err)
	}

	// vrnemo se v sql RETURNING id, ki samo pove kateri ID nam pripada, to si zapisemo skozi scan, in pol ga napisemo kot response POSTU,
	// tok da vemo kaj smo vstavil
	// pa idje itak sam vstavlja

	json.NewEncoder(w).Encode(UserID)
}


func borrowBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook BorrowBook
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

	row := db.QueryRow("select * from borrowbook($1,$2)", borrowBook.UserID, borrowBook.BookID)

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

func returnBook(w http.ResponseWriter, r *http.Request) {
	var borrowBook BorrowBook
	var Bookid int




	json.NewDecoder(r.Body).Decode(&borrowBook)
	err := db.QueryRow("DELETE FROM borrowedbooks WHERE user_id = $1 AND book_id = $2 RETURNING id", // placeholderje za valueje k jih postamo,
		borrowBook.UserID, borrowBook.BookID).Scan(&Bookid)

	if err != nil {
		_, err = w.Write([]byte("No such user/book relationship"))
		logFatal(err)
		return
	} else {
		json.NewEncoder(w).Encode("Success!")
	}




}