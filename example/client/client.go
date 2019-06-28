package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/PTFS/connSocket/connC"
	"github.com/PTFS/connSocket/models"
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
				fmt.Println("Receive Msg", cmd.Content)
			case <-ctx.Done():
				return
			}
		}
	}()
	//Write
	for i := 0; i < 20; i++ {
		re := models.Login{
			ID: "qwerty",
		}
		data, err := re.String()
		if err != nil {
			continue
		}
		report := models.Report{
			Type:    models.Type_Login,
			Content: string(data),
		}
		connc.Write(&report)
		time.Sleep(2 * time.Second)
	}
	cancel()
}
