package main

import (
	"log"
	"net/http"

	"github.com/dshum/go-http-postgresql/pkg/service"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=db port=5432 user=postgres password=secret dbname=go_test1 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userService := service.NewUsersService(db)
	router := httprouter.New()

	router.GET("/", service.Welcome)
	router.GET("/users", userService.GetUsers)
	router.POST("/users", userService.CreateUser)
	router.GET("/users/:id", userService.GetUser)
	router.PUT("/users/:id", userService.UpdateUser)
	router.DELETE("/users/:id", userService.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
