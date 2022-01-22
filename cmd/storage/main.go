package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"

	"github.com/EclesioMeloJunior/go-solidity/internals/client"
	"github.com/EclesioMeloJunior/go-solidity/internals/simple_storage"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
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

func init() {
	evalEnvVars()
}

func main() {
	ethcli, err := client.NewClient(envvars["GANACHE_HOST"])
	if err != nil {
		log.Fatalf("cannot start eth client: %s", err)
	}

	addr := ethcommon.HexToAddress(envvars["CONTRACT_ADDR"])
	simpleStorageInst, err := simple_storage.NewSimpleStorage(addr, ethcli.Client)
	if err != nil {
		log.Fatalf("cannot create a simple_storage instance: %s", err)
	}

	storedValue, err := simpleStorageInst.Get(nil)
	if err != nil {
		log.Fatalf("cannot call simple_storage Get function: %s", err)
	}

	fmt.Printf("simple_storage current value ==> %d\n", storedValue.Uint64())

	// writing operations need a signature
	privkey, pubaddr, err := client.FromPrivateKey(envvars["ACCOUNT_PRIVATE_KEY"])
	if err != nil {
		log.Fatalf("cannot generate signature address: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nonceToUse, err := ethcli.Client.PendingNonceAt(ctx, *pubaddr)
	if err != nil {
		log.Fatalf("cannot retrieve the pending nonce: %s", err)
	}

	gasPrice, err := ethcli.Client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("cannot retrieve the suggested fas price: %s", err)
	}

	chainId, err := strconv.ParseInt(envvars["GANHACE_NETWORK_ID"], 10, 64)
	if err != nil {
		log.Fatalf("cannot convert ganache network id %s to int: %s",
			envvars["GANHACE_NETWORK_ID"], err)
	}

	auth, err := ethbind.NewKeyedTransactorWithChainID(privkey, big.NewInt(chainId))
	if err != nil {
		log.Fatalf("cannot create a keyed transactor: %s", err)
	}

	auth.Nonce = big.NewInt(int64(nonceToUse))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	// increment the current stored value by 1
	storedValue = big.NewInt(0).Add(storedValue, big.NewInt(1))

	tx, err := simpleStorageInst.Set(auth, storedValue)
	if err != nil {
		log.Fatalf("cannot call simple_storage set with value %d: %s", storedValue, err)
	}
	fmt.Printf("simple_storage set called\ntransaction_id: %s\n", tx.Hash().Hex())
}
