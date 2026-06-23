package main

import (
	"log"
	"monyLonger/internal/domain"
	"monyLonger/internal/handler"
	"monyLonger/internal/storage"
	"net/http"
)

func main() {

	db, err := storage.NewDB("file://migrations")
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}
	defer db.Close()

	events := make(chan domain.TransactionCreatedEvent, 100)

	txStorage := storage.NewTransactionStorage(db)
	vltStorage := storage.NewVaultsStorage(db)

	txHandler := handler.NewTransactionHandler(txStorage, events)
	vltHandler := handler.NewVaultHandler(vltStorage)
	dHandler := handler.NewDashboardHandler(vltStorage, txStorage)

	http.HandleFunc("GET /api/dashboard", dHandler.Load)

	http.HandleFunc("GET /api/budgets", vltHandler.GetAll)
	http.HandleFunc("POST /api/budgets", vltHandler.Create)
	http.HandleFunc("PATCH /api/budgets", vltHandler.Update)

	http.HandleFunc("POST /api/transactions", txHandler.Create)
	http.HandleFunc("GET /api/transactions", txHandler.GetAll)

	http.HandleFunc("/budgets", pageHandler("./public/budgets.html"))
	http.HandleFunc("/transactions/new", pageHandler("./public/transactions-new.html"))
	http.HandleFunc("/transactions", pageHandler("./public/transactions.html"))

	http.Handle("/", http.FileServer(http.Dir("./public")))

	go func() {
		for event := range events {
			vltHandler.UpdateByEvent(event)
		}
	}()

	http.ListenAndServe(":8090", nil)
}

func pageHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
