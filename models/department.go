package models

//Department is the structure for the departments, it has an id, name and a color
type Department struct {
	Name  string `json:"name"            yaml:"name"`
	Color string `json:"color"         yaml:"color"`
	ID    string `json:"id"        yaml:"id"`
}
