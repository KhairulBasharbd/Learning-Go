package models

// Book represents our data model. The backtick tags tell the json
// encoder how to format the keys when converting to/from JSON.
type Book struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
