# Models in Go

## The Mental Model: "Dumb Data" vs "Smart Objects"

**In Java (Spring Boot):**
Your `@Entity` or DTO classes are often "heavy." They have private fields, getters, setters, constructors, and `@OneToMany` annotations that magically link database tables together. They are "smart" objects.

**In Go:**
The mental model for a Go struct is **"Dumb Data."** 
A Go struct should just be a plain container holding data in memory. It doesn't know how to save itself to a database, it doesn't have getters and setters, and it doesn't have hidden magical behaviors. It is explicitly designed to be simple, predictable, and transparent.

---

## Line-by-Line Syntax Breakdown

### 1. Package Declaration
```go
package models
```
*   **What:** The package declaration.
*   **Why:** It tells the Go compiler that this file belongs to the `models` namespace. 
*   **When/How:** It must be the very first line of code in every Go file. 

### 2. Struct Definition & Capitalization
```go
type Book struct {
```
*   **What:** We are defining a new custom shape in memory. 
*   **Capitalization Rule (CRITICAL):** Notice the `B` in `Book` is **Capitalized**. In Go, if a variable, function, or struct starts with a capital letter, it is **Exported (Public)**. If it was `type book struct`, the Handlers and Services would not be able to see it!

### 3. Fields & Visibility
```go
	ID     string
	Title  string
```
*   **What:** The fields of the struct and their data types.
*   **Capitalization Rule:** `ID` and `Title` are capitalized so they are **Public**. If they were lowercase (`id string`), they would be **Private**, and the JSON encoder would not be allowed to read them. 

### 4. Struct Tags
```go
	`json:"id"`
```
*   **What:** This is called a **Struct Tag**. It is a raw string literal surrounded by backticks.
*   **Why:** This is the Go equivalent of `@JsonProperty("id")`. It tells the Go JSON encoder to turn the capitalized `ID` field into a lowercase JSON key like `{"id": "123"}`. 
