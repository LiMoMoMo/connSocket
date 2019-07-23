package models

import "encoding/json"

// ReportType type of report Message
type ReportType int32

const (
	// Type_Login client connect
	Type_Login ReportType = iota
	// Type_Logout client disconnect
	Type_Logout
)

// Report post message
type Report struct {
	Type    ReportType
	Content interface{}
}

// Login client connect
type Login struct {
	ID string
}

// Fill implement from BaseContent
func (l *Login) Fill(m map[string]interface{}) {
	FillStruct(m, l)
}

// String struct to bytearray
func (r *Report) String() ([]byte, error) {
	return json.Marshal(r)
}

// Unmarshal bytearray to struct
func (r *Report) Unmarshal(content json.RawMessage) {
	var val interface{}
	json.Unmarshal(content, &val)
	i, ok := reportMap[r.Type]
	if ok {
		t := i()
		t.Fill(val.(map[string]interface{}))
		r.Content = (t).(BaseContent)
	} else {
		r.Content = val
	}
}
