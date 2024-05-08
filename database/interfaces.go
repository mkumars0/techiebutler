package database

import "example.com/m/Assesment/models"

type Employee interface {
	Create(employee models.Employee) (int64, error)
	Update(employee models.Employee, id int64) error
	Get(id int64) (models.Employee, error)
	GetAll(page, pageLimit int) ([]models.Employee, error)
	Delete(id int64) error
}
