package main

import (
	"log"

	"github.com/mvo5/godash"
)

const SENSITIVITY = 16

func main() {
	dash := godash.New()
	if err := dash.Connect(); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer dash.Disconnect()

	driving := false
	for {
		switch {
		case dash.ProxLeft > SENSITIVITY:
			driving = false
			if err := dash.Stop(); err != nil {
				log.Fatalf("stop: %s", err)
			}
			if err := dash.Turn(90); err != nil {
				log.Fatalf("turn: %s", err)
			}
		default:
			if !driving {
				driving = true
				if err := dash.Drive(50); err != nil {
					log.Fatalf("drive: %s", err)
				}
			}
		}

	}
}
