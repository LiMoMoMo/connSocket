package main

import (
	"fmt"

	connS "github.com/PTFS/connsocket/connS"
	"github.com/PTFS/connsocket/models"
)

func readMsg(s *connS.ConnS) {
	for {
		select {
		case msg := <-s.GetRepChan():
			switch msg.Type {
			case models.Type_AddressInfo:
				addrinfo := (msg.Content).(models.AddressInfo)
				fmt.Println(addrinfo)
				// test echo
				cmd := models.Command{
					Type:    models.Command_Connect,
					Content: addrinfo,
				}
				s.WriteTo(addrinfo.ID, &cmd)
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
