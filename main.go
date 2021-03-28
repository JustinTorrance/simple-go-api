package main

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm" 

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Person struct {
	gorm.Model

	Name string
	Books []Book
}

// Book Struct (Model)
type Book struct {
	gorm.Model
	ID 			string `json:"id"`
	Isbn 		string `json:"isbn"`
	Title 	string `json:"title"`
	Author 	*Author `json:"author"`
}

// Author Struct
type Author struct {
	gorm.Model
	Firstname 	string `json:"firstname"`
	Lastname 		string `json:"lastname"`
	Isbn				string `json:"isbn"`
}


// Init books var as a slice Book struct
// slice is a variable length array
var books []Book

func setupCorsResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
}

// Get All Books

// func getBooks(w http.ResponseWriter, router *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
//    if (*router).Method == "OPTIONS" {
//       return
//    }
// 	json.NewEncoder(w).Encode(books)
// }

func getBooks(w http.ResponseWriter, router *http.Request) {
	var books []Book
	db.Find(&favBooks)
	json.NewEncoder(w).Encode(&books)
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

// func createBook(w http.ResponseWriter, router *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var book Book
// 	_ = json.NewDecoder(router.Body).Decode(&book)
// 	book.ID = strconv.Itoa(rand.Intn(10000000)) //Mock ID - not safe
// 	books = append(books, book)
// 	json.NewEncoder(w).Encode(book)
// }

func createBook(w http.ResponseWriter, router *http.Request) {
	var book Book
	json.NewDecoder(router.Body).Decode(&book)

	createdPerson := db.Create(&book)
	err = createdPerson.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdBook)
}

// Update Book
func updateBook(w http.ResponseWriter, router *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(router)
	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book
			_ = json.NewDecoder(router.Body).Decode(&book)
			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}
}

// Delete Book

// func deleteBook(w http.ResponseWriter, router *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(router)
// 	for index, item := range books {
// 		if item.ID == params["id"] {
// 			books = append(books[:index], books[index+1:]...)
// 			break
// 		}
// 	}
// }

func deleteBook(w http.ResponseWriter, router *http.Request) {
	params := mux.Vars(router)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}

var (
	favBooks = &Book{ID: "269", Isbn: "0987", Title: "Norse Mythology"}
	authors = &Author{Firstname: "Neil", Lastname: "Gaiman", Isbn: "0987"}
)

var db *gorm.DB
var err error

func main() {
	//Init Router
	router := mux.NewRouter()

	//Loading environment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	//Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	//Opening connection to database
	db, err = gorm.Open(dialect, dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database!")
	}

	// Close connection to database when the main function finishes
	defer db.Close()

	//Make migrations to database if they have not already been created
	db.AutoMigrate(&Book{})
	db.AutoMigrate(&Author{})


	db.Create(favBooks)
	// db.Create(authors)


	// for idx := range books {
	// 	db.Create(&favBooks)
	// }

	//Route Handlers / Endpoints
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))

}