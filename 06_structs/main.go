package main

import "fmt"

type Account struct {
	id      int
	balance int
}

func (a *Account) Deposit(amount int) {
	a.balance += amount
}
func (a *Account) Withdraw(amount int) error {
	if a.balance < amount {
		return fmt.Errorf("insufficient balance")
	}
	a.balance -= amount
	return nil
}
func (a Account) Balance() int {
	return a.balance
}

func NewAccount(id, balance int) *Account {
	return &Account{
		id,
		balance,
	}
}

func main() {
	a := NewAccount(1, 100)
	a.Withdraw(10)
	fmt.Println(a.Balance())
	a.Deposit(10)
	fmt.Println(a.Balance())
	err := a.Withdraw(200)
	if err != nil {
		fmt.Println(err)
	}
}
