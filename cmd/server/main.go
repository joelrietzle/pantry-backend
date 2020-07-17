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

type Variables struct {
	qr        string
	uuid      string
	firstName string
	lastName  string
	email     string
}

func main() {

	var Variable Variables
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

	firstnameMerchant, lastnameMerchant, emailMerchant, _ := merchant.ListMerchants(context.Background(), Variable.firstName, Variable.lastName, Variable.email)
	firstnameCustomer, lastnameCustomer, emailCustomer, _ := customer.ListCustomers(context.Background(), Variable.firstName, Variable.lastName, Variable.email)
	qr, uuid, _ := store.ListDevice(context.Background(), Variable.qr, Variable.uuid)

	logger.Info("Test", zap.Any("firstnameMerchant", firstnameMerchant), zap.Any("lastnameMerchant: ", lastnameMerchant), zap.Any("emailMerchant: ", emailMerchant))
	logger.Info("Test", zap.Any("firstnameCustomer", firstnameCustomer), zap.Any("lastnameCustomer: ", lastnameCustomer), zap.Any("emailCustomer: ", emailCustomer))
	logger.Info("Test", zap.Any("Testing", qr))
	logger.Info("Test", zap.Any("Testing", uuid))

	coapsvc := coapctrl.NewController(logger, store, merchant, customer)

	r := mux.NewRouter()
	r.Use(coapsvc.Logger)
	r.HandleFunc("/pantry/", mux.HandlerFunc(coapsvc.HandlePantry))
	r.HandleFunc("/customer/", mux.HandlerFunc(coapsvc.HandleCustomer))
	r.HandleFunc("/merchant/", mux.HandlerFunc(coapsvc.HandleMerchant))
	r.HandleFunc("/unlock/", mux.HandlerFunc(coapsvc.HandleQRUnlock))
	r.HandleFunc("/lock/", mux.HandlerFunc(coapsvc.HandleQRLock))
	r.HandleFunc("/unlock/", mux.HandlerFunc(coapsvc.HandleUUIDUnlock))
	r.HandleFunc("/lock/", mux.HandlerFunc(coapsvc.HandleUUIDLock))

	logger.Fatal("Failed to start CoAP server", zap.Error(coap.ListenAndServe("udp", ":5688", r)))
}
