package models

import "encoding/json"

// CommandType type of command Message
type CommandType int32

const (
	// Command_Start start test
	Command_Start CommandType = iota
)

// Command post message
type Command struct {
	Type    CommandType
	Content interface{}
}

// Start Command Message test
type Start struct {
	Val string
}

// Fill implement from BaseContent
func (s *Start) Fill(m map[string]interface{}) {
	FillStruct(m, s)
}

// String struct to bytearray
func (r *Command) String() ([]byte, error) {
	return json.Marshal(r)
}

// Unmarshal bytearray to struct
func (r *Command) Unmarshal(content json.RawMessage) {
	var val interface{}
	json.Unmarshal(content, &val)
	i, ok := commandMap[r.Type]
	if ok {
		t := i()
		bytes, _ := json.Marshal(val.(map[string]interface{}))
		json.Unmarshal(bytes, t)
		r.Content = (t).(BaseContent)
	} else {
		r.Content = val
	}
}
