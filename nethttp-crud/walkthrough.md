# The Go Web Flow: Spring Boot to Go Translation

Here is exactly how our Go project maps to Java Spring Boot:

## 1. The Direct Spring Boot Mapping

*   **`models/book.go` = `@Entity` or DTO Class**
    *   *Spring Boot:* `public class Book { private String id; ... }`
    *   *Go:* Just a struct defining data. No logic.

*   **`store/memory.go` = `@Repository`**
    *   *Spring Boot:* `public interface BookRepository extends JpaRepository<Book, String>`
    *   *Go:* The database layer. It holds the data and contains functions like `Create()`, `Get()`. It knows *nothing* about HTTP or business rules.

*   **`services/book.go` = `@Service`**
    *   *Spring Boot:* `@Service public class BookService { ... }`
    *   *Go:* The business logic layer. It enforces rules. It talks to the Store but knows *nothing* about HTTP requests.

*   **`handlers/book.go` = `@RestController`**
    *   *Spring Boot:* `@RestController @RequestMapping("/books") public class BookController { ... }`
    *   *Go:* The HTTP controllers. They take the incoming JSON, hand it to the Service layer to process, and send JSON back. 

*   **`cmd/api/main.go` = `@SpringBootApplication`**
    *   *Spring Boot:* `public static void main(String[] args) { SpringApplication.run(App.class, args); }`
    *   *Go:* The entrypoint. It wires all 3 layers together manually: `Store -> Service -> Handler`.

---

## 2. The File Relationships (The Restaurant Story)

1. **`models/book.go` (The Recipe Card)**: Defines *what* a book is.
2. **`store/memory.go` (The Fridge)**: The physical place where food is kept.
3. **`services/book.go` (The Head Chef)**: Enforces business rules. He checks if the recipe is valid. If it's valid, he puts it in the Fridge.
4. **`handlers/book.go` (The Waiters)**: Waiters talk to customers. They take the customer's order, hand it to the Head Chef (`service`), and wait for the food to give back to the customer. 
5. **`cmd/api/main.go` (The Manager)**: Hires the Fridge, hires the Chef, gives the Chef access to the Fridge, hires the Waiter, gives the Waiter access to the Chef, and opens the doors.

---

## 3. Step-by-Step Flow of a Web Request (Create a Book)
Let's trace exactly what happens when a user sends: `POST http://localhost:8080/books`

### Step 1: The Front Door (`main.go`)
```go
mux.HandleFunc("POST /books", app.CreateBookHandler)
```
The router sees a `POST` request and sends it to the Waiter.

### Step 2: The Waiter Takes the Order (`handlers/book.go`)
```go
func (a *App) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
    var book models.Book
    json.NewDecoder(r.Body).Decode(&book) 
    
    err := a.service.CreateBook(book)
```
The Waiter reads the JSON payload, converts it to a Go struct, and immediately hands it off to the Service layer.

### Step 3: The Chef Enforces Rules (`services/book.go`)
```go
func (s *BookService) CreateBook(book models.Book) error {
	if book.ID == "" || book.Title == "" {
		return errors.New("ID and Title are required") 
	}
	return s.store.Create(book) 
}
```
The Service layer checks the business rules. If it fails, it returns an error back to the Waiter. If it passes, it asks the Store to save it.

### Step 4: The Fridge Saves It (`store/memory.go`)
```go
func (s *BookStore) Create(book models.Book) error {
	s.mu.Lock()
	defer s.mu.Unlock()
}
```
The store locks the map, saves the data, and returns `nil` (no error).

### Step 5: The Waiter Serves the Result (`handlers/book.go`)
Back in the Waiter's function, if there was an error from the Chef, he tells the customer (`sendError(400)`). If there was no error, he sends the food (`sendJSON(201 Created)`).
