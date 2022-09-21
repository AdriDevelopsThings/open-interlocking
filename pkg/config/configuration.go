package config

import "os"

var JWTSecret string = ""
var JWTVersion string = "1"

func LoadConfiguration() {
	JWTSecret = os.Getenv("AUTH_JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "development-1"
	}
}
