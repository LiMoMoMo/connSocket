package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PTFS/connsocket/connC"
	"github.com/PTFS/connsocket/models"
)

func main() {
	connc, err := connC.NewConnC("8421", "127.0.0.1")
	if err != nil {
		log.Println(err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	//Read
	go func() {
		for {
			select {
			case cmd := <-connc.GetCmdChan():
				fmt.Println("Receive Msg", cmd.Content.(models.AddressInfo))
			case <-ctx.Done():
				return
			}
		}
	}()
	//Write
	for i := 0; i < 1; i++ {
		re := models.Login{
			ID: "qwerty",
		}
		report := models.Report{
			Type:    models.Type_Login,
			Content: re,
		}

		connc.Write(&report)
		time.Sleep(2 * time.Second)
	}

	addr := models.Report{
		Type: models.Type_AddressInfo,
		Content: models.AddressInfo{
			ID:        "qwerty",
			Addresses: []string{"0.0.0.0", "1.1.1.1"},
		},
	}
	connc.Write(&addr)
	time.Sleep(10 * time.Second)

	cancel()
}
