package models

//AgazaType structure defining the user, it has name, id, department, remaining_sick_leaves, remaining_annual_leave, leaves_taken
type AgazaType struct {
	ID   string `json:"id"    yaml:"id"`
	Name string `json:"name"     yaml:"name"`
}
