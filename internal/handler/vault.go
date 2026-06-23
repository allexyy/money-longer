package handler

import (
	"encoding/json"
	"fmt"
	"monyLonger/internal/domain"
	"net/http"
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
	}
	w.Header().Set("Content-Type", "application/json")
	response := VaultResponse{
		"Алексей",
		"Июнь 2026",
		80000,
		23180,
		ConvertToBudgetResponseData(vaults),
	}
	json.NewEncoder(w).Encode(response)
}

func (h *VaultHandler) Create(w http.ResponseWriter, r *http.Request) {
	var v domain.Vault
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := h.repo.Create(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *VaultHandler) UpdateByEvent(event domain.TransactionCreatedEvent) {
	fmt.Println("process event")
	v, _ := h.repo.GetById(event.VaultID)
	v.LeftAmount = v.LeftAmount - event.Amount
	if err := h.repo.Update(*v); err != nil {
		fmt.Println(err)
		return
	}
}
func (h *VaultHandler) Update(w http.ResponseWriter, r *http.Request) {
	var v domain.Vault
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	if err := h.repo.Update(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func ConvertToBudgetResponseData(vaults []domain.Vault) []BudgetsResponse {
	var vr []BudgetsResponse
	for _, v := range vaults {
		vr = append(vr, BudgetsResponse{
			v.ID,
			v.Name,
			v.Icon,
			v.Color,
			v.LeftAmount,
			v.Limit,
			v.Expire.Month().String(),
		})
	}
	return vr
}
