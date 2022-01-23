package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/EclesioMeloJunior/go-solidity/internals/client"
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
	"GANACHE_HOST":        "",
	"CONTRACT_ADDR":       "",
	"GANHACE_NETWORK_ID":  "",
	"ACCOUNT_PRIVATE_KEY": "",
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
go run ./main.go -transfer -to 0x... -amount 10
`)

	flag.BoolVar(&listen, "listen", false, `
listen and logs events from smart contract (Minted, Sent). example:
go run ./main.go -listen
`)

	flag.BoolVar(&mint, "mint", false, `
create more token supply (only minter can execute this function). example:
go run ./main.go  -mint -amount 10
`)

	flag.StringVar(&to, "to", "", `
flag to store the public address from the receiver, used with -transfer and -mint
`)

	flag.UintVar(&amount, "amount", 0, `
flag to store the amount of tokens to be minted or transfered, used with -transfer and -mint
`)

}

func main() {
	flag.Parse()

	ethcli, err := client.NewClient(envvars["GANACHE_HOST"])
	if err != nil {
		log.Fatalf("cannot start ethereum client connection")
	}

	if listen {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		doneCh := make(chan error)
		go coin.Listen(ctx, ethcli, envvars, doneCh)

		// setup a shutdown action
		exit := make(chan struct{})
		sigCh := make(chan os.Signal, 1)

		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			defer func() {
				signal.Stop(sigCh)
				exit <- struct{}{}
			}()

			select {
			case <-sigCh:
			case err := <-doneCh:
				if err != nil {
					log.Fatalf("problems while listen: %s", err)
				}
			}
		}()

		<-exit
		return
	}

	if transafer {
		err := coin.Transfer(ethcli, envvars, to, amount)
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
