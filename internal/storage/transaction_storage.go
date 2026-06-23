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

const txColumns = `SELECT id, name, amount, is_income, vault_id, note, date FROM transactions`

func scanTx(row func(dest ...any) error) (domain.Transaction, error) {
	var t domain.Transaction
	var vaultID sql.NullInt64
	if err := row(&t.ID, &t.Name, &t.Amount, &t.IsIncome, &vaultID, &t.Note, &t.Date); err != nil {
		return t, err
	}
	if vaultID.Valid {
		t.VaultId = int(vaultID.Int64)
	}
	return t, nil
}

func (s *TransactionStorage) GetAll() ([]domain.Transaction, error) {
	rows, err := s.db.Query(txColumns + ` ORDER BY date DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		t, err := scanTx(rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *TransactionStorage) GetForPeriod(period time.Time) ([]domain.Transaction, error) {
	rows, err := s.db.Query(txColumns+`
        WHERE EXTRACT(YEAR FROM date) = $1
        AND EXTRACT(MONTH FROM date) = $2
        ORDER BY date DESC`,
		period.Year(), period.Month())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Transaction
	for rows.Next() {
		t, err := scanTx(rows.Scan)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}
	return result, nil
}

func (s *TransactionStorage) GetByID(id int) (*domain.Transaction, error) {
	t, err := scanTx(s.db.QueryRow(txColumns+` WHERE id = $1`, id).Scan)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *TransactionStorage) Create(t domain.Transaction) error {
	vaultID := sql.NullInt64{Int64: int64(t.VaultId), Valid: t.VaultId != 0}
	_, err := s.db.Exec(`
        INSERT INTO transactions (name, amount, is_income, vault_id, note, date)
        VALUES ($1, $2, $3, $4, $5, $6)
    `, t.Name, t.Amount, t.IsIncome, vaultID, t.Note, t.Date)
	return err
}

func (s *TransactionStorage) Update(t domain.Transaction) error {
	vaultID := sql.NullInt64{Int64: int64(t.VaultId), Valid: t.VaultId != 0}
	_, err := s.db.Exec(`
        UPDATE transactions SET name=$1, amount=$2, is_income=$3, vault_id=$4, note=$5, date=$6
        WHERE id=$7
    `, t.Name, t.Amount, t.IsIncome, vaultID, t.Note, t.Date, t.ID)
	return err
}

func (s *TransactionStorage) Delete(id int) error {
	_, err := s.db.Exec(`DELETE FROM transactions WHERE id = $1`, id)
	return err
}
