package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	// Настраиваем логирование в stdout, чтобы можно было перенаправить в файл
	logger := log.New(os.Stdout, "", log.LstdFlags)

	// Эндпоинт 1: Возвращает 200 OK
	http.HandleFunc("/api/v1/transactions/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "processed"})
	})

	// Эндпоинт 2: Возвращает 201 Created
	http.HandleFunc("/internal/webhook/billing-sync", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"result": "synced"})
	})

	// Эндпоинт 3: Эмулирует сбой, возвращает 500 и пишет стек-трейс в лог
	http.HandleFunc("/api/v1/transactions/route", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("[ERROR] E2E_Test_04 failed")
		logger.Println("[FATAL] panic: runtime error: invalid memory address or nil pointer dereference")
		logger.Println("goroutine 34 [running]:")
		logger.Println("main.ProcessPayment(...)")
		logger.Println("\t/app/internal/services/payment.go:42 +0x1a4")

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	})

	logger.Println("Starting mock transaction-service on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
