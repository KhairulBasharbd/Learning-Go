# Store in Go

This file represents your Database Layer (often called the Data Access Layer or Repository Layer).

## The Mental Model: "Explicit over Magic"

**In Java (Spring Boot):**
You create an interface like `interface BookRepository extends JpaRepository`. Spring Boot uses reflection and magic to instantly generate SQL queries, handle database connection pools, and map database rows to your `@Entity` objects.

**In Go:**
Go hates magic. The mental model for a Go Store/Repository is **"Explicit Data Access."**
The job of this file is to act as a shield. The rest of your application should have no idea if you are using PostgreSQL, MongoDB, or an in-memory Map. Inside this file, you must explicitly manage the state, write SQL queries, and manually map rows to Go structs.

---

## Section-by-Section Syntax Breakdown

### 1. The State Container (The "Database")
```go
type BookStore struct {
	mu    sync.RWMutex
	books map[string]models.Book
}
```
*   **What:** We define a struct that holds the state of our database.
*   **When/How:** Notice `BookStore` is capitalized (Public), so `main.go` can see it. But `mu` and `books` are lowercase (Private). This prevents the Handler or Service from accidentally bypassing our functions and mutating the map directly!

### 2. The Concurrency Protector (Mutex)
```go
	mu    sync.RWMutex
```
*   **Why (CRITICAL CONCEPT):** Every HTTP request in Go runs in its own concurrent thread (Goroutine). **Go maps are not thread-safe.** If two threads write to a map at once, the Go program will instantly crash. We use `mu` to force Goroutines to wait in line.

### 3. The Constructor Pattern
```go
func NewBookStore() *BookStore { ... }
```
*   **Why:** Go does not have `new` constructors. The community standard is to write a regular function starting with `New...` to initialize your structs. 
*   **The Pointer:** It returns a pointer (`*BookStore`). By returning a pointer, `main.go`, `services`, and `handlers` all share the exact same physical memory location.

### 4. The Method Receiver
```go
func (s *BookStore) Create(book models.Book) error {
```
*   **Why:** Go doesn't have classes. To make `Create` a method belonging to `BookStore`, we use `(s *BookStore)` before the function name. This is called a **Receiver**. 

### 5. Locking and Deferring
```go
	s.mu.Lock() 
	defer s.mu.Unlock()
```
*   **Why:** `s.mu.Lock()` locks the door to the database. `defer s.mu.Unlock()` tells Go: *"Schedule this unlock operation to run at the absolute last microsecond before the function returns."* This guarantees the door unlocks even if the function crashes.

### 6. Read Locks vs Full Locks
```go
	s.mu.RLock()
```
*   **Why:** `RLock()` is a "Read Lock". It allows infinite simultaneous readers, but if someone tries to write, it forces the writer to wait. This makes the web server insanely fast.
