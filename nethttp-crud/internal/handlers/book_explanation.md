# Handlers in Go

This file represents your Controller Layer (the HTTP interface).

## The Mental Model: "The Waiter (HTTP Translator)"

**In Java (Spring Boot):**
You use `@RestController`. Spring does a lot of invisible work for you: if you put `@RequestBody Book book`, Spring automatically reads the HTTP stream, parses the JSON, and hands it to you. If you return a `Book`, Spring automatically converts it back to JSON.

**In Go:**
The mental model is **"The Waiter (Translator)."** 
The Handler's *only* job is to speak HTTP. It translates "Internet Bytes" into "Go Structs," hands those structs to the Service (The Chef), and then translates the Chef's response back into "Internet Bytes" (JSON).

**You** must tell Go to read the JSON, and **you** must tell Go exactly which HTTP status code to return.

---

## Section-by-Section Syntax Breakdown

### 1. Controller State & Dependency Injection
```go
type App struct {
	service *services.BookService
}
```
*   **Why:** This is the Go equivalent of a `@RestController` class with an `@Autowired` Service.

### 2. Private Helper Functions
```go
func sendJSON(w http.ResponseWriter, status int, data any) {
```
*   **Why:** Because Go doesn't magically assume you want to send JSON, you have to write the serialization code yourself. To avoid copying and pasting, we create a helper.
*   **When/How:** Notice it starts with a **lowercase `s`**. This means it is **Private**. Only the handlers inside this specific package can use it. It takes `data any` (like Java's `Object`).

### 3. The Handler Signature
```go
func (a *App) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
```
*   **What:** The strictly enforced signature for *any* HTTP handler.
*   **Why:** 
    *   `w http.ResponseWriter`: Your "Serving Plate." You use it to set headers and stream data out.
    *   `r *http.Request`: The "Customer's Order." It contains the URL, headers, and JSON body. It is a pointer (`*`) so we don't copy massive network requests in memory.

### 4. Extracting Path Variables
```go
id := r.PathValue("id")
```
*   **Why:** The Go equivalent of `@PathVariable("id")`. 

### 5. Decoding JSON (Parsing `@RequestBody`)
```go
var book models.Book
if err := json.NewDecoder(r.Body).Decode(&book); err != nil { ... }
```
*   **Why:** The manual equivalent of `@RequestBody`. 
*   **When/How:** We tell the decoder to read the live `r.Body` stream and pour the data directly into the memory address of our variable (`&book`).

### 6. Communicating with the Service
```go
if err := a.service.CreateBook(book); err != nil { ... }
```
*   **Why:** Notice there is **zero business logic here**. We just pass the struct to the Service. If the Service rejects it (returns an error), we translate that error into an HTTP 400 status.
