package main

//
import (
	"database/sql"
	"fmt"

	"github.com/ilya-korotya/solid/server"

	"github.com/ilya-korotya/solid/database/postgres"
	"github.com/ilya-korotya/solid/server/handler"
	"github.com/ilya-korotya/solid/server/router"
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
	dbi := postgres.NewUserStore(db)
	uc := usecase.NewUserInteractor(dbi)
	hc := handler.New(uc)
	router := &router.Config{hc}
	server := server.Config{
		Port: "8080",
	}
	if err := server.Run(router.Install()); err != nil {
		fmt.Println("Cannot run server:", err)
	}
}

/** TODO: creare custom context as:
struct {
	http Context
}
func Log
func Response
and other

Create middle ware, when be reaturn basic context
and set it in our struct
*/
