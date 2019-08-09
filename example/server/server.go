package main

import (
	"fmt"

	connS "github.com/LiMoMoMo/go-connSocket/connS"
	exModels "github.com/LiMoMoMo/go-connSocket/example/models"
	"github.com/LiMoMoMo/go-connSocket/models"
)

func readMsg(s *connS.ConnS) {
	for {
		select {
		case msg := <-s.GetRepChan():
			switch msg.Type {
			case exModels.Type_Addr:
				addr := msg.Content.(*exModels.Addr)
				fmt.Println("Receive Report Msg is", addr)
				cmd := models.Command{
					Type: models.Command_Start,
					Content: models.Start{
						Val: "This is AddrCommand",
					},
				}
				s.WriteTo(addr.ID, &cmd)
			case exModels.Type_Show:
				show := msg.Content.(*exModels.Show)
				fmt.Println("Receive Report Msg is", show.Te)
			case exModels.Type_HeartBeat:
				hb := msg.Content.(*exModels.HeatbeatInfo)
				fmt.Println("Receive Report Msg is", hb)
				fmt.Println("Parent len is", len(hb.ParentTraffic))
				fmt.Println("Child len is", len(hb.ChildTraffics))
			case models.Type_Logout:
				hb := msg.Content.(models.Register)
				fmt.Println("Disconnected.", hb)
			}
		}
	}
}

func main() {
	s := connS.NewConnS("8976", "0.0.0.0")
	s.Start()
	go readMsg(s)
	select {}
}
