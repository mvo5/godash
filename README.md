# Unofficial golang client for the Dash kids robot

This is an unofficial (and unsanctioned) golang based client
to control the "dash" kids robot via bluetooth le.

## How to use

There is a simple API to control the robot.

E.g.:
```go

package main

import "github.com/mvo5/godash"

func main() {
	dash := godash.New()
	if err := dash.Connect(); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer dash.Disconnect()

	if err := dash.Eye(0x1b1b); err != nil {
		fmt.Println("drive err: ", err)
	}
        select {}
}
```
