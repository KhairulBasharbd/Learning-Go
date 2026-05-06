package main

import "fmt"

// 1. The Implementation (Notice: No idea that any interface exists!)
type CreditCard struct {
	CardNumber string
}

func (c CreditCard) Charge(amount float64) bool {
	fmt.Printf("Charging $%.2f to card %s\n", amount, c.CardNumber)
	return true
}

type PayPal struct {
	Email string
}

func (p PayPal) Charge(amount float64) bool {
	fmt.Printf("Charging $%.2f to PayPal account %s\n", amount, p.Email)
	return true
}

// 2. The Consumer's Interface
// We define what we need right next to where we use it.
type PaymentMethod interface {
	Charge(amount float64) bool
}

// 3. The Polymorphic Function
// It accepts anything that satisfies the PaymentMethod interface.
func ProcessCheckout(p PaymentMethod, total float64) {
	success := p.Charge(total)
	if success {
		fmt.Println("Checkout complete!")
	}
}

func main() {
	myCard := CreditCard{CardNumber: "1234-5678-9012"}
	myPayPal := PayPal{Email: "user@example.com"}

	// Go automatically checks the methods.
	// Both structs have a Charge(float64) bool method, so both are valid!
	ProcessCheckout(myCard, 50.00)
	ProcessCheckout(myPayPal, 25.50)
}
