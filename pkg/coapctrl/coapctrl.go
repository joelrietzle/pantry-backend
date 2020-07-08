package coapctrl

import (
	"bytes"
	"context"
	"time"

	"github.com/go-ocf/go-coap/v2/message"
	"github.com/go-ocf/go-coap/v2/message/codes"
	"github.com/go-ocf/go-coap/v2/mux"
	"github.com/joelrietzle/pantry/pkg/postgres"
	"go.uber.org/zap"
)

type CoAPController struct {
	logger *zap.Logger
	store  *postgres.PantryStore
}

func NewController(logger *zap.Logger, store *postgres.PantryStore) *CoAPController {
	return &CoAPController{
		logger: logger,
		store:  store,
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	list, err := c.store.ListPantries(ctx)
	if err != nil {
		c.logger.Error("Cannot list pantries", zap.Error(err))
	}

	// Incomplete example! Do something with list of pantries instead of logging!
	c.logger.Info("Found pantries", zap.Any("Pantries", list))
}
