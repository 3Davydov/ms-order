package main

import (
	"log"

	"github.com/3Davydov/ms-order/config"
	"github.com/3Davydov/ms-order/internal/adapters/db"
	"github.com/3Davydov/ms-order/internal/adapters/grpc"
	"github.com/3Davydov/ms-order/internal/application/core/api"
)

func main() {
	DbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database : %v", err)
	}

	application := api.NewApplication(DbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
