package handler

import (
	"encoding/json"
	"monyLonger/internal/domain"
	"net/http"
	"strconv"
	"time"
)

type DashboardHandler struct {
	vRepo domain.VaultRepository
	tRepo domain.TransactionRepository
}

type DashboardData struct {
	UserName        string                `json:"user_name"`
	MonthPeriod     string                `json:"month_period"`
	Balance         int                   `json:"balance"`
	SpentInMonth    int                   `json:"spent_in_month"`
	Incomes         int                   `json:"incomes"`
	OverBudgetCount int                   `json:"over_budget_count"`
	OverBudgetNames []string              `json:"over_budget_names"`
	Budgets         []BudgetsResponse     `json:"budgets"`
	Transactions    []TransactionResponse `json:"transactions"`
	MonthlyStats    []MonthlyStat         `json:"monthly_stats"`
}

type MonthlyStat struct {
	Label   string `json:"label"`
	Income  int    `json:"income"`
	Expense int    `json:"expense"`
}

func NewDashboardHandler(vRepo domain.VaultRepository, tRepo domain.TransactionRepository) *DashboardHandler {
	return &DashboardHandler{vRepo: vRepo, tRepo: tRepo}
}

func (h *DashboardHandler) Load(w http.ResponseWriter, r *http.Request) {
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	v, _ := h.vRepo.GetAll()
	t, _ := h.tRepo.GetForPeriod(period)

	d := DashboardData{
		"Алексей",
		"Июнь 2026",
		80000,
		23180,
		60000,
		1,
		[]string{"Cafe"},
		ConvertToBudgetResponseData(v),
		ConvertToTransactionResponse(t),
		[]MonthlyStat{{"Jan", 60000, 28000}},
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)

}
