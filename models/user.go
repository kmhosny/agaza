package models

//User structure defining the user, it has name, id, department, remaining_sick_leaves, remaining_annual_leave, leaves_taken
type User struct {
	ID                  string   `json:"ID"    yaml:"ID"`
	Name                string   `json:"name"     yaml:"name"`
	Department          string   `json:"department" yaml:"department"`
	RemainingSickLeaves int      `json:"remaining_sick_leave" yaml:"remaining_sick_leave"`
	TakenLeaves         []string `json:"leaves" yaml:"leaves"`
}
