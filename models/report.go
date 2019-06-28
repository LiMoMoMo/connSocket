package models

import "encoding/json"

type ReportType int32

const (
	// Login
	Type_Login  ReportType = 0
	Type_Logout ReportType = 1
)

// Report post message
type Report struct {
	Type    ReportType
	Content string
}

type Login struct {
	ID string
}

func (r *Report) String() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Login) String() ([]byte, error) {
	return json.Marshal(r)
}
