package main

import (
	"fmt"
	"time"

	"github.com/adridevelopsthings/open-interlocking/pkg"
	"github.com/adridevelopsthings/open-interlocking/pkg/components"
)

func CheckFahrstrassenEverySecond() {
	for {
		components.CheckConnections()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	err := components.ReadTemplating()
	if err != nil {
		fmt.Printf("Error while reading template config: %v\n", err)
		return
	}
	go CheckFahrstrassenEverySecond()
	pkg.RunServer()
}
