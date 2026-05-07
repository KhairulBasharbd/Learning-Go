# Services in Go

This file represents your Service Layer (or Business Logic Layer). 

## The Mental Model: "The Pure Middleman"

**In Java (Spring Boot):**
You use the `@Service` annotation. It sits between your Controller and your Repository. You use `@Autowired` to magically inject the Repository, and `@Transactional` to magically wrap the method in a database transaction.

**In Go:**
The mental model is exactly the same: **The Service is the Head Chef.** It enforces the business rules. 

However, in Go, the Service layer must remain **Pure**:
1. **No HTTP:** It should never import `net/http`. It doesn't know what JSON or URLs are.
2. **No Databases:** It should not write SQL queries.
3. **No Magic:** There is no `@Autowired`. We explicitly hand the Store to the Service.

---

## Section-by-Section Syntax Breakdown

### 1. Dependency Injection (The Go Way)
```go
type BookService struct {
	store *store.BookStore
}
```
*   **What:** We define the `BookService` struct and give it a field that holds a pointer (`*`) to the Database layer.
*   **Why:** This is how Go does Dependency Injection. Instead of an `@Autowired` annotation, the Service struct explicitly demands a store to do its job.

### 2. The Constructor
```go
func NewBookService(s *store.BookStore) *BookService {
```
*   **What:** The standard Go pattern for creating a new instance.
*   **Why:** Because Go doesn't have an IoC Container, `main.go` creates the Store, passes it into this function, and this function locks the Store inside the new `BookService`.

### 3. Passthrough Methods
```go
func (s *BookService) GetAllBooks() []models.Book {
	return s.store.GetAll()
}
```
*   **Why:** Sometimes, getting data requires no business logic. But we *must* write this passthrough method anyway because the HTTP Handler is only allowed to talk to the Service, not the Store directly.

### 4. Business Logic and Error Handling
```go
func (s *BookService) CreateBook(book models.Book) error {
	if book.ID == "" || book.Title == "" {
		return errors.New("ID and Title are required")
	}
	return s.store.Create(book)
}
```
*   **Why:** In Java, you might use `@NotNull` on the Entity, or `throw new IllegalArgumentException()`. In Go, we explicitly check the data and return an `error` value.
*   **When/How:** If the validation fails, the function immediately returns the error back to the HTTP Handler. The Store is never touched. If validation passes, it calls the Store.
