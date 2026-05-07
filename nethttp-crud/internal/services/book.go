package services

import (
	"errors"
	"nethttp-crud/internal/models"
	"nethttp-crud/internal/store"
)

// BookService acts as the Business Logic layer (equivalent to @Service).
// It talks to the Store (@Repository) but knows nothing about HTTP or JSON.
// This section is like java @Service Dependency injection
type BookService struct {
	store *store.BookStore
}

// NewBookService creates a new service instance, injecting the store dependency.
// This section is like java initialization / constructor / @Autowired annotation
func NewBookService(s *store.BookStore) *BookService {
	return &BookService{store: s}
}

func (s *BookService) GetAllBooks() []models.Book {
	return s.store.GetAll()
}

func (s *BookService) GetBook(id string) (models.Book, bool) {
	return s.store.Get(id)
}

func (s *BookService) CreateBook(book models.Book) error {
	// BUSINESS LOGIC: We moved the validation out of the HTTP handler!
	// Now the service is responsible for business rules.
	if book.ID == "" || book.Title == "" {
		return errors.New("ID and Title are required")
	}

	// We could add more logic here: "Does this user have permission?" Like check in DB or any file validation or
	// "Send an email that a book was created", etc.

	return s.store.Create(book)
}

func (s *BookService) UpdateBook(book models.Book) error {
	if book.Title == "" {
		return errors.New("Title cannot be empty")
	}
	return s.store.Update(book)
}

func (s *BookService) DeleteBook(id string) error {
	return s.store.Delete(id)
}
