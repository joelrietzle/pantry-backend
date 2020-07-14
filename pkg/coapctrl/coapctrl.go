package coapctrl

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-ocf/go-coap/v2/message"
	"github.com/go-ocf/go-coap/v2/message/codes"
	"github.com/go-ocf/go-coap/v2/mux"
	"github.com/joelrietzle/pantry/pkg/postgres"
	"go.uber.org/zap"
)

var qr string

type CoAPController struct {
	logger   *zap.Logger
	store    *postgres.PantryStore
	merchant *postgres.MerchantStore
	customer *postgres.CustomerStore
}

func NewController(logger *zap.Logger, store *postgres.PantryStore, merchant *postgres.MerchantStore, customer *postgres.CustomerStore) *CoAPController {
	return &CoAPController{
		logger:   logger,
		store:    store,
		merchant: merchant,
		customer: customer,
	}
}

func (c *CoAPController) Logger(next mux.Handler) mux.Handler {
	return mux.HandlerFunc(func(w mux.ResponseWriter, r *mux.Message) {
		c.logger.Info("Received request",
			zap.Any("RemoteAddr", w.Client().RemoteAddr()),
			zap.Any("Request", r),
		)
		next.ServeCOAP(w, r)
	})
}

func (c *CoAPController) HandleA(w mux.ResponseWriter, r *mux.Message) {
	err := w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader([]byte("hello world")))
	if err != nil {
		c.logger.Error("Cannot set response", zap.Error(err))
	}
}

func (c *CoAPController) HandleB(w mux.ResponseWriter, r *mux.Message) {
	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte("B hello world")),
	}

	optsBuf := make([]byte, 32)

	opts, used, err := customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	if err == message.ErrTooSmall {
		optsBuf = append(optsBuf, make([]byte, used)...)
		opts, _, err = customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	}

	if err != nil {
		c.logger.Error("Cannot set options to response", zap.Error(err))
		return
	}

	customResp.Options = opts

	err = w.Client().WriteMessage(&customResp)
	if err != nil {
		c.logger.Error("Cannot set response", zap.Error(err))
	}
}

func (c *CoAPController) HandleC(w mux.ResponseWriter, r *mux.Message) {

	path, err := r.Options.Path()
	QR := strings.TrimPrefix(path, "unlock/")
	c.logger.Info("Trimmed: ", zap.Any("Path", QR))
	qrstring := postgres.getString(QR)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	listPantry, err := c.store.ListPantries(ctx)
	if err != nil {
		c.logger.Error("Cannot list pantries", zap.Error(err))
	}

	listMerchant, err := c.merchant.ListMerchants(ctx)
	if err != nil {
		c.logger.Error("Cannot list merchants", zap.Error(err))
	}

	listCustomer, err := c.customer.ListCustomers(ctx)
	if err != nil {
		c.logger.Error("Cannot list merchants", zap.Error(err))
	}

	searchPantry, uuid, err := c.store.SearchPantry(ctx)
	if err != nil {
		c.logger.Error("Cannot find corresponding UUID", zap.Error(err))
	}

	// Incomplete example! Do something with list of pantries instead of logging!
	c.logger.Info("Found pantries BOY", zap.Any("Pantries", listPantry))
	c.logger.Info("Found merchants", zap.Any("Merchants", listMerchant))
	c.logger.Info("Found customers", zap.Any("Customers", listCustomer))
	c.logger.Info("Found UUID", zap.Any("UUID", uuid))
	c.logger.Info("Searched pantry", zap.Any("Search", searchPantry))

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("List of pantries %v", listPantry))),
	}

	optsBuf := make([]byte, 32)

	opts, used, err := customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	if err == message.ErrTooSmall {
		optsBuf = append(optsBuf, make([]byte, used)...)
		opts, _, err = customResp.Options.SetContentFormat(optsBuf, message.TextPlain)
	}

	if err != nil {
		c.logger.Error("Cannot set options to response", zap.Error(err))
		return
	}

	customResp.Options = opts
	c.logger.Info("BODY: ", zap.Any("body: ", customResp.Body))
	err = w.Client().WriteMessage(&customResp)
	if err != nil {
		c.logger.Error("Cannot set pantry response", zap.Error(err))
	}

}
