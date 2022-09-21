package main

import (
	"fmt"
	"os"
	"time"

	"github.com/adridevelopsthings/open-interlocking/pkg"
	"github.com/adridevelopsthings/open-interlocking/pkg/authorization"
	"github.com/adridevelopsthings/open-interlocking/pkg/components"
	"github.com/adridevelopsthings/open-interlocking/pkg/config"
	"github.com/akamensky/argparse"
)

func generateAuthToken(cid string, permissions []string) {
	if cid == "" {
		fmt.Printf("Set a computer-id with -c to create an auth token.\n")
		return
	}

	if len(permissions) == 0 {
		fmt.Printf("Setting up default permissions: .+ \n")
		permissions = append(permissions, ".+")
	}

	token, err := authorization.CreateToken(&cid, &permissions)
	if err != nil {
		fmt.Printf("Error while creating token: %v\n", err)
		return
	}
	fmt.Printf("Your token:\n\nBearer %s\n", token)
}

func CheckFahrstrassenEverySecond() {
	for {
		components.CheckConnections()
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	config.LoadConfiguration()
	parser := argparse.NewParser("open-interlocking", "An open source interlocking simulation")
	generateAuthTokenFlag := parser.Flag("g", "generate-auth-token", &argparse.Options{})
	permissions := parser.List("p", "generate-auth-token-permission", &argparse.Options{})
	cid := parser.String("c", "computer-id", &argparse.Options{})
	host := parser.String("l", "listen-host", &argparse.Options{Default: ":8000"})

	err := parser.Parse(os.Args)

	if *generateAuthTokenFlag {
		generateAuthToken(*cid, *permissions)
		return
	}
	err = components.ReadTemplating()

	if err != nil {
		fmt.Printf("Error while reading template config: %v\n", err)
		return
	}
	go CheckFahrstrassenEverySecond()
	pkg.RunServer(*host)
}
