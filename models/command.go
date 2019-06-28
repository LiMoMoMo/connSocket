package models

import "encoding/json"

type CommandType int32

const (
	// Command_Talk
	Command_Talk CommandType = 0
)

// Report post message
type Command struct {
	Type    CommandType
	Content string
}

type Talk struct {
	Val string
}

func (r *Command) String() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Talk) String() ([]byte, error) {
	return json.Marshal(r)
}
