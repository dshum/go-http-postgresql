package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dshum/go-http-postgresql/internal/db"
	"github.com/dshum/go-http-postgresql/internal/models"
	"github.com/julienschmidt/httprouter"
)

// Welcome handles index route
func Welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	response := makeResponse("message", "Welcome!")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetUsers displays all stored persons from db
func GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := db.GetInstance()

	people := []models.Person{}
	db.Select(&people, "SELECT * FROM persons ORDER BY id ASC")

	response, err := json.Marshal(people)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// GetUser displays stored person by ID
func GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	db := db.GetInstance()

	person := models.Person{}
	db.Get(&person, "SELECT * FROM persons WHERE id = $1", id)

	response, err := json.Marshal(person)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// UpdateUser updates person by ID
func UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	db := db.GetInstance()

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var person models.Person
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	person.ID = idInt

	db.NamedExec("UPDATE persons SET first_name = :first_name, last_name = :last_name, email = :email WHERE id = :id", person)

	response := makeResponse("message", "Person updated")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// CreateUser creates new person
func CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	db := db.GetInstance()

	var person models.Person
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&person)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db.NamedExec("INSERT INTO persons (first_name, last_name, email) VALUES (:first_name, :last_name, :email)", person)

	response := makeResponse("message", "Person created")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// DeleteUser deletes person by ID
func DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	db := db.GetInstance()

	db.MustExec("DELETE FROM persons WHERE id = $1", id)

	response := makeResponse("message", "Person deleted")
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// makeResponse returns JSON with message
func makeResponse(key, value string) []byte {
	message := make(map[string]string)
	message[key] = value
	response, err := json.Marshal(message)

	if err != nil {
		log.Fatal(err)
	}

	return response
}
