package coin

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/EclesioMeloJunior/go-solidity/internals/client"
	ethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

func Mint(ethcli *client.ETH, envvars map[string]string, to string, amount uint) error {
	chainID, err := strconv.ParseInt(envvars["GANHACE_NETWORK_ID"], 10, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %s to chain id: %w",
			envvars["GANHACE_NETWORK_ID"], err)
	}

	if strings.TrimSpace(to) == "" {
		return errors.New("-to flag must have a value")
	}

	coinContractAddr := ethcommon.HexToAddress(envvars["CONTRACT_ADDR"])
	coinContractInst, err := NewCoin(coinContractAddr, ethcli.Client)
	if err != nil {
		return fmt.Errorf("cannot create coin contract inst: %w", err)
	}

	privkey, pubkey, err := client.FromPrivateKey(envvars["ACCOUNT_PRIVATE_KEY"])
	if err != nil {
		return fmt.Errorf("cannot decode private key: %w", err)
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

	auth, err := ethbind.NewKeyedTransactorWithChainID(privkey, big.NewInt(chainID))
	if err != nil {
		return fmt.Errorf("cannot create a transactor: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonceToUse))
	auth.Value = big.NewInt(0)
	auth.GasLimit = GasLimit
	auth.GasPrice = gasPrice

	toAddr := ethcommon.HexToAddress(to)
	amountToSend := big.NewInt(int64(amount))

	messageToAsk := fmt.Sprintf("Mint: %d\nTo: %s\nGas Price: %d\n",
		amountToSend, toAddr, gasPrice)
	messageToAsk += fmt.Sprintf("Are you sure? (yes/no): ")

	awnser, err := askBeforeGo(messageToAsk)
	if err != nil {
		return fmt.Errorf("problems while asking: %w", err)
	}

	if !awnser {
		fmt.Println("transaction canceled!")
		return nil
	}

	// TODO: improve the error handling to output the tx failure reason!
	tx, err := coinContractInst.Mint(auth, toAddr, amountToSend)
	if err != nil {
		return fmt.Errorf("failed to send tokens: %w", err)
	}

	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	return nil
}
