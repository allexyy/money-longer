package handler

import (
	"encoding/json"
	"monyLonger/internal/domain"
	"net/http"
	"sort"
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
	repo      domain.TransactionRepository
	vaultRepo domain.VaultRepository
	events    chan domain.TransactionCreatedEvent
}

func NewTransactionHandler(
	repo domain.TransactionRepository,
	vaultRepo domain.VaultRepository,
	events chan domain.TransactionCreatedEvent,
) *TransactionHandler {
	return &TransactionHandler{repo: repo, vaultRepo: vaultRepo, events: events}
}

func (h *TransactionHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	vaults, _ := h.vaultRepo.GetAll()
	vaultMap := vaultsByID(vaults)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildTransactionsResponse(transactions, vaultMap, vaults))
}

func (h *TransactionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	t, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if t == nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(t)
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var t domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if t.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if t.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	if !t.IsIncome && t.VaultId != 0 {
		h.events <- domain.TransactionCreatedEvent{VaultID: t.VaultId, Amount: t.Amount, IsIncome: t.IsIncome}
	}
}

func (h *TransactionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var t domain.Transaction
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if t.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if t.Amount <= 0 {
		http.Error(w, "amount must be positive", http.StatusBadRequest)
		return
	}
	t.ID = id
	if err := h.repo.Update(t); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TransactionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func vaultsByID(vaults []domain.Vault) map[int]domain.Vault {
	m := make(map[int]domain.Vault, len(vaults))
	for _, v := range vaults {
		m[v.ID] = v
	}
	return m
}

func toTransactionResponse(t domain.Transaction, vaultMap map[int]domain.Vault) TransactionResponse {
	v := vaultMap[t.VaultId]
	dateKey := ""
	if len(t.Date) >= 7 {
		dateKey = t.Date[:7]
	}
	return TransactionResponse{
		Id:          t.ID,
		Name:        t.Name,
		Icon:        v.Icon,
		Note:        t.Note,
		BudgetId:    t.VaultId,
		BudgetName:  v.Name,
		BudgetColor: v.Color,
		Amount:      t.Amount,
		IsIncome:    t.IsIncome,
		Date:        t.Date,
		DateLabel:   dateKey,
		DateDisplay: t.Date,
	}
}

func buildTransactionsResponse(
	transactions []domain.Transaction,
	vaultMap map[int]domain.Vault,
	vaults []domain.Vault,
) TransactionsResponse {
	txResp := make([]TransactionResponse, 0, len(transactions))
	monthSet := make(map[string]MonthLabel)

	for _, t := range transactions {
		tr := toTransactionResponse(t, vaultMap)
		txResp = append(txResp, tr)

		if tr.DateLabel != "" {
			if _, exists := monthSet[tr.DateLabel]; !exists {
				year, _ := strconv.Atoi(tr.DateLabel[:4])
				month, _ := strconv.Atoi(tr.DateLabel[5:7])
				names := [...]string{"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
					"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
				lbl := ""
				if month >= 1 && month <= 12 {
					lbl = names[month] + " " + strconv.Itoa(year)
				}
				monthSet[tr.DateLabel] = MonthLabel{Value: tr.DateLabel, Label: lbl}
			}
		}
	}

	months := make([]MonthLabel, 0, len(monthSet))
	for _, m := range monthSet {
		months = append(months, m)
	}
	sort.Slice(months, func(i, j int) bool { return months[i].Value > months[j].Value })

	return TransactionsResponse{
		UserName:     "Алексей",
		Transactions: txResp,
		Budgets:      ConvertToBudgetResponseData(vaults),
		Months:       months,
	}
}
