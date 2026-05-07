package main

import (
	"fmt"
	"log"
	"net/http"

	"nethttp-crud/internal/handlers"
	"nethttp-crud/internal/services"
	"nethttp-crud/internal/store"
)

func main() {
	// 1. Initialize the Database Layer (@Repository)
	bookStore := store.NewBookStore()

	// 2. Initialize the Business Logic Layer (@Service)
	// We inject the store into the service.
	bookService := services.NewBookService(bookStore)

	// 3. Initialize the HTTP Layer (@RestController)
	// We inject the service into the handler.
	app := handlers.NewApp(bookService)

	// 4. Create a new HTTP Multiplexer (Router)
	mux := http.NewServeMux()

	// 5. Register Routes
	mux.HandleFunc("GET /books", app.ListBooksHandler)
	mux.HandleFunc("GET /books/{id}", app.GetBookHandler)
	mux.HandleFunc("POST /books", app.CreateBookHandler)
	mux.HandleFunc("PUT /books/{id}", app.UpdateBookHandler)
	mux.HandleFunc("DELETE /books/{id}", app.DeleteBookHandler)

	// 6. Start the Server
	port := ":8080"
	fmt.Printf("Starting pure net/http server on http://localhost%s\n", port)

	err := http.ListenAndServe(port, mux)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
