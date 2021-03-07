package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"google.golang.org/grpc/grpclog"

	"github.com/ubozov/grpc-atlant/data"
	"github.com/ubozov/grpc-atlant/grpc/server"
)

type config struct {
	db   *data.Config
	addr string
}

func getConfig() (*config, error) {
	return &config{
		addr: os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		db: &data.Config{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			DBName:   os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}, nil
}

func main() {
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	conf, err := getConfig()
	if err != nil {
		log.Fatalln("Failed to read .env:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, err := data.NewDB(ctx, *conf.db)
	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}
	defer db.Close(ctx)

	if err != server.Start(db, log, conf.addr) {
		log.Fatalln("Failed to serve:", err)
	}
}
