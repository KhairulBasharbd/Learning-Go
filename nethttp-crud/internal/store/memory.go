package store

import (
	"errors"
	"sync"

	"nethttp-crud/internal/models"
)

// BookStore is a thread-safe in-memory "database".
type BookStore struct {
	mu    sync.RWMutex
	books map[string]models.Book
}

func NewBookStore() *BookStore {
	return &BookStore{
		books: make(map[string]models.Book),
	}
}

func (s *BookStore) GetAll() []models.Book {
	s.mu.RLock() // Read lock
	defer s.mu.RUnlock()

	var result []models.Book
	for _, b := range s.books {
		result = append(result, b)
	}
	return result
}

func (s *BookStore) Get(id string) (models.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	book, exists := s.books[id]
	return book, exists
}

func (s *BookStore) Create(book models.Book) error {
	s.mu.Lock() // Write lock
	defer s.mu.Unlock()

	if _, exists := s.books[book.ID]; exists {
		return errors.New("book already exists")
	}
	s.books[book.ID] = book
	return nil
}

func (s *BookStore) Update(book models.Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[book.ID]; !exists {
		return errors.New("book not found")
	}
	s.books[book.ID] = book
	return nil
}

func (s *BookStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.books[id]; !exists {
		return errors.New("book not found")
	}
	delete(s.books, id)
	return nil
}
