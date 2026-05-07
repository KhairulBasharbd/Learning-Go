package handlers

import (
	"encoding/json"
	"net/http"

	"nethttp-crud/internal/models"
	"nethttp-crud/internal/services"
)

// App acts as our "Controller". 
// Notice it now depends on the SERVICE, not the STORE directly!
type App struct {
	service *services.BookService
}

func NewApp(s *services.BookService) *App {
	return &App{service: s}
}

// ------------------------------------------------------------------
// HELPER FUNCTIONS 
// ------------------------------------------------------------------

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{"error": message})
}

// ------------------------------------------------------------------
// HTTP HANDLERS 
// ------------------------------------------------------------------

func (a *App) ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	// The handler asks the service for data, not the database.
	books := a.service.GetAllBooks()
	if books == nil {
		books = []models.Book{} 
	}
	sendJSON(w, http.StatusOK, books)
}

func (a *App) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	book, exists := a.service.GetBook(id)
	if !exists {
		sendError(w, http.StatusNotFound, "Book not found")
		return
	}
	sendJSON(w, http.StatusOK, book)
}

func (a *App) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	// Notice we removed the business logic (validation) from here.
	// We just pass the struct to the service layer.
	if err := a.service.CreateBook(book); err != nil {
		// If the service returns an error (like "ID required"), we send it to the user.
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusCreated, book)
}

func (a *App) UpdateBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		sendError(w, http.StatusBadRequest, "Invalid JSON format")
		return
	}

	book.ID = id 

	if err := a.service.UpdateBook(book); err != nil {
		sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, book)
}

func (a *App) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	if err := a.service.DeleteBook(id); err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}
