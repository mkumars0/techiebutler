package database

import (
	"example.com/m/Assesment/models"
	"github.com/stretchr/testify/mock"
)

type MockDatabase struct {
	mock.Mock
	CreateF func(employee models.Employee) (int64, error)
	UpdateF func(employee models.Employee, id int64) error
	GetF    func(id int64) (models.Employee, error)
	GetAllF func(page, pageLimit int) ([]models.Employee, error)
	DeleteF func(id int64) error
}

func (m *MockDatabase) Create(employee models.Employee) (int64, error) {
	return m.CreateF(employee)
}

func (m *MockDatabase) Update(employee models.Employee, id int64) error {
	return m.UpdateF(employee, id)
}

func (m *MockDatabase) Get(id int64) (models.Employee, error) {
	return m.GetF(id)
}

func (m *MockDatabase) GetAll(page, pageLimit int) ([]models.Employee, error) {
	return m.GetAllF(page, pageLimit)
}

func (m *MockDatabase) Delete(id int64) error {
	return m.DeleteF(id)
}
