package main

import (
	"log"
	"net/http"

	"github.com/dshum/go-http-postgresql/internal/db"
	"github.com/dshum/go-http-postgresql/internal/server"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

type Person struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
}

func main() {
	db := db.Connect("postgres", "host=db port=5432 user=postgres password=secret dbname=go_test1 sslmode=disable")
	defer db.Close()

	router := httprouter.New()

	router.GET("/", server.Welcome)
	router.GET("/users", server.GetUsers)
	router.POST("/users", server.CreateUser)
	router.GET("/users/:id", server.GetUser)
	router.PUT("/users/:id", server.UpdateUser)
	router.DELETE("/users/:id", server.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}
