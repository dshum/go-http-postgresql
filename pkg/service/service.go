package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/dshum/go-http-postgresql/pkg/errors"
	"github.com/dshum/go-http-postgresql/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type UsersService struct {
	db *sqlx.DB
}

type ResponseWelcomeMessage struct {
	Service string `json:"service"`
	Version string `json:"version"`
	Message string `json:"message"`
}

type ResponseMessage struct {
	Message string        `json:"message"`
	Person  models.Person `json:"person"`
}

type ResponseError struct {
	Error string `json:"error"`
}

// NewUsersService is called in main() to initalize UsersService and pass it a reference to db object
func NewUsersService(db *sqlx.DB) *UsersService {
	return &UsersService{db: db}
}

func newResponseWelcomeMessage(service, version, message string) *ResponseWelcomeMessage {
	return &ResponseWelcomeMessage{Service: service, Version: version, Message: message}
}

func newResponseMessage(message string, person models.Person) *ResponseMessage {
	return &ResponseMessage{Message: message, Person: person}
}

func newResponseError(err error) *ResponseError {
	return &ResponseError{Error: err.Error()}
}

// Welcome handles index route
func Welcome(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponseWelcomeMessage("http-postgresql", "0.1.0", "Welcome!"))
}

// GetUsers displays all stored persons from db
func (service *UsersService) GetUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	people := []models.Person{}
	err := service.db.Select(&people, "SELECT * FROM persons ORDER BY id ASC")
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(people)
}

// GetUser displays stored person by ID
func (service *UsersService) GetUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	person := models.Person{}
	err = service.db.Get(&person, "SELECT * FROM persons WHERE id = $1 LIMIT 1", idInt)
	if err != nil {
		err := errors.PersonNotFound("Person not found")
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

// UpdateUser updates person by ID
func (service *UsersService) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	var person models.Person

	err = service.db.Get(&person, "SELECT * FROM persons WHERE id = $1 LIMIT 1", idInt)
	if err != nil {
		err := errors.PersonNotFound("Person not found")
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	_, err = service.db.NamedExec("UPDATE persons SET first_name = :first_name, last_name = :last_name, email = :email WHERE id = :id", person)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponseMessage("Person updated", person))
}

// CreateUser creates new person
func (service *UsersService) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var person models.Person

	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	sql := "INSERT INTO persons (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING id"
	err := service.db.QueryRow(sql, person.FirstName, person.LastName, person.Email).Scan(&person.ID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponseMessage("Person created", person))
}

// DeleteUser deletes person by ID
func (service *UsersService) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")

	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	person := models.Person{}
	err = service.db.Get(&person, "SELECT * FROM persons WHERE id = $1 LIMIT 1", idInt)
	if err != nil {
		err := errors.PersonNotFound("Person not found")
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	_, err = service.db.NamedExec("DELETE FROM persons WHERE id = :id", person)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(newResponseError(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newResponseMessage("Person deleted", person))
}
