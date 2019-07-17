package models

import "encoding/json"

type CommandType int32

const (
	// Command_Talk
	Command_Talk CommandType = 0

	// 下发节点地址信息给目的节点
	Command_Connect CommandType = 1

	// 指挥节点从源站开始/停止拉流
	Command_PullStream CommandType = 2
)

// Report post message
type Command struct {
	Type    CommandType
	Content interface{}
}

type Talk struct {
	Val string
}

// 从源站拉流指令信息结构
// 从源站开始拉流指令内容：Url为源站地址，Type为Type_Start
// 停止从源站拉流指令内容：Url为空，Type为Type_Stop
type CodeType int32

const (
	// 开始
	Type_Start CodeType = 0
	// 停止
	Type_Stop CodeType = 1
)

type PullStream struct {
	Url  string
	Type CodeType
}

// end CodeType

func (r *Command) String() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Command) Unmarshal(content json.RawMessage) {
	switch r.Type {
	case Command_Connect:
		info := AddressInfo{}
		json.Unmarshal(content, &info)
		r.Content = info
	case Command_PullStream:
		stream := PullStream{}
		json.Unmarshal(content, &stream)
		r.Content = stream
	}
}
