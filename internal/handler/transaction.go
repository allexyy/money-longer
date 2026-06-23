package handler

import (
	"encoding/json"
	"fmt"
	"monyLonger/internal/domain"
	"net/http"
	"strconv"
)

type TransactionsResponse struct {
	UserName     string                `json:"user_name"`
	Transactions []TransactionResponse `json:"transactions"`
	Budgets      []BudgetsResponse     `json:"budgets"`
	Months       []MonthLabel          `json:"months"`
}

type TransactionResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Note        string `json:"note"`
	BudgetId    int    `json:"budget_id"`
	BudgetName  string `json:"budget_name"`
	BudgetColor string `json:"budget_color"`
	Amount      int    `json:"amount"`
	IsIncome    bool   `json:"is_income"`
	Date        string `json:"date"`
	DateLabel   string `json:"date_label"`
	DateDisplay string `json:"date_display"`
}

type MonthLabel struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type TransactionHandler struct {
	repo   domain.TransactionRepository
	events chan domain.TransactionCreatedEvent
}

func NewTransactionHandler(repo domain.TransactionRepository, events chan domain.TransactionCreatedEvent) *TransactionHandler {
	return &TransactionHandler{repo: repo, events: events}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ConvertToTransactionsResponse(transactions))
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := h.repo.Create(t); err != nil {
		fmt.Println(t)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	h.events <- domain.TransactionCreatedEvent{t.VaultId, t.Amount}
}

func ConvertToTransactionsResponse(transactions []domain.Transaction) TransactionsResponse {
	return TransactionsResponse{
		"",
		ConvertToTransactionResponse(transactions),
		[]BudgetsResponse{{
			1,
			"",
			"",
			"",
			1000,
			2300,
			""},
		},
		[]MonthLabel{{
			"",
			""},
		},
	}
}

func ConvertToTransactionResponse(transactions []domain.Transaction) []TransactionResponse {
	var resp []TransactionResponse
	for _, t := range transactions {
		id, _ := strconv.Atoi(t.ID)
		resp = append(resp, TransactionResponse{
			id,
			t.Name,
			"",
			t.Note,
			t.VaultId,
			"",
			"",
			t.Amount,
			t.IsIncome,
			t.Date,
			t.Date,
			t.Date,
		})
	}
	return resp
}
