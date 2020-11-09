package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var instance *sqlx.DB

func Connect(driver, data string) *sqlx.DB {
	db, err := sqlx.Connect(driver, data)
	if err != nil {
		log.Fatalln(err)
	}
	instance = db
	return instance
}

func GetInstance() *sqlx.DB {
	return instance
}
