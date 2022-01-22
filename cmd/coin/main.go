package main

import (
	"log"
	"os"
)

var envvars = map[string]string{
	"GANACHE_HOST":        "",
	"CONTRACT_ADDR":       "",
	"ACCOUNT_PRIVATE_KEY": "",
	"GANHACE_NETWORK_ID":  "",
}

func evalEnvVars() {
	for key := range envvars {
		value, exists := os.LookupEnv(key)
		if !exists {
			log.Fatalf("variable %s must exists", key)
		}
		envvars[key] = value
	}
}
