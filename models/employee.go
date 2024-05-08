package models

type Employee struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Position string  `josn:"position"`
	Salary   float64 `json:"salary"`
}
