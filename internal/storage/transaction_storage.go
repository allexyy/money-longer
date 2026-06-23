package storage

import (
	"database/sql"
	"monyLonger/internal/domain"
	"time"
)

type TransactionStorage struct {
	db *sql.DB
}

func NewTransactionStorage(db *sql.DB) *TransactionStorage {
	return &TransactionStorage{db: db}
}

func (s *TransactionStorage) GetAll() ([]domain.Transaction, error) {
	rows, err := s.db.Query(`
        SELECT id, name, amount, is_income, vault_id, note, date
        FROM transactions
        ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.Name, &t.Amount, &t.IsIncome, &t.VaultId, &t.Note, &t.Date); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *TransactionStorage) GetForPeriod(period time.Time) ([]domain.Transaction, error) {
	rows, err := s.db.Query(`
        SELECT id, name, amount, is_income, vault_id, note, date
        FROM transactions
        WHERE EXTRACT(YEAR FROM date) = $1
        AND EXTRACT(MONTH FROM date) = $2
        ORDER BY date DESC
    `, period.Year(), period.Month())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(&t.ID, &t.Name, &t.Amount, &t.IsIncome, &t.VaultId, &t.Note, &t.Date); err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *TransactionStorage) Create(t domain.Transaction) error {
	_, err := s.db.Exec(`
        INSERT INTO transactions (name, amount, is_income, vault_id, note, date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, t.Name, t.Amount, t.IsIncome, t.VaultId, t.Note, t.Date)
	return err
}

func (s *TransactionStorage) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM transactions WHERE id = $1`, id)
	return err
}
