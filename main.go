package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);unique_index"`
	Books []Book
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int
	UserID   int
}

var db *gorm.DB
var err error

func main() {
	// Loading enviroment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbname := os.Getenv("NAME")
	dbpassword := os.Getenv("PASSWORD")

	// Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbname, dbpassword, dbPort)

	// Openning connection to database
	db, err = gorm.Open(dialect, dbURI)

	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connected to database successfully")
	}

	// Close the databse connection when the main function closes
	defer db.Close()

	// Make migrations to the database if they haven't been made already
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Book{})

	/*----------- API routes ------------*/
	router := mux.NewRouter()

	// router.HandleFunc("/books", GetBooks).Methods("GET")
	// router.HandleFunc("/book/{id}", GetBook).Methods("GET")
	// router.HandleFunc("/users", GetPeople).Methods("GET")
	// router.HandleFunc("/user/{id}", GetUser).Methods("GET")

	// router.HandleFunc("/create/user", CreateUser).Methods("POST")
	// router.HandleFunc("/create/book", CreateBook).Methods("POST")

	// router.HandleFunc("/delete/user/{id}", DeleteUser).Methods("DELETE")
	// router.HandleFunc("/delete/book/{id}", DeleteBook).Methods("DELETE")

	router.HandleFunc("/api/books", GetBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", GetBook).Methods("GET")
	router.HandleFunc("/api/books", CreateBook).Methods("POST")
	// router.HandleFunc("/api/books/{id}", UpdateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

/*-------- API Controllers --------*/

/*----- People ------*/
func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	var books []Book

	db.First(&user, params["id"])
	db.Model(&user).Related(&books)

	user.Books = books

	json.NewEncoder(w).Encode(&user)
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	var users []User

	db.Find(&users)

	json.NewEncoder(w).Encode(&users)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	createdUser := db.Create(&user)
	err = createdUser.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdUser)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var user User

	db.First(&user, params["id"])
	db.Delete(&user)

	json.NewEncoder(w).Encode(&user)
}

/*------- Books ------*/

func GetBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var book Book

	db.First(&book, params["id"])

	json.NewEncoder(w).Encode(&book)
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	var books []Book

	db.Find(&books)

	json.NewEncoder(w).Encode(&books)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	createdBook := db.Create(&book)
	err = createdBook.Error
	if err != nil {
		fmt.Println(err)
	}

	json.NewEncoder(w).Encode(&createdBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book

	db.First(&book, params["id"])
	db.Delete(&book)

	json.NewEncoder(w).Encode(&book)
}