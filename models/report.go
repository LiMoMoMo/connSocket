package models

import "encoding/json"

// ReportType type of report Message
type ReportType int32

const (
	// Type_Login client connect
	Type_Register ReportType = iota
	// Type_Logout client disconnect
	Type_Logout
)

// Report post message
type Report struct {
	Type    ReportType
	Content interface{}
}

// Register client connect
type Register struct {
	ID string
}

// Fill implement from BaseContent
func (l *Register) Fill(m map[string]interface{}) {
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
		bytes, _ := json.Marshal(val.(map[string]interface{}))
		json.Unmarshal(bytes, t)
		r.Content = (t).(BaseContent)
	} else {
		r.Content = val
	}
}
