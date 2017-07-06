package models

//Department is the structure for the departments, it has an id, name and a color
type Department struct {
	DepartmentName  string `json:"name"            yaml:"name"`
	DepartmentColor string `json:"color"         yaml:"color"`
	ID              string `json:"id"        yaml:"id"`
}
