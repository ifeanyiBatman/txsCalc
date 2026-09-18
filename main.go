package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type BankTx struct {
	TransactionID   string    `json:"transactionID"`
	TransactionType string    `json:"transactionType"`
	TransactionDate time.Time `json:"transactionDate"`
	Amount          int       `json:"transactionAmount"`
	Sender          string    `json:"sender"`
	Recipient       string    `json:"recipient"`
	Notes           string    `json:"notes"`
}

type ResponseTotal struct {
	SnapshotDate time.Time `json:"snapshotDate"`
	IncomeTotal  int       `json:"incomeTotal"`
	ExpenseTotal int       `json:"expenseTotal"`
}

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("POST /transactions/total", calculateTotal)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      serveMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func calculateTotal(w http.ResponseWriter, req *http.Request) {
	var txs []BankTx

	req.Body = http.MaxBytesReader(w, req.Body, 1<<20)
	if err := json.NewDecoder(req.Body).Decode(&txs); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	income, expenses := calculateTxs(txs)
	response := ResponseTotal{
		SnapshotDate: time.Now(),
		IncomeTotal:  income,
		ExpenseTotal: expenses,
	}

	dat, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dat)
}

func calculateTxs(txs []BankTx) (income, expenses int) {
	for _, tx := range txs {
		switch tx.TransactionType {
		case "credit", "income", "incoming transfer":
			income += tx.Amount
		case "debit", "expense", "interbank transfer":
			expenses += tx.Amount
		}
	}
	return income, expenses
}
