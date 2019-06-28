package main

import (
	connS "github.com/PTFS/connSocket/connS"
)

func main() {
	s := connS.NewConnS("8421", "0.0.0.0")
	s.Start()
	select {}
}
