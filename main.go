package main

import (
	"database/sql"
	"fmt"

	"github.com/ilya-korotya/solid/database/postgres"
	_ "github.com/ilya-korotya/solid/handler"
	"github.com/ilya-korotya/solid/server"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", "postgres://lowcoder:@localhost:5432/solid?sslmode=disable")
	if err != nil {
		panic(fmt.Sprint("Invalid open connect to database:", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprint("Invalid connect to database:", err))
	}
	// store := postgres.NewUserStore(db)
	// usecase := usecase.NewUserInteractor(store)
	server.InstallDB(postgres.NewUserStore(db))
	server.Run(":8081")
}
