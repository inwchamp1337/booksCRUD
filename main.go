package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/inwchamp1337/CRYD/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Book struct {
    ID     int     `json:"id"`
    Isbn   string  `json:"isbn"`
    Title  string  `json:"title"`
    Author *Author `json:"author"`
}

type Author struct {
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

var books []Book
// getBooks godoc
// @Summary      Get all books
// @Description  get all books
// @Tags         books
// @Produce      json
// @Success      200  {array}   Book
// @Router       /books [get]
func getBooks(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(books)
}
// getBookByID godoc
// @Summary      Get book by ID
// @Description  Get a single book by its ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Book ID"
// @Success      200  {object}  Book
// @Failure      404  {string}  string  "Book not found"
// @Router       /books/{id} [get]
func getBook(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    log.Printf("GET /books/%s called", params["id"])
    for _, item := range books {
        if strconv.Itoa(item.ID) == params["id"] {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    http.Error(w, "Book not found", http.StatusNotFound)
}
// CreateBook godoc
// @Summary      Create a new book
// @Description  add by json book
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      Book  true  "Add book"
// @Success      201   {object}  Book
// @Router       /books [post]
func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	book.ID = len(books) + 1
	books = append(books, book)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}
// UpdateBook godoc
// @Summary      Update a book
// @Description  Update a book by ID
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Book ID"
// @Param        book  body      Book    true  "Update book"
// @Success      200   {object}  Book
// @Failure      404   {string}  string  "Book not found"
// @Router       /books/{id} [put]
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if strconv.Itoa(item.ID) == params["id"] {
			var book Book
			err := json.NewDecoder(r.Body).Decode(&book)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			book.ID = item.ID // Keep the same ID
			books[index] = book
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}
// DeleteBook godoc
// @Summary      Delete a book
// @Description  Delete a book by ID
// @Tags         books
// @Param        id   path      string  true  "Book ID"
// @Success      204  {string}  string  "No Content"
// @Failure      404  {string}  string  "Book not found"
// @Router       /books/{id} [delete]
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range books {
		if strconv.Itoa(item.ID) == params["id"] {
			books = append(books[:index], books[index+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}
// @title           Book API
// @version         1.0
// @description     This is a sample server for managing books.
// @host            localhost:8000
// @BasePath        /
func main() {
    r := mux.NewRouter()

    books = append(books, Book{ID: 1, Isbn: "12345", Title: "Book One", Author: &Author{FirstName: "John", LastName: "Doe"}})
    books = append(books, Book{ID: 2, Isbn: "67890", Title: "Book Two", Author: &Author{FirstName: "Jane", LastName: "Smith"}})

    r.HandleFunc("/books", getBooks).Methods("GET")
    r.HandleFunc("/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/books", CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

    log.Println("Server is running on port 8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}