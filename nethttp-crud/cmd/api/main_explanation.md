# Main in Go

This file is the beating heart of your application.

## The Mental Model: "The Manager" (Manual Wiring)

**In Java (Spring Boot):**
Your `Application.java` uses `@SpringBootApplication` and `SpringApplication.run()`. Spring does "Component Scanning," magically wiring all dependencies together (IoC Container), mapping routes, and starting Tomcat. 

**In Go:**
The mental model is **"The Manager" (or The Master Weaver).**
Go has absolutely zero Component Scanning and zero magic Dependency Injection. The `main` function is fully responsible for manually creating the database, service, and handlers, wiring them together, defining the URLs, and turning on the server.

You can look at one single file (`main.go`) and instantly understand the entire architecture.

---

## Section-by-Section Syntax Breakdown

### 1. The Magic `main` Package
```go
package main

func main() {
```
*   **Why:** In Go, any file in a package named `main` tells the compiler: *"Compile me into an executable binary."* The `func main()` is the exact function the operating system calls to start the app.

### 2. Manual Dependency Injection (The Lego Bricks)
```go
bookStore := store.NewBookStore()
bookService := services.NewBookService(bookStore)
app := handlers.NewApp(bookService)
```
*   **Why:** This replaces Java's `@Autowired`.
*   **When/How:** We build the lowest layer first (Store), pass it into the Service, and pass the Service into the Handler. Because the constructors return pointers (`*`), we share the exact same memory addresses.

### 3. The Router (Multiplexer)
```go
mux := http.NewServeMux()
```
*   **What:** Creates a new HTTP Router.
*   **Why:** "Mux" stands for Multiplexer. Its job is to look at incoming HTTP requests (`GET /books`) and redirect traffic to the correct Handler function.

### 4. Explicit Route Registration
```go
mux.HandleFunc("GET /books", app.ListBooksHandler)
```
*   **Why:** This replaces `@RequestMapping("/books")`. 
*   **When/How (CRITICAL):** Notice that we write `app.ListBooksHandler` and **NOT** `app.ListBooksHandler()`. 
    *   If we used `()`, Go would execute the function immediately during startup!
    *   By omitting the parenthesis, we are passing the *memory address of the function*. We are telling the router: *"Here is the GPS location of the code. Save it and run it later when someone visits this URL."*

### 5. Starting the Server
```go
err := http.ListenAndServe(port, mux)
```
*   **Why:** Go has a production-grade web server built directly into the standard library.
*   **When/How:** `ListenAndServe` is a blocking call. It traps your program in an infinite loop listening for internet traffic. The code below this line will *never* execute unless the server crashes.

### 6. Fatal Error Logging
```go
if err != nil {
	log.Fatalf("Server failed to start: %v", err)
}
```
*   **Why:** If port 8080 is already being used, `ListenAndServe` immediately returns an error. `log.Fatalf` prints the error and immediately kills the program (like a fatal `Exception` bringing down the JVM).
