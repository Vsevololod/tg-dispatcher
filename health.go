package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"tg-dispatcher/domain"
	"tg-dispatcher/lib/logger/sl"
	"time"
)

type DatabaseChecker interface {
	IsDBOk() (string, error)
}
type AmqpChecker interface {
	IsQueueOk() (string, error)
}

type Health struct {
	srv         *http.Server
	dbChecker   DatabaseChecker
	amqpChecker AmqpChecker
	log         *slog.Logger
}

func NewHealth(dbChecker DatabaseChecker, amqpChecker AmqpChecker, log *slog.Logger) *Health {
	srv := &http.Server{Addr: ":8080"}
	return &Health{srv: srv, dbChecker: dbChecker, amqpChecker: amqpChecker, log: log}
}

func (h *Health) Start() {
	http.HandleFunc("/health", h.healthHandler)
	go func() {
		h.log.Info("Starting health server on :8080")
		if err := h.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			h.log.Error("Server failed:", sl.Err(err))
		}
	}()
}

func (h *Health) StopProcessing() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := h.srv.Shutdown(ctx); err != nil {
		h.log.Error("Server forced to shutdown: %v", sl.Err(err))
	}
}

func (h *Health) healthHandler(w http.ResponseWriter, r *http.Request) {

	textError := ""
	dbOk, err := h.dbChecker.IsDBOk()
	if err != nil {
		h.log.Error("db check failed", sl.Err(err))
		textError = err.Error()
	}

	amqpOk, err := h.amqpChecker.IsQueueOk()
	if err != nil {
		h.log.Error("amqp check failed", sl.Err(err))
		textError += err.Error()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	status := domain.AppStatus{AppStatus: "UP", DbStatus: dbOk, AmqpStatus: amqpOk, TextError: textError}
	body, err := json.Marshal(status)
	if err != nil {
		_, _ = w.Write([]byte(`{"status": "UP"}`))
	}
	_, _ = w.Write(body)
}
