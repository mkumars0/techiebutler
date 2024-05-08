package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"example.com/m/Assesment/database"
	"example.com/m/Assesment/models"
	"github.com/gorilla/mux"
)

type Handler struct {
	EmployeeDB database.Employee
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	err = json.Unmarshal(data, &employee)
	if err != nil {
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	// checking mandatory fields
	employee.Name = strings.TrimSpace(employee.Name)
	if employee.Name == "" {
		http.Error(w, "error employee name missing", http.StatusBadRequest)
		return
	}

	employee.Position = strings.TrimSpace(employee.Position)
	if employee.Name == "" {
		http.Error(w, "error employee name missing", http.StatusBadRequest)
		return
	}

	if employee.Salary == 0 {
		http.Error(w, "error employee name missing", http.StatusBadRequest)
		return
	}

	id, err := h.EmployeeDB.Create(employee)
	if err != nil {
		http.Error(w, "error creating employee", http.StatusInternalServerError)
		return
	}

	employee.ID = id

	response, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Write(response)
}

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringID := vars["id"]

	intID, err := strconv.Atoi(stringID)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	// check for empty id
	id := int64(intID)
	if id == 0 {
		http.Error(w, "error empty id", http.StatusBadRequest)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
		return
	}

	var employee models.Employee
	err = json.Unmarshal(data, &employee)
	if err != nil {
		http.Error(w, "error unmarshalling body", http.StatusBadRequest)
		return
	}

	// mandatory check for atleast one field
	employee.Name = strings.TrimSpace(employee.Name)
	employee.Position = strings.TrimSpace(employee.Position)
	if employee.Name == "" && employee.Position == "" && employee.Salary == 0 {
		http.Error(w, "error no fields to update", http.StatusBadRequest)
		return
	}

	err = h.EmployeeDB.Update(employee, id)
	if err != nil {
		http.Error(w, "error creating employee", http.StatusInternalServerError)
		return
	}

	employee, err = h.EmployeeDB.Get(id)
	if err != nil {
		http.Error(w, "error fetching empoyee details", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Write(response)
}

func (h Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringID := vars["id"]

	intID, err := strconv.Atoi(stringID)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	// check for empty id
	id := int64(intID)
	if id == 0 {
		http.Error(w, "error empty id", http.StatusBadRequest)
		return
	}

	employee, err := h.EmployeeDB.Get(id)
	if err != nil {
		http.Error(w, "error fetching empoyee details", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(employee)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Write(response)
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	pageParam := queryParams.Get("page")

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		http.Error(w, "error invalid page value "+pageParam, http.StatusBadRequest)
		return
	}

	pageLimitParam := queryParams.Get("pagelimit")

	pageLimit, err := strconv.Atoi(pageLimitParam)
	if err != nil {
		http.Error(w, "error invalid pagelimit", http.StatusBadRequest)
		return
	}

	employees, err := h.EmployeeDB.GetAll(page, pageLimit)
	if err != nil {
		http.Error(w, "error fetching all empoyee details", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(employees)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Write(response)
}

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	stringID := vars["id"]

	intID, err := strconv.Atoi(stringID)
	if err != nil {
		http.Error(w, "error invalid id", http.StatusBadRequest)
		return
	}

	// check for empty id
	id := int64(intID)
	if id == 0 {
		http.Error(w, "error empty id", http.StatusBadRequest)
		return
	}

	err = h.EmployeeDB.Delete(id)
	if err != nil {
		http.Error(w, "error deleting employee", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal("employee deleted sucessfully")
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusBadRequest)
		return
	}

	w.Write(response)
}
