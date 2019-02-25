package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ilya-korotya/solid/database/postgres"
	_ "github.com/ilya-korotya/solid/handler"
	"github.com/ilya-korotya/solid/server"
	"github.com/ilya-korotya/solid/usecase"
	_ "github.com/lib/pq"
)

func main() {
	done := make(chan struct{})
	db, err := sql.Open("postgres", "postgres://lowcoder:@localhost:5432/solid?sslmode=disable")
	if err != nil {
		panic(fmt.Sprint("Invalid open connect to database:", err))
	}
	if err := db.Ping(); err != nil {
		panic(fmt.Sprint("Invalid connect to database:", err))
	}
	store := postgres.NewUserStore(db)
	userUsecase := usecase.NewUserInteractor(store)
	server.InstallUserUsecase(userUsecase)
	if err := server.Run(":8080", done); err != nil {
		log.Println("Server is broke:", err)
		// if we call Shutdown or Close
		if err == http.ErrServerClosed {
			<-done
		}
	}
}
