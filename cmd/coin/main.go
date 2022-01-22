package main

import (
	"flag"
	"log"
	"os"

	"github.com/EclesioMeloJunior/go-solidity/internals/coin"
)

// transfer flags

var (
	transafer bool
	listen    bool
	mint      bool
	to        string
	amount    uint
)

var envvars = map[string]string{
	"GANACHE_HOST":       "",
	"CONTRACT_ADDR":      "",
	"GANHACE_NETWORK_ID": "",
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

func init() {
	evalEnvVars()

	flag.BoolVar(&transafer, "transfer", false, `
execute a coin transfer from sender to another account. example:
make it-coin -transfer -to 0x... -amount 10
`)

	flag.BoolVar(&transafer, "listen", false, `
listen and logs events from smart contract (Minted, Sent). example:
make it-coin -listen
`)

	flag.BoolVar(&transafer, "mint", false, `
create more token supply (only minter can execute this function). example:
make it-coin -mint -amount 10
`)
}

func main() {
	flag.Parse()

	if listen {
		err := coin.Listen()
		if err != nil {
			log.Fatalf("cannot listen: %s", err)
		}
	}

	if transafer {
		err := coin.Transfer(to, amount)
		if err != nil {
			log.Fatalf("cannot transfer: %s", err)
		}
	}

	if mint {
		err := coin.Mint(amount)
		if err != nil {
			log.Fatalf("cannot mint: %s", err)
		}
	}
}
