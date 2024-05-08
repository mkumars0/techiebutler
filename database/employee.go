package database

import (
	"database/sql"
	"strings"

	"example.com/m/Assesment/models"
)

type Database struct {
	DB *sql.DB
}

func New(db *sql.DB) Database {
	return Database{DB: db}
}

func (d Database) Create(employee models.Employee) (int64, error) {
	var id int64

	query := CreateQuery
	result, err := d.DB.Exec(query, employee.Name, employee.Position, employee.Salary)
	if err != nil {
		return id, err
	}

	id, err = result.LastInsertId()
	if err != nil {
		return id, err
	}

	return id, err
}

func (d Database) Update(employee models.Employee, id int64) error {
	// building query to accomodate partial update
	query := "update employee set "
	var args []interface{}

	if employee.Name != "" {
		query = query + "name = ?,"
		args = append(args, employee.Name)
	}

	if employee.Position != "" {
		query = query + "position = ?,"
		args = append(args, employee.Position)
	}

	if employee.Salary != 0 {
		query = query + "salary = ?,"
		args = append(args, employee.Salary)
	}

	query = strings.TrimSuffix(query, ",")
	query = query + " where id = ?"
	args = append(args, id)

	_, err := d.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	return err
}

func (d Database) Get(id int64) (models.Employee, error) {
	var employee models.Employee

	rows, err := d.DB.Query(GetQuery, id)
	if err != nil {
		return employee, err
	}

	defer rows.Close()

	if rows.Err() != nil {
		return employee, err
	}

	for rows.Next() {
		err = rows.Scan(&employee.ID, &employee.Name, &employee.Position, &employee.Salary)
		if err != nil {
			return employee, err
		}
	}

	return employee, err
}

func (d Database) GetAll(page, pageLimit int) ([]models.Employee, error) {
	var employee []models.Employee

	offset := (page - 1) * pageLimit

	rows, err := d.DB.Query(GetAllQuery, pageLimit, offset)
	if err != nil {
		return employee, err
	}

	defer rows.Close()

	if rows.Err() != nil {
		return employee, err
	}

	for rows.Next() {
		var e models.Employee
		err = rows.Scan(&e.ID, &e.Name, &e.Position, &e.Salary)
		if err != nil {
			return employee, err
		}

		employee = append(employee, e)
	}

	return employee, err
}

func (d Database) Delete(id int64) error {
	_, err := d.DB.Exec(DeleteQuery, id)
	if err != nil {
		return err
	}

	return err
}
