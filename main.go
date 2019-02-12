package main

import (
	"database/sql"
	"fmt"

	"github.com/ilya-korotya/solid/database/postgres"
	"github.com/ilya-korotya/solid/server"
	"github.com/ilya-korotya/solid/usecase"
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
	store := postgres.NewUserStore(db)
	usecase := usecase.NewUserInteractor(store)
	server := server.Server{
		Port:     "8080",
		Handlers: usecase,
	}
	server.Run()
}
