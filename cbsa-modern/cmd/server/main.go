package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/cicsdev/cbsa-modern/internal/config"
	"github.com/cicsdev/cbsa-modern/internal/database"
	"github.com/cicsdev/cbsa-modern/internal/handler"
	"github.com/cicsdev/cbsa-modern/internal/middleware"
	"github.com/cicsdev/cbsa-modern/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		migrationsDir = "/migrations"
	}
	if err := database.RunMigrations(db, migrationsDir); err != nil {
		slog.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}

	customerSvc := service.NewCustomerService(db, cfg.SortCode)
	accountSvc := service.NewAccountService(db, cfg.SortCode)
	txnSvc := service.NewTransactionService(db, cfg.SortCode)
	creditSvc := service.NewCreditService()

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/company", func(w http.ResponseWriter, r *http.Request) {
			handler.JSON(w, http.StatusOK, map[string]string{"name": cfg.CompanyName})
		})
		r.Get("/sortcode", func(w http.ResponseWriter, r *http.Request) {
			handler.JSON(w, http.StatusOK, map[string]string{"sortCode": cfg.SortCode})
		})

		customerHandler := handler.NewCustomerHandler(customerSvc, creditSvc)
		r.Route("/customers", func(r chi.Router) {
			r.Get("/", customerHandler.List)
			r.Post("/", customerHandler.Create)
			r.Get("/{customerNumber}", customerHandler.Get)
			r.Put("/{customerNumber}", customerHandler.Update)
			r.Delete("/{customerNumber}", customerHandler.Delete)
		})

		accountHandler := handler.NewAccountHandler(accountSvc)
		r.Route("/accounts", func(r chi.Router) {
			r.Get("/", accountHandler.List)
			r.Post("/", accountHandler.Create)
			r.Get("/{accountNumber}", accountHandler.Get)
			r.Put("/{accountNumber}", accountHandler.Update)
			r.Delete("/{accountNumber}", accountHandler.Delete)
			r.Get("/customer/{customerNumber}", accountHandler.ListByCustomer)
		})

		txnHandler := handler.NewTransactionHandler(txnSvc)
		r.Route("/transactions", func(r chi.Router) {
			r.Get("/", txnHandler.List)
			r.Put("/debit", txnHandler.Debit)
			r.Put("/credit", txnHandler.Credit)
			r.Put("/transfer", txnHandler.Transfer)
		})
	})

	slog.Info("starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
