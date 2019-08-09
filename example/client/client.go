package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LiMoMoMo/go-connSocket/connC"
	exModels "github.com/LiMoMoMo/go-connSocket/example/models"
	"github.com/LiMoMoMo/go-connSocket/models"
)

func main() {
	connc, err := connC.NewConnC("8976", "127.0.0.1")
	if err != nil {
		log.Println(err)
		return
	}
	connc.SetReconnect(func() {
		re := models.Register{
			ID: "qwerty",
		}
		report := models.Report{
			Type:    models.Type_Register,
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
		re := models.Register{
			ID: "qwerty",
		}
		report := models.Report{
			Type:    models.Type_Register,
			Content: re,
		}

		connc.Write(&report)
		time.Sleep(2 * time.Second)
	}
	for i := 0; i < 20; i++ {
		pTrafic := []exModels.TrafficInfo{}
		cTrafic := []exModels.ChildTraffic{}
		pTrafic = append(pTrafic, exModels.TrafficInfo{"qwerty", 60, 50})
		pTrafic = append(pTrafic, exModels.TrafficInfo{"asdfgh", 78, 35})
		cTrafic = append(cTrafic, pTrafic)
		addr := models.Report{
			Type: exModels.Type_HeartBeat,
			Content: exModels.HeatbeatInfo{
				ID:            "qwerty",
				ParentTraffic: pTrafic,
				ChildTraffics: cTrafic,
			},
		}
		connc.Write(&addr)
		time.Sleep(2 * time.Second)

		// test := exModels.Test{
		// 	A: "This is A part",
		// 	B: 156,
		// 	C: []string{"qwer", "asdf"},
		// }
		// show := models.Report{
		// 	Type: exModels.Type_Show,
		// 	Content: exModels.Show{
		// 		Name: "This is show",
		// 		Te:   test,
		// 	},
		// }
		// connc.Write(&show)
	}

	time.Sleep(10 * time.Second)
	cancel()
}
