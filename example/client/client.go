package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PTFS/connsocket/connC"
	exModels "github.com/PTFS/connsocket/example/models"
	"github.com/PTFS/connsocket/models"
)

func main() {
	connc, err := connC.NewConnC("8421", "127.0.0.1")
	if err != nil {
		log.Println(err)
		return
	}
	connc.SetReconnect(func() {
		re := models.Login{
			ID: "qwerty",
		}
		report := models.Report{
			Type:    models.Type_Login,
			Content: re,
		}
		connc.Write(&report)
	})
	ctx, cancel := context.WithCancel(context.Background())
	//Read
	go func() {
		for {
			select {
			case cmd := <-connc.GetCmdChan():
				fmt.Println("Receive Msg", cmd.Content.(*models.Start).Val)
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
	for i := 0; i < 20; i++ {
		addr := models.Report{
			Type: exModels.Type_Addr,
			Content: exModels.Addr{
				ID:   "qwerty",
				Name: "inkli",
			},
		}
		connc.Write(&addr)
		time.Sleep(2 * time.Second)

		test := exModels.Test{
			A: "This is A part",
			B: 156,
			C: []string{"qwer", "asdf"},
		}
		show := models.Report{
			Type: exModels.Type_Show,
			Content: exModels.Show{
				Name: "This is show",
				Te:   test,
			},
		}
		connc.Write(&show)
	}

	time.Sleep(10 * time.Second)
	cancel()
}
