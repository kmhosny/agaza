package models

//User structure defining the user, it has name, id, department, remaining_sick_leaves, remaining_annual_leave, leaves_taken
type User struct {
	ID                    string `json:"ID"    yaml:"ID"`
	Name                  string `json:"name"     yaml:"name"`
	DepartmentID          string `json:"department_id" yaml:"department_id"`
	RemainingAnnualLeaves int    `json:"remaining_annual_leave" yaml:"remaining_annual_leave"`
}
