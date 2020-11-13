package models

type Person struct {
	ID        int    `db:"id" json:"id"`
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
	Email     string `db:"email" json:"email"`
}

var schema = `
CREATE TABLE persons (
	id int,
    first_name text,
    last_name text,
    email text
)`
