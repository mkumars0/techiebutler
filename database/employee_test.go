package database

import (
	"errors"
	"testing"

	"example.com/m/Assesment/models"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	database := Database{DB: db}

	employee := models.Employee{Name: "John Doe", Position: "Software Engineer", Salary: 70000}

	// success case
	mock.ExpectExec(CreateQuery).
		WithArgs(employee.Name, employee.Position, employee.Salary).
		WillReturnResult(sqlmock.NewResult(1, 1))

	_, err = database.Create(employee)
	if err != nil {
		t.Error(err)
	}

	// lastInsertID error case
	mock.ExpectExec(CreateQuery).
		WithArgs(employee.Name, employee.Position, employee.Salary).
		WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))

	_, err = database.Create(employee)
	if err == nil {
		t.Error(err)
	}

	// error from db case
	mock.ExpectExec(CreateQuery).
		WithArgs(employee.Name, employee.Position, employee.Salary).
		WillReturnError(errors.New("test error"))

	_, err = database.Create(employee)
	if err == nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	database := Database{DB: db}

	employee := models.Employee{ID: 1, Name: "John Doe", Position: "Software Engineer", Salary: 70000}

	// success case
	mock.ExpectQuery(GetQuery).
		WithArgs(employee.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).AddRow(employee.ID, employee.Name, employee.Position, employee.Salary))

	resp, err := database.Get(employee.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp, employee)

	// error case
	mock.ExpectQuery(GetQuery).
		WithArgs(employee.ID).
		WillReturnError(errors.New("test error"))

	_, err = database.Get(employee.ID)
	if err == nil {
		t.Error(err)
	}

	// rowscan error case
	mock.ExpectQuery(GetQuery).
		WithArgs(employee.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).AddRow("apple", employee.Name, employee.Position, employee.Salary))

	_, err = database.Get(employee.ID)
	if err == nil {
		t.Error(err)
	}
}

func TestGetAll(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	database := Database{DB: db}

	page := 1
	pageLimit := 5
	offset := (page - 1) * pageLimit
	employee := models.Employee{ID: 1, Name: "John Doe", Position: "Software Engineer", Salary: 70000}
	resp := []models.Employee{employee}

	// success case
	mock.ExpectQuery(GetAllQuery).
		WithArgs(pageLimit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).AddRow(employee.ID, employee.Name, employee.Position, employee.Salary))

	result, err := database.GetAll(page, pageLimit)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, resp, result)

	// error case
	mock.ExpectQuery(GetAllQuery).
		WithArgs(pageLimit, offset).
		WillReturnError(errors.New("test error"))

	_, err = database.GetAll(page, pageLimit)
	if err == nil {
		t.Error(err)
	}

	// rowscan error case
	mock.ExpectQuery(GetAllQuery).
		WithArgs(pageLimit, offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "position", "salary"}).AddRow("apple", employee.Name, employee.Position, employee.Salary))

	_, err = database.GetAll(page, pageLimit)
	if err == nil {
		t.Error(err)
	}
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	database := Database{DB: db}

	var id int64 = 1

	// success case
	mock.ExpectExec(DeleteQuery).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.Delete(id)
	if err != nil {
		t.Error(err)
	}

	// error from db case
	mock.ExpectExec(DeleteQuery).
		WithArgs(id).
		WillReturnError(errors.New("test error"))

	err = database.Delete(id)
	if err == nil {
		t.Error(err)
	}
}

func TestUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	assert.NoError(t, err)
	defer db.Close()

	database := Database{DB: db}

	var id int64 = 1
	employee := models.Employee{Name: "John Doe", Position: "SDE-2", Salary: 20000}

	// success case
	mock.ExpectExec("update employee set name = ?,position = ?,salary = ? where id = ?").
		WithArgs(employee.Name, employee.Position, employee.Salary, id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = database.Update(employee, id)
	if err != nil {
		t.Error(err)
	}

	// error from db case
	mock.ExpectExec("update employee set name = ?,position = ?,salary = ? where id = ?").
		WithArgs(employee.Name, employee.Position, employee.Salary, id).
		WillReturnError(errors.New("test error"))

	err = database.Update(employee, id)
	if err == nil {
		t.Error(err)
	}
}
