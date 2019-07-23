package examplemodels

import (
	"encoding/json"

	"github.com/PTFS/connsocket/models"
)

func init() {
	models.RegisterReport(Type_Addr, func() models.BaseContent { return &Addr{} })
	models.RegisterReport(Type_Show, func() models.BaseContent { return &Show{} })
}

const (
	Type_Addr models.ReportType = models.Type_Logout + 1 + iota
	Type_Show
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

// Fill implement from BaseContent
func (a *Addr) Fill(m map[string]interface{}) {
	// models.FillStruct(m, a)
	bytes, _ := json.Marshal(m)
	json.Unmarshal(bytes, a)
}

// Fill implement from BaseContent
func (s *Show) Fill(m map[string]interface{}) {
	// models.FillStruct(m, s)
	bytes, _ := json.Marshal(m)
	json.Unmarshal(bytes, s)
}
