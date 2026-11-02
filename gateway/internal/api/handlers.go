package api

import (
	"encoding/json"
	"ledger"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/transactions", createTransaction)
	mux.HandleFunc("GET /api/transactions", listTransactions)
	mux.HandleFunc("POST /api/budgets", createBudget)
	mux.HandleFunc("GET /api/budgets", listBudgets)
	mux.HandleFunc("/", notFound)
	return mux
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func createTransaction(w http.ResponseWriter, r *http.Request) {
	var req CreateTransactionRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный JSON в запросе")
		return
	}

	tx := toTransaction(req)
	err = tx.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ledger.AddTransaction(tx)
	if err == ledger.ErrBudgetExceeded {
		writeError(w, http.StatusConflict, "budget exceeded")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	list := ledger.ListTransactions()
	created := list[len(list)-1]
	writeJSON(w, http.StatusCreated, toTransactionResponse(created))
}

func listTransactions(w http.ResponseWriter, r *http.Request) {
	result := []TransactionResponse{}
	for _, tx := range ledger.ListTransactions() {
		result = append(result, toTransactionResponse(tx))
	}
	writeJSON(w, http.StatusOK, result)
}

func createBudget(w http.ResponseWriter, r *http.Request) {
	var req CreateBudgetRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "неправильный JSON в запросе")
		return
	}

	b := toBudget(req)
	err = b.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = ledger.SetBudget(b)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
		return
	}

	writeJSON(w, http.StatusCreated, toBudgetResponse(b))
}

func listBudgets(w http.ResponseWriter, r *http.Request) {
	result := []BudgetResponse{}
	for _, b := range ledger.ListBudgets() {
		result = append(result, toBudgetResponse(b))
	}
	writeJSON(w, http.StatusOK, result)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound, "not found")
}
