package main

import (
	"fmt"

	connS "github.com/PTFS/connsocket/connS"
	exModels "github.com/PTFS/connsocket/example/models"
	"github.com/PTFS/connsocket/models"
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
			case models.Type_Logout:
				fmt.Println("Disconnected.")
			}
		}
	}
}

func main() {
	s := connS.NewConnS("8421", "0.0.0.0")
	s.Start()
	go readMsg(s)
	select {}
}
