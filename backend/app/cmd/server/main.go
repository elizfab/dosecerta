package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"dose-certa-backend/config"
	"dose-certa-backend/internal/handler"
	"dose-certa-backend/internal/middleware"
	"dose-certa-backend/internal/repository/mongodb"
	"dose-certa-backend/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("falha ao conectar no MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("erro ao desconectar MongoDB: %v", err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB não responde: %v", err)
	}
	log.Printf("conectado ao MongoDB: %s (db: %s)", cfg.MongoURI, cfg.DBName)

	db := client.Database(cfg.DBName)

	mux := http.NewServeMux()

	// Recurso: Medication (/api/v1/medications)
	medicationRepo := mongodb.NewMedicationRepository(db)
	medicationUC := usecase.NewMedicationUseCase(medicationRepo)
	handler.NewMedicationHandler(medicationUC).RegisterRoutes(mux)

	// Recurso: Medidas corporais (/api/v1/weight-records)
	weightRepo := mongodb.NewWeightRepository(db)
	weightUC := usecase.NewWeightUseCase(weightRepo)
	handler.NewWeightHandler(weightUC).RegisterRoutes(mux)

	// Recurso: Exam (/api/v1/exams)
	examRepo := mongodb.NewExamRepository(db)
	examUC := usecase.NewExamUseCase(examRepo)
	handler.NewExamHandler(examUC).RegisterRoutes(mux)

	// Recurso: Doses tomadas (/api/v1/doses)
	doseRepo := mongodb.NewDoseLogRepository(db)
	doseUC := usecase.NewDoseLogUseCase(doseRepo)
	handler.NewDoseLogHandler(doseUC).RegisterRoutes(mux)

	// Recurso: Configurações key/value (/api/v1/settings) — ex.: meta de peso
	settingRepo := mongodb.NewSettingRepository(db)
	settingUC := usecase.NewSettingUseCase(settingRepo)
	handler.NewSettingHandler(settingUC).RegisterRoutes(mux)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"` + cfg.DBName + `"}`))
	})

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      middleware.CORS(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("servidor iniciado na porta %s (env: %s)", cfg.ServerPort, cfg.AppEnv)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("erro ao iniciar servidor: %v", err)
	}
}
