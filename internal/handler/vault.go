package handler

import (
	"encoding/json"
	"log"
	"monyLonger/internal/domain"
	"net/http"
	"strconv"
	"time"
)

type VaultHandler struct {
	repo domain.VaultRepository
}

type VaultResponse struct {
	UserName    string            `json:"user_name"`
	MonthPeriod string            `json:"month_period"`
	TotalLimit  int               `json:"total_limit"`
	TotalSpent  int               `json:"total_spent"`
	Budgets     []BudgetsResponse `json:"budgets"`
}

type BudgetsResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Color       string `json:"color"`
	Spent       int    `json:"spent"`
	Limit       int    `json:"limit"`
	PeriodLabel string `json:"period_label"`
}

func NewVaultHandler(repo domain.VaultRepository) *VaultHandler {
	return &VaultHandler{repo: repo}
}

func (h *VaultHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	vaults, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var totalLimit, totalSpent int
	for _, v := range vaults {
		totalLimit += v.Limit
		totalSpent += v.Limit - v.LeftAmount
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VaultResponse{
		UserName:    "Алексей",
		MonthPeriod: monthLabel(time.Now()),
		TotalLimit:  totalLimit,
		TotalSpent:  totalSpent,
		Budgets:     ConvertToBudgetResponseData(vaults),
	})
}

func (h *VaultHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var v domain.Vault
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if v.LeftAmount == 0 {
		v.LeftAmount = v.Limit
	}
	if v.Expire.IsZero() {
		now := time.Now()
		v.Expire = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	}

	if err := h.repo.Create(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *VaultHandler) UpdateByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	var v domain.Vault
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	existing, err := h.repo.GetById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		http.NotFound(w, r)
		return
	}
	v.ID = id
	v.LeftAmount = existing.LeftAmount
	v.Expire = existing.Expire

	if err := h.repo.Update(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *VaultHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *VaultHandler) UpdateByEvent(event domain.TransactionCreatedEvent) {
	v, err := h.repo.GetById(event.VaultID)
	if err != nil {
		log.Printf("UpdateByEvent: GetById %d: %v", event.VaultID, err)
		return
	}
	if v == nil {
		log.Printf("UpdateByEvent: vault %d not found", event.VaultID)
		return
	}
	if event.IsIncome {
		v.LeftAmount += event.Amount
	} else {
		v.LeftAmount -= event.Amount
	}
	if err := h.repo.Update(*v); err != nil {
		log.Printf("UpdateByEvent: update vault %d: %v", event.VaultID, err)
	}
}

func monthLabel(t time.Time) string {
	names := [...]string{"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь",
		"Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	return names[t.Month()] + " " + strconv.Itoa(t.Year())
}

func ConvertToBudgetResponseData(vaults []domain.Vault) []BudgetsResponse {
	result := make([]BudgetsResponse, 0, len(vaults))
	for _, v := range vaults {
		result = append(result, BudgetsResponse{
			Id:          v.ID,
			Name:        v.Name,
			Icon:        v.Icon,
			Color:       v.Color,
			Spent:       v.Limit - v.LeftAmount,
			Limit:       v.Limit,
			PeriodLabel: v.Expire.Month().String(),
		})
	}
	return result
}
