package main

import (
	"encoding/json"
	"log"
	"net/http"
	"math/rand"
	"strconv"
	"github.com/gorilla/mux"
)

// Book Struct (Model)
type Book struct {
	ID 			string `json:"id"`
	Isbn 		string `json:"isbn"`
	Title 	string `json:"title"`
	Author 	*Author `json:"author"`
}

// Author Struct
type Author struct {
	Firstname 	string `json:"firstname"`
	Lastname 		string `json:"lastname"`
}

// Init books var as a slice Book struct
// slice is a variable length array
var books []Book

// Get All Books
func getBooks(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Get Single Book
func getBook(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router) // Get params
	//Loop through books and find id
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

// Create New Book
func createBook(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(router.body).Decode(&book)
	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - not safe
}

// Update Book
func updateBook(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
}

// Delete Book
func deleteBook(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
}

func main() {
	//Init Router
	router := mux.NewRouter()

	// Mock Data @todo - implement DB
	books = append(books, Book{ID: "1", Isbn: "33857", Title: "Gone With The Wind", Author: &Author {Firstname: "Bob", Lastname: "Johnson"}})
	books = append(books, Book{ID: "2", Isbn: "89867", Title: "Call of the Wild", Author: &Author {Firstname: "Jack", Lastname: "London"}})

	//Route Handlers / Endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}