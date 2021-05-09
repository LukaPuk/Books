package main

import (
	"github.com/LukaPuk/Books/internal/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func InitRoutes() http.Handler {

	router := mux.NewRouter()
	router.HandleFunc("/", handlers.Homepage).Methods("GET")
	router.HandleFunc("/books/{id}", handlers.GetBook).Methods("GET")
	router.HandleFunc("/booksapi", handlers.GetBooksApi).Methods("GET")
	router.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	router.HandleFunc("/usersapi", handlers.GetUsersApi).Methods("GET")
	router.HandleFunc("/users", handlers.GetUsers).Methods("GET")
	router.HandleFunc("/add-user", handlers.AddUser).Methods("POST")
	router.HandleFunc("/add-user", handlers.AddUserPage).Methods("GET")
	router.HandleFunc("/borrowbook", handlers.BorrowBook).Methods("POST")
	router.HandleFunc("/returnbook", handlers.ReturnBook).Methods("POST")
	router.HandleFunc("/delete", handlers.DeleteUser).Methods("POST")
	router.HandleFunc("/userapi/{id}", handlers.GetUserBooksApi).Methods("GET")
	router.HandleFunc("/user/{id}", handlers.GetUserBooks).Methods("GET")
	router.HandleFunc("/all-borrowed-booksapi", handlers.GetAllBorrowedBooksApi).Methods("GET")
	router.HandleFunc("/all-borrowed-books", handlers.GetAllBorrowedBooks).Methods("GET")

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../static/"))))

	return router
}
