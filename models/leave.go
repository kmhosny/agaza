package models

import (
	"encoding/json"
	"time"
)

//Leave structure defining the leave taken by users, it contains, ID, userID, Reason, from and to
type Leave struct {
	ID           string    `json:"id"    yaml:"id"`
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
	ID           string    `json:"id"    yaml:"id"`
	UserID       string    `json:"user_id"     yaml:"user_id"`
	From         time.Time `json:"from" yaml:"from"`
	To           time.Time `json:"to" yaml:"to"`
	DepartmentID string    `json:"department_id" yaml:"department_id"`
	UserName     string    `json:"user_name" yaml:"user_name"`
}

//MarshalJSON custom marshall
func (l *Leave) MarshalJSON() ([]byte, error) {
	type Alias Leave
	return json.Marshal(&struct {
		*Alias
		From string `json:"from"`
		To   string `json:"to"`
	}{
		Alias: (*Alias)(l),
		From:  l.From.Format("2006-Jan-01"),
		To:    l.To.Format("2006-Jan-01"),
	})
}

//MarshalJSON custom marshall
func (l *ExposedLeave) MarshalJSON() ([]byte, error) {
	type Alias ExposedLeave
	return json.Marshal(&struct {
		*Alias
		From string `json:"from"`
		To   string `json:"to"`
	}{
		Alias: (*Alias)(l),
		From:  l.From.Format("2006-Jan-01"),
		To:    l.To.Format("2006-Jan-01"),
	})
}
