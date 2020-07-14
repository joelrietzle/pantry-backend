package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/go-ocf/go-coap/v2/udp"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	co, err := udp.Dial("localhost:5688")
	if err != nil {
		logger.Fatal("Error dialing", zap.Error(err))
	}

	path := "/"

	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := co.Get(ctx, path)
	if err != nil {
		logger.Fatal("Error sending request", zap.Error(err))
	}

	body, err := ioutil.ReadAll(resp.Body())

	logger.Info("Response payload", zap.ByteString("response", body))
	//logger.Info("Reponse body", zap.String("reponse", resp.Body()))
}
