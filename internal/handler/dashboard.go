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
	now := time.Now()
	month, _ := strconv.Atoi(r.URL.Query().Get("month"))
	year, _ := strconv.Atoi(r.URL.Query().Get("year"))
	if month < 1 || month > 12 {
		month = int(now.Month())
	}
	if year < 2000 || year > 2100 {
		year = now.Year()
	}
	period := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	vaults, err := h.vRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transactions, err := h.tRepo.GetForPeriod(period)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var spentInMonth, incomes int
	for _, t := range transactions {
		if t.IsIncome {
			incomes += t.Amount
		} else {
			spentInMonth += t.Amount
		}
	}

	overBudgetNames := []string{}
	var totalLeftAmount int
	for _, v := range vaults {
		totalLeftAmount += v.LeftAmount
		if v.LeftAmount <= 0 {
			overBudgetNames = append(overBudgetNames, v.Name)
		}
	}

	vaultMap := vaultsByID(vaults)

	d := DashboardData{
		UserName:        "Алексей",
		MonthPeriod:     monthLabel(period),
		Balance:         totalLeftAmount,
		SpentInMonth:    spentInMonth,
		Incomes:         incomes,
		OverBudgetCount: len(overBudgetNames),
		OverBudgetNames: overBudgetNames,
		Budgets:         ConvertToBudgetResponseData(vaults),
		Transactions:    toTransactionResponses(transactions, vaultMap),
		MonthlyStats:    []MonthlyStat{{Label: monthLabel(period), Income: incomes, Expense: spentInMonth}},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func toTransactionResponses(txs []domain.Transaction, vaultMap map[int]domain.Vault) []TransactionResponse {
	result := make([]TransactionResponse, 0, len(txs))
	for _, t := range txs {
		result = append(result, toTransactionResponse(t, vaultMap))
	}
	return result
}
