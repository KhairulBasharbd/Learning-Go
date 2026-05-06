package handlers

import (
	"encoding/json"
	"net/http"

	"nethttp-crud/internal/models"
	"nethttp-crud/internal/store"
)

// App acts as our "Controller". It holds all dependencies our handlers need.
// Notice it is capitalized (Exported) so main.go can use it.
type App struct {
	store *store.BookStore
}

// NewApp is a constructor function to create our App struct.
func NewApp(s *store.BookStore) *App {
	return &App{store: s}
}

// ------------------------------------------------------------------
// HELPER FUNCTIONS 
// Notice these are lowercase. In Go, this means they are "Private" 
// and can only be used inside the `handlers` package.
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
// Notice these are now Capitalized. They must be exported (public) 
// so that `main.go` can access them to register them in the router.
// ------------------------------------------------------------------

func (a *App) ListBooksHandler(w http.ResponseWriter, r *http.Request) {
	books := a.store.GetAll()
	if books == nil {
		books = []models.Book{} 
	}
	sendJSON(w, http.StatusOK, books)
}

func (a *App) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	book, exists := a.store.Get(id)
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

	if book.ID == "" || book.Title == "" {
		sendError(w, http.StatusBadRequest, "ID and Title are required")
		return
	}

	if err := a.store.Create(book); err != nil {
		sendError(w, http.StatusConflict, err.Error())
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

	if err := a.store.Update(book); err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	sendJSON(w, http.StatusOK, book)
}

func (a *App) DeleteBookHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	
	if err := a.store.Delete(id); err != nil {
		sendError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent) 
}
