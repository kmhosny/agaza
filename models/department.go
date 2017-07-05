package models

//Department is the structure for the departments, it has an id, name and a color
type Department struct {
	DepartmentName  string `json:"department_name"            yaml:"department_name"`
	DepartmentColor string `json:"department_color"         yaml:"department_color"`
	ID              string `json:"id"        yaml:"id"`
}
