Books Program

Zagon:

- cmd/web -> go build -> .\web

http://localhost:8080/ Za Homepage

Baza:

- spisana v Postgres z Soda migrations (Vse tabele so v migrations folderju)
- Tables:
  - Books (ID, Title, Availability)
  - Users (ID, First Name, Last Name)
  - Borrowed Books (ID, User ID, Books ID)

- Dodani Triggerji Update Books on Delete/Insert v Borrowed_books (ce si sposodis knjigo, availability zmanjsa za 1 ter obratno)
- Dodani foreign keys Users -> borrowed books, Books -> Borrowed books (ce se user zbrise se vrnejo sposojene knjige)
- Dodan unique za Borrowed books, da si en user ne more veckrat iste knjige sposoditi
- Dodani funkciji Book_availability_check (lahko sposodis ce je availability > 0) ter Add_user_check (pogleda ce je user ze v bazi)

Handlers:

```go
router := mux.NewRouter()
	"/", ("GET") - Homepage
	"/books/{id}", ("GET") -JSON z posamezno Book ID, Title, Availability
	"/booksapi", ("GET") -JSON z vsemi Book ID, Title, Availability
	"/books", ("GET") - Spletna stran z vsemi Books
	"/usersapi", ("GET") - JSON z vsemi userji(Id,first name, last name, borrow count)
	"/users", ("GET") - Spletna stran z Userji
	"/userbooksapi/{id}", ("GET") - JSON z knjigami od Userja 
	"/userbooks/{id}", ("GET") - Stran z knjigami od Userja
	"/add-user", ("POST") - Add user (first_name, last_name) - Lahko FORM ali JSON
	"/add-user", ("GET") - Spletna stran za dodajanje userjev
	"/borrowbook", ("POST") - Borrow book (user_id, book_id)
	"/returnbook", ("POST") - Return Book Post (user_id, book_id)
	"/delete", ("POST") - Delete User (user_id)
	"/all-borrowed-booksapi", ("GET") - JSON z vsemi sposojenemi knjigami
	"/all-borrowed-books", ("GET") - Stran z vsemi sposojenemi knjigami

```



