package main

import (
	"fmt"
	"log"

	"github.com/mvo5/godash"
)

func main() {
	dash := godash.New()
	if err := dash.Connect(); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer dash.Disconnect()

	if err := dash.Eye(0x1b1b); err != nil {
		fmt.Println("eye err: ", err)
	}
	driving := false
	for {
		switch {
		case dash.ProxRear > 16:
			driving = true
			if err := dash.Drive(int(dash.ProxRear) * 3); err != nil {
				log.Fatalf("drive: %s", err)
			}
		default:
			if driving {
				driving = false
				dash.Stop()
			}
		}
	}
}
