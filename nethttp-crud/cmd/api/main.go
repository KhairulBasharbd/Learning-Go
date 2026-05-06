package main

import (
	"fmt"
	"log"
	"net/http"

	"nethttp-crud/internal/handlers"
	"nethttp-crud/internal/store"
)

func main() {
	// 1. Initialize our Data Store
	bookStore := store.NewBookStore()

	// 2. Initialize our App struct that holds the dependencies
	app := handlers.NewApp(bookStore)

	// 3. Create a new HTTP Multiplexer (Router)
	mux := http.NewServeMux()

	// 4. Register Routes
	mux.HandleFunc("GET /books", app.ListBooksHandler)
	mux.HandleFunc("GET /books/{id}", app.GetBookHandler)
	mux.HandleFunc("POST /books", app.CreateBookHandler)
	mux.HandleFunc("PUT /books/{id}", app.UpdateBookHandler)
	mux.HandleFunc("DELETE /books/{id}", app.DeleteBookHandler)

	// 5. Start the Server
	port := ":8080"
	fmt.Printf("Starting pure net/http server on http://localhost%s\n", port)
	
	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
