package main

import (
	"jwt-auth/cmd/api"
	"jwt-auth/db"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres:64062639@localhost/user_db?sslmode=disable"
	db, err := db.NewSQLStorage(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := api.NewAPIServer(":8383", db)
	server.Run()
}
