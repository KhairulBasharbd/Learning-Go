package main

func main() {

	type Book struct {
		Id     string  `json : "id"`
		Title  string  `json : "title"`
		Author string  `json : "author"`
		Price  float64 `json : "price"`
	}

}
