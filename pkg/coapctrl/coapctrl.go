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

/*func (c *CoAPController) HandleA(w mux.ResponseWriter, r *mux.Message) {
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
}*/

func (c *CoAPController) HandlePantry(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	pantry := strings.TrimPrefix(path, "pantry/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	qr, uuid, err := c.store.ListDevice(ctx, pantry, pantry)
	if err != nil {
		c.logger.Error("Cannot find corresponding UUID", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Pantry details: QR: %v, UUID: %v", qr, uuid))),
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

func (c *CoAPController) HandleCustomer(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	customer := strings.TrimPrefix(path, "customer/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	firstnameCustomer, lastnameCustomer, emailCustomer, err := c.customer.ListCustomers(ctx, customer, customer, customer)
	if err != nil {
		c.logger.Error("Cannot list customers", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Customer details: FirstName: %v, LastName: %v, Email: %v", firstnameCustomer, lastnameCustomer, emailCustomer))),
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
		c.logger.Error("Cannot set customer response", zap.Error(err))
	}
}

func (c *CoAPController) HandleMerchant(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	merchant := strings.TrimPrefix(path, "merchant/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	firstname, lastname, email, err := c.merchant.ListMerchants(ctx, merchant, merchant, merchant)
	if err != nil {
		c.logger.Error("Cannot list customers", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Merchant details: FirstName: %v, LastName: %v, Email: %v", firstname, lastname, email))),
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
		c.logger.Error("Cannot set merchants response", zap.Error(err))
	}
}

func (c *CoAPController) HandleQRUnlock(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	qr := strings.TrimPrefix(path, "unlock/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uuid, err := c.store.LockUnlockQR(ctx, qr)
	if err != nil {
		c.logger.Error("Cannot list customers", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Unlock pantry with UUID: %v", uuid))),
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
		c.logger.Error("Cannot set merchants response", zap.Error(err))
	}
}

func (c *CoAPController) HandleQRLock(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	qr := strings.TrimPrefix(path, "lock/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	uuid, err := c.store.LockUnlockQR(ctx, qr)
	if err != nil {
		c.logger.Error("Cannot list customers", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Lock pantry with UUID: %v", uuid))),
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
		c.logger.Error("Cannot set merchants response", zap.Error(err))
	}
}

func (c *CoAPController) HandleUUIDUnlock(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	uuid := strings.TrimPrefix(path, "unlock/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	qr, err := c.store.LockUnlockUUID(ctx, uuid)
	if err != nil {
		c.logger.Error("Cannot unlock pantry", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Lock pantry with QR: %v", qr))),
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
		c.logger.Error("Cannot set unlock response", zap.Error(err))
	}
}

func (c *CoAPController) HandleUUIDLock(w mux.ResponseWriter, r *mux.Message) {
	path, err := r.Options.Path()
	uuid := strings.TrimPrefix(path, "lock/")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	qr, err := c.store.LockUnlockUUID(ctx, uuid)
	if err != nil {
		c.logger.Error("Cannot lock pantry", zap.Error(err))
	}

	customResp := message.Message{
		Code:    codes.Content,
		Token:   r.Token,
		Context: r.Context,
		Options: make(message.Options, 0, 16),
		Body:    bytes.NewReader([]byte(fmt.Sprintf("Lock pantry with QR: %v", qr))),
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
		c.logger.Error("Cannot set lock response", zap.Error(err))
	}
}
