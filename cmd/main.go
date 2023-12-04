package main

import (
	badger "github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"log"
	"remote-control/internal"
)

func main() {
	logger := zap.NewExample()
	db, err := badger.Open(badger.DefaultOptions("./db-badger"))
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	srv := internal.InitRouter(logger, db)

	logger.Info("Server started", zap.String("addr", srv.Addr))
	log.Fatal(srv.ListenAndServe())
}
