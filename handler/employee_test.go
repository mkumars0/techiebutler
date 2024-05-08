package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/m/Assesment/database"
	"example.com/m/Assesment/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	testCases := []struct {
		name           string
		body           models.Employee
		response       models.Employee
		result         int64
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful Create Request",
			body:           models.Employee{Name: "John", Position: "SDE", Salary: 30000},
			response:       models.Employee{ID: 1, Name: "John", Position: "SDE", Salary: 30000},
			err:            nil,
			result:         1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from db",
			body:           models.Employee{Name: "John", Position: "SDE", Salary: 30000},
			err:            errors.New("TestError"),
			result:         0,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Mandatory salary check error",
			body:           models.Employee{Name: "John", Position: "SDE"},
			err:            nil,
			result:         0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Mandatory name check error",
			body:           models.Employee{Position: "SDE"},
			err:            nil,
			result:         0,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Mandatory position check error",
			body:           models.Employee{Name: "John"},
			err:            nil,
			result:         0,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mock for dependency
			testDatabase := new(database.MockDatabase)
			testDatabase.CreateF = func(models.Employee) (int64, error) {
				return tc.result, tc.err
			}

			mockHandler := Handler{EmployeeDB: testDatabase}

			data, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatal(err)
			}

			// Create the request
			req, err := http.NewRequest(http.MethodPost, "employee", bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mockHandler.Create(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				resp := models.Employee{}
				data, err = io.ReadAll(rr.Body)
				if err != nil {
					t.Error(err)
				}

				err = json.Unmarshal(data, &resp)
				if err != nil {
					t.Error(string(data))
					t.Error(err)
				}

				assert.Equal(t, tc.response, resp)
			}

		})
	}
}

func TestUpdate(t *testing.T) {
	testCases := []struct {
		name           string
		body           models.Employee
		response       models.Employee
		id             string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful Update Request",
			body:           models.Employee{Name: "John", Position: "SDE-2", Salary: 30000},
			response:       models.Employee{ID: 1, Name: "John", Position: "SDE-2", Salary: 30000},
			err:            nil,
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from db",
			body:           models.Employee{Name: "John", Position: "SDE", Salary: 30000},
			err:            errors.New("TestError"),
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Empty update check",
			body:           models.Employee{},
			err:            nil,
			id:             "1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty id check",
			body:           models.Employee{},
			err:            nil,
			id:             "0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid id check",
			body:           models.Employee{},
			err:            nil,
			id:             "a",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mock for dependency
			testDatabase := new(database.MockDatabase)
			testDatabase.UpdateF = func(models.Employee, int64) error {
				return tc.err
			}

			testDatabase.GetF = func(id int64) (models.Employee, error) {
				return tc.response, nil
			}

			mockHandler := Handler{EmployeeDB: testDatabase}

			data, err := json.Marshal(tc.body)
			if err != nil {
				t.Fatal(err)
			}

			// Create the request
			req, err := http.NewRequest(http.MethodPut, "employee", bytes.NewReader(data))
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.id,
			})

			rr := httptest.NewRecorder()
			mockHandler.Update(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				resp := models.Employee{}
				data, err = io.ReadAll(rr.Body)
				if err != nil {
					t.Error(err)
				}

				err = json.Unmarshal(data, &resp)
				if err != nil {
					t.Error(string(data))
					t.Error(err)
				}

				assert.Equal(t, tc.response, resp)
			}

		})
	}
}

func TestGet(t *testing.T) {
	testCases := []struct {
		name           string
		response       models.Employee
		id             string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful Get Request",
			response:       models.Employee{ID: 1, Name: "John", Position: "SDE-2", Salary: 30000},
			err:            nil,
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from db",
			err:            errors.New("TestError"),
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Empty id check",
			err:            nil,
			id:             "0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid id check",
			err:            nil,
			id:             "a",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mock for dependency
			testDatabase := new(database.MockDatabase)

			testDatabase.GetF = func(id int64) (models.Employee, error) {
				return tc.response, tc.err
			}

			mockHandler := Handler{EmployeeDB: testDatabase}

			// Create the request
			req, err := http.NewRequest(http.MethodPut, "employee", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.id,
			})

			rr := httptest.NewRecorder()
			mockHandler.Get(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				resp := models.Employee{}
				data, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Error(err)
				}

				err = json.Unmarshal(data, &resp)
				if err != nil {
					t.Error(string(data))
					t.Error(err)
				}

				assert.Equal(t, tc.response, resp)
			}

		})
	}
}

func TestGetAll(t *testing.T) {
	testCases := []struct {
		name           string
		response       []models.Employee
		queryParams    string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful Get Request",
			response:       []models.Employee{{ID: 1, Name: "John", Position: "SDE-2", Salary: 30000}},
			err:            nil,
			queryParams:    "?page=2&pagelimit=20",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from db",
			err:            errors.New("TestError"),
			queryParams:    "?page=2&pagelimit=20",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Invalid page check",
			err:            nil,
			queryParams:    "?page=apple&pagelimit=20",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid pagelimit check",
			err:            nil,
			queryParams:    "?page=2&pagelimit=apple",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mock for dependency
			testDatabase := new(database.MockDatabase)

			testDatabase.GetAllF = func(page, pageLimit int) ([]models.Employee, error) {
				return tc.response, tc.err
			}

			mockHandler := Handler{EmployeeDB: testDatabase}

			// Create the request
			req, err := http.NewRequest(http.MethodPut, "employee"+tc.queryParams, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			mockHandler.GetAll(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				resp := []models.Employee{}
				data, err := io.ReadAll(rr.Body)
				if err != nil {
					t.Error(err)
				}

				err = json.Unmarshal(data, &resp)
				if err != nil {
					t.Error(string(data))
					t.Error(err)
				}

				assert.Equal(t, tc.response, resp)
			}

		})
	}
}

func TestDelete(t *testing.T) {
	testCases := []struct {
		name           string
		id             string
		err            error
		expectedStatus int
	}{
		{
			name:           "Successful Delete Request",
			err:            nil,
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Error from db",
			err:            errors.New("TestError"),
			id:             "1",
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Empty id check",
			err:            nil,
			id:             "0",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid id check",
			err:            nil,
			id:             "a",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//mock for dependency
			testDatabase := new(database.MockDatabase)

			testDatabase.DeleteF = func(id int64) error {
				return tc.err
			}

			mockHandler := Handler{EmployeeDB: testDatabase}

			// Create the request
			req, err := http.NewRequest(http.MethodDelete, "employee", nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{
				"id": tc.id,
			})

			rr := httptest.NewRecorder()
			mockHandler.Delete(rr, req)

			// Check the response status code
			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

		})
	}
}
