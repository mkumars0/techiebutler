package main

import (
	"database/sql"
	"net/http"

	"example.com/m/Assesment/database"
	"example.com/m/Assesment/handler"
	"github.com/gorilla/mux"
)

func main() {

	// connecting to db
	dsn := "username:password@/dbname"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return
	}
	defer db.Close()

	empDB := database.New(db)
	eh := handler.Handler{EmployeeDB: empDB}

	r := mux.NewRouter()

	r.HandleFunc("/employee/{id}", eh.Get).Methods(http.MethodGet)
	r.HandleFunc("/employee/", eh.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/employee", eh.Create).Methods(http.MethodPost)
	r.HandleFunc("/employee/{id}", eh.Update).Methods(http.MethodPut)
	r.HandleFunc("/employee/{id}", eh.Delete).Methods(http.MethodDelete)

}
