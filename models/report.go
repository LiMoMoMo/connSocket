package models

import "encoding/json"

type ReportType int32

const (
	// Login
	Type_Login  ReportType = 0
	Type_Logout ReportType = 1
	// 上报的地址信息
	Type_AddressInfo ReportType = 2
	// 上报的节点状态信息
	Type_NodeStatInfo ReportType = 3
	// 上报的下级节点的连接状态信息
	Type_ConnectionInfo ReportType = 4
)

// Report post message
type Report struct {
	Type    ReportType
	Content interface{}
}

type Login struct {
	ID string
}

type AddressInfo struct {
	ID            string
	P2PAddresses  []string
	LiveAddresses []string
}

type NodeStatInfo struct {
	ID   string
	Info StatInfo
}

type StatInfo struct {
	InnerIp   string
	OuterIp   string
	Longitude float64
	Latitude  float64
	ISP       string
	City      string
	Province  string
	SpeedUp   float64
	SpeedDown float64
	Load      float64
}

type STATU int32

const (
	// 连接
	Statu_Connected STATU = iota
	// 断开连接
	Statu_Disconnected
)

// 上报给后台的下级节点的连接状态信息
type ConnectionInfo struct {
	Self  string // 节点自己的ID
	Peer  string // 下级节点ID
	Statu STATU  // 下级节点的连接状态
}

func (r *Report) String() ([]byte, error) {
	return json.Marshal(r)
}

func (r *Report) Unmarshal(content json.RawMessage) {
	switch r.Type {
	case Type_Login:
		login := Login{}
		json.Unmarshal(content, &login)
		r.Content = login
	case Type_AddressInfo:
		info := AddressInfo{}
		json.Unmarshal(content, &info)
		r.Content = info
	case Type_NodeStatInfo:
		status := NodeStatInfo{}
		json.Unmarshal(content, &status)
		r.Content = status
	case Type_ConnectionInfo:
		coninfo := ConnectionInfo{}
		json.Unmarshal(content, &coninfo)
		r.Content = coninfo
	}
}
