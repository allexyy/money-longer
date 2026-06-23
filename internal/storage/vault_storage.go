package storage

import (
	"database/sql"
	"monyLonger/internal/domain"
)

type VaultsStorage struct {
	db *sql.DB
}

func NewVaultsStorage(db *sql.DB) *VaultsStorage {
	return &VaultsStorage{db: db}
}

func (s *VaultsStorage) GetAll() ([]domain.Vault, error) {
	rows, err := s.db.Query("SELECT * FROM vaults")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Vault
	for rows.Next() {
		var v domain.Vault
		if err := rows.Scan(&v.ID, &v.Name, &v.Limit, &v.LeftAmount, &v.Icon, &v.Color, &v.Expire); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (s *VaultsStorage) GetById(id int) (*domain.Vault, error) {
	rows, err := s.db.Query("SELECT * FROM vaults where id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result domain.Vault
	for rows.Next() {
		var v domain.Vault
		if err := rows.Scan(&v.ID, &v.Name, &v.Limit, &v.LeftAmount, &v.Icon, &v.Color, &v.Expire); err != nil {
			return nil, err
		}
		result = v
	}
	return &result, nil
}

func (s *VaultsStorage) Create(v domain.Vault) error {
	_, err := s.db.Exec(`
        INSERT INTO vaults (name, limit_amount, left_amount, icon, color, expire)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, v.Name, v.Limit, v.LeftAmount, v.Icon, v.Color, v.Expire)
	return err
}

func (s *VaultsStorage) Update(v domain.Vault) error {
	_, err := s.db.Exec(`
        UPDATE vaults SET name = $1, limit_amount = $2, left_amount = $3, icon = $4, color = $5, expire = $6
        WHERE id = $7
    `, v.Name, v.Limit, v.LeftAmount, v.Icon, v.Color, v.Expire, v.ID)
	return err
}
