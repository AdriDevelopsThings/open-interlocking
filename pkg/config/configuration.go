package config

import "os"

var JWTSecret string = ""
var JWTVersion string = "1"
var IgnoreAcknowledgements = false

func LoadConfiguration() {
	JWTSecret = os.Getenv("AUTH_JWT_SECRET")
	if JWTSecret == "" {
		JWTSecret = "development-1"
	}
	IGNORE_ACKNOWLEDGEMENTS := os.Getenv("IGNORE_ACKNOWLEDGEMENTS")
	if IGNORE_ACKNOWLEDGEMENTS == "true" {
		IgnoreAcknowledgements = true
	}
}
