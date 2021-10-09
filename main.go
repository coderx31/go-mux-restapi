package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Book struct

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct

type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type CustomError struct {
	Error      error `json:"error"`
	StatusCode int16 `json:"statusCode"`
}

var books []Book

// get all books
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// get single book
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params
	log.Println(params)
	for _, item := range books {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return // if found
		}
	}

	// if not found
	json.NewEncoder(w).Encode(&Book{})

}

// create a new book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book

	err := json.NewDecoder(r.Body).Decode(&book)

	if err != nil {
		json.NewEncoder(w).Encode(&CustomError{Error: err, StatusCode: 500})
		return
	}

	book.ID = strconv.Itoa(rand.Intn(10000000))
	books = append(books, book)

	json.NewEncoder(w).Encode(&book)

}

// update book
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			var book Book

			err := json.NewDecoder(r.Body).Decode(&book)

			if err != nil {
				json.NewEncoder(w).Encode(&CustomError{Error: err, StatusCode: 500})
				return
			}

			book.ID = params["id"]
			books = append(books, book)
			json.NewEncoder(w).Encode(book)
			return
		}
	}

}

// delete book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r) // get params

	for index, item := range books {
		if item.ID == params["id"] {
			books = append(books[:index], books[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func main() {
	// init the router
	r := mux.NewRouter()

	// mock data
	books = append(books, Book{ID: "1", Isbn: "58654645", Title: "Book one",
		Author: &Author{Firstname: "John", Lastname: "Doe"}})

	books = append(books, Book{ID: "2", Isbn: "7897854", Title: "Book two",
		Author: &Author{Firstname: "Steve", Lastname: "Smith"}})

	// route handlers - endpoints
	r.HandleFunc("/api/books", getBooks).Methods("GET")
	r.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	r.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// to start server
	log.Fatal(http.ListenAndServe(":8000", r))
}
