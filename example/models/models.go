package examplemodels

import (
	"github.com/LiMoMoMo/go-connSocket/models"
)

func init() {
	models.RegisterReport(Type_Addr, func() models.BaseContent { return &Addr{} })
	models.RegisterReport(Type_Show, func() models.BaseContent { return &Show{} })
	models.RegisterReport(Type_HeartBeat, func() models.BaseContent { return &HeatbeatInfo{} })
}

const (
	Type_Addr models.ReportType = models.Type_Logout + 1 + iota
	Type_Show
	Type_HeartBeat
)

// Addr struct
type Addr struct {
	ID   string
	Name string
}

type Test struct {
	A string
	B int
	C []string
}

type Show struct {
	Name string
	Te   Test
}

type HeatbeatInfo struct {
	ID            string
	ParentTraffic []TrafficInfo
	ChildTraffics []ChildTraffic
}

type ChildTraffic []TrafficInfo

type TrafficInfo struct {
	ID           string
	TotalTraffic int64
	Speed        int64
}

// Fill implement from BaseContent
func (a *Addr) Fill(m map[string]interface{}) {
}

// Fill implement from BaseContent
func (s *Show) Fill(m map[string]interface{}) {
}

func (s *HeatbeatInfo) Fill(m map[string]interface{}) {
}
