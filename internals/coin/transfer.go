package coin

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/EclesioMeloJunior/go-solidity/internals/client"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	GasLimit  uint64 = 300000
	YesAwnser        = "yes"
)

func Transfer(ethcli *client.ETH, envvars map[string]string, to string, amount uint) error {
	chainID, err := strconv.ParseInt(envvars["GANHACE_NETWORK_ID"], 10, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %s to chain id: %w",
			envvars["GANHACE_NETWORK_ID"], err)
	}

	if strings.TrimSpace(to) == "" {
		return errors.New("-to flag must have a value")
	}

	privkey, pubkey, err := client.FromPrivateKey(envvars["ACCOUNT_PRIVATE_KEY"])
	if err != nil {
		return fmt.Errorf("cannot decode private key: %w", err)
	}

	coinContractAddr := ethcommon.HexToAddress(envvars["CONTRACT_ADDR"])
	coinContractInst, err := NewCoin(coinContractAddr, ethcli.Client)
	if err != nil {
		return fmt.Errorf("cannot create coin contract inst: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nonceToUse, err := ethcli.Client.PendingNonceAt(ctx, *pubkey)
	if err != nil {
		return fmt.Errorf("cannot get pending nonce at %s: %w", pubkey, err)
	}

	gasPrice, err := ethcli.Client.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("cannot get suggested gas price: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privkey, big.NewInt(chainID))
	if err != nil {
		return fmt.Errorf("cannot create a transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonceToUse))
	auth.Value = big.NewInt(0)
	auth.GasLimit = GasLimit
	auth.GasPrice = gasPrice

	toAddr := ethcommon.HexToAddress(to)
	amountToSend := big.NewInt(int64(amount))

	awnser, err := askBeforeGo(*pubkey, toAddr, amountToSend, gasPrice)
	if err != nil {
		return fmt.Errorf("problems while asking: %w", err)
	}

	if !awnser {
		fmt.Println("transaction canceled!")
		return nil
	}

	tx, err := coinContractInst.Send(auth, toAddr, amountToSend)
	if err != nil {
		return fmt.Errorf("failed to send tokens: %w", err)
	}

	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	return nil
}

func askBeforeGo(from, to ethcommon.Address, value, gasPrice *big.Int) (aws bool, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("From: %s\nTo: %s\nValue: %d\nGas Price: %d\n",
		from, to, value, gasPrice)
	fmt.Printf("Are you sure? (yes/no): ")

	line, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	line = strings.TrimRight(line, "\r\n")
	return line == YesAwnser, nil
}
