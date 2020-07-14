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
	merchant := postgres.NewMerchantStore(pool)
	customer := postgres.NewCustomerStore(pool)
	qrstring := postgres.NewQrString(QR)

	pantries, _ := store.ListPantries(context.Background())
	merchants, _ := merchant.ListMerchants(context.Background())
	customers, _ := customer.ListCustomers(context.Background())
	uuid, search, _ := store.SearchPantry(context.Background())
	qrstring, _ := qrstring.GetString(context.Background())
	logger.Info("Test", zap.Any("Testing", pantries))
	logger.Info("Test", zap.Any("Testing", merchants))
	logger.Info("Test", zap.Any("Testing", customers))
	logger.Info("Test", zap.Any("Testing", uuid))
	logger.Info("Tets", zap.Any("Testing", search))

	coapsvc := coapctrl.NewController(logger, store, merchant, customer, qrstring)

	r := mux.NewRouter()
	r.Use(coapsvc.Logger)
	r.Handle("/a", mux.HandlerFunc(coapsvc.HandleA))
	r.Handle("/b", mux.HandlerFunc(coapsvc.HandleB))
	r.HandleFunc("/unlock/", mux.HandlerFunc(coapsvc.HandleC))
	//r.Handle("/unlock", mux.HandlerFunc(coapsvc.HandleC))

	logger.Fatal("Failed to start CoAP server", zap.Error(coap.ListenAndServe("udp", ":5688", r)))
}
