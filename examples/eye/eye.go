package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/mvo5/godash"
)

func main() {
	dash := godash.New()
	if err := dash.Connect(); err != nil {
		log.Fatalf("cannot connect: %v", err)
	}
	defer dash.Disconnect()

	if err := dash.Eye(2546); err != nil {
		fmt.Println("drive err: ", err)
	}

	print("press enter when finished")
	r := bufio.NewReader(os.Stdin)
	r.ReadLine()
}
