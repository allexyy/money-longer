package domain

import "time"

type Transaction struct {
	ID       int
	Name     string
	VaultId  int `json:"budget_id"`
	Date     string
	Amount   int
	IsIncome bool
	Note     string
}

type TransactionRepository interface {
	GetForPeriod(period time.Time) ([]Transaction, error)
	GetAll() ([]Transaction, error)
	GetByID(id int) (*Transaction, error)
	Create(t Transaction) error
	Update(t Transaction) error
	Delete(id int) error
}

type TransactionCreatedEvent struct {
	VaultID  int
	Amount   int
	IsIncome bool
}
