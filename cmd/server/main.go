package main

import (
	"context"
	"os"
	"time"

	"github.com/go-ocf/go-coap/v2"
	"github.com/go-ocf/go-coap/v2/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joelrietzle/pantry/pkg/coapctrl"
	"github.com/joelrietzle/pantry/pkg/postgres"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	pool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to connect to database", zap.Error(err))
	}
	defer pool.Close()

	store := postgres.NewPantryStore(pool)

	pantries, _ := store.ListPantries(context.Background())
	logger.Info("Test", zap.Any("Testing", pantries))

	coapsvc := coapctrl.NewController(logger, store)

	r := mux.NewRouter()
	r.Use(coapsvc.Logger)
	r.Handle("/a", mux.HandlerFunc(coapsvc.HandleA))
	r.Handle("/b", mux.HandlerFunc(coapsvc.HandleB))
	r.Handle("/c", mux.HandlerFunc(coapsvc.HandleC))

	logger.Fatal("Failed to start CoAP server", zap.Error(coap.ListenAndServe("udp", ":5688", r)))
}
