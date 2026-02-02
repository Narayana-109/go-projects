package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

var DB *pgxpool.Pool

func Connect() {
	var connString = "postgres://postgres:root@localhost:5432/trailog?sslmode=disable"
	fmt.Print("Connecting..")
	var err error
	DB, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal("unable to use data source name", err)
	}
}
