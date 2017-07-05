package models

import "time"

//Leave structure defining the leave taken by users, it contains, ID, userID, Reason, from and to
type Leave struct {
	ID           string    `json:"ID"    yaml:"ID"`
	UserID       string    `json:"user_id"     yaml:"user_id"`
	Reason       string    `json:"reason" yaml:"reason"`
	From         time.Time `json:"from" yaml:"from"`
	To           time.Time `json:"to" yaml:"to"`
	Status       string    `json:"status" yaml:"status"`
	DepartmentID string    `json:"department_id" yaml:"department_id"`
	UserName     string    `json:"user_name" yaml:"user_name"`
	Type         int       `json:"type" yaml:"type"`
}

//ExposedLeave is used to return the leave for listing APIs
type ExposedLeave struct {
	ID           string    `json:"ID"    yaml:"ID"`
	UserID       string    `json:"user_id"     yaml:"user_id"`
	From         time.Time `json:"from" yaml:"from"`
	To           time.Time `json:"to" yaml:"to"`
	DepartmentID string    `json:"department_id" yaml:"department_id"`
	UserName     string    `json:"user_name" yaml:"user_name"`
}
