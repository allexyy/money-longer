package domain

import "time"

type Vault struct {
	ID         int
	Name       string
	Limit      int
	LeftAmount int
	Icon       string
	Color      string
	Expire     time.Time
}

type VaultRepository interface {
	GetAll() ([]Vault, error)
	GetById(id int) (*Vault, error)
	Create(v Vault) error
	Update(v Vault) error
	Delete(id int) error
}
