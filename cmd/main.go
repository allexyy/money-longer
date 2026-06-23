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

	txHandler := handler.NewTransactionHandler(txStorage, vltStorage, events)
	vltHandler := handler.NewVaultHandler(vltStorage)
	dHandler := handler.NewDashboardHandler(vltStorage, txStorage)

	// Dashboard
	http.HandleFunc("GET /api/dashboard", dHandler.Load)

	// Budgets
	http.HandleFunc("GET /api/budgets", vltHandler.GetAll)
	http.HandleFunc("POST /api/budgets", vltHandler.Create)
	http.HandleFunc("PUT /api/budgets/{id}", vltHandler.UpdateByID)
	http.HandleFunc("DELETE /api/budgets/{id}", vltHandler.Delete)

	// Transactions
	http.HandleFunc("GET /api/transactions", txHandler.GetAll)
	http.HandleFunc("POST /api/transactions", txHandler.Create)
	http.HandleFunc("GET /api/transactions/{id}", txHandler.GetByID)
	http.HandleFunc("PUT /api/transactions/{id}", txHandler.Update)
	http.HandleFunc("DELETE /api/transactions/{id}", txHandler.Delete)

	// Pages
	http.HandleFunc("GET /budgets", pageHandler("./public/budgets.html"))
	http.HandleFunc("GET /transactions/new", pageHandler("./public/transactions-new.html"))
	http.HandleFunc("GET /transactions/{id}/edit", pageHandler("./public/transactions-edit.html"))
	http.HandleFunc("GET /transactions", pageHandler("./public/transactions.html"))

	http.Handle("/", http.FileServer(http.Dir("./public")))

	go func() {
		for event := range events {
			vltHandler.UpdateByEvent(event)
		}
	}()

	log.Println("Listening on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func pageHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	}
}
