package coin

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/EclesioMeloJunior/go-solidity/internals/client"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	EventSent   = []byte("Sent(address,address,uint256)")
	EventMinted = []byte("Minted(address,uint256)")

	logEventSentSignatureHash   = crypto.Keccak256Hash(EventSent)
	logEventMintedSignatureHash = crypto.Keccak256Hash(EventMinted)
)

type EventSentData struct {
	From   ethcommon.Address
	To     ethcommon.Address
	Amount *big.Int
}

type EventMintedData struct {
	Owner  ethcommon.Address
	Amount *big.Int
}

func Listen(ctx context.Context, ethcli *client.ETH, envvars map[string]string, done chan<- error) {
	coinContractAddr := ethcommon.HexToAddress(envvars["CONTRACT_ADDR"])
	query := ethereum.FilterQuery{
		Addresses: []ethcommon.Address{
			coinContractAddr,
		},
	}

	logCh := make(chan ethtypes.Log)
	subscriber, err := ethcli.Client.SubscribeFilterLogs(ctx, query, logCh)
	if err != nil {
		done <- fmt.Errorf("cannot subscribe to contract logs: %w", err)
		return
	}

	coinContractABI, err := abi.JSON(strings.NewReader(string(CoinABI)))
	if err != nil {
		done <- fmt.Errorf("problems to parse contract ABI: %w", err)
		return
	}

	fmt.Println("-> start listen <-")
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("finish listener: %s\n", ctx.Err())
			return
		case err := <-subscriber.Err():
			subscriber.Unsubscribe()
			done <- fmt.Errorf("problems with subscriber: %w", err)
			return
		case log := <-logCh:
			output := "Event: %s\n==> Sender: %s\n"

			switch log.Topics[0].Hex() {
			case logEventSentSignatureHash.Hex():
				var event EventSentData
				err := coinContractABI.UnpackIntoInterface(&event, "Sent", log.Data)
				if err != nil {
					done <- fmt.Errorf("cannot unpack event Sent into struct: %w", err)
				}

				output += "==> From: %s\n==> To: %s\n==> Amount: %d\n"

				fmt.Printf(output,
					"Sent",
					log.Address,
					event.From,
					event.To,
					event.Amount,
				)
			case logEventMintedSignatureHash.Hex():
				var event EventMintedData
				err := coinContractABI.UnpackIntoInterface(&event, "Minted", log.Data)
				if err != nil {
					done <- fmt.Errorf("cannot unpack event Sent into struct: %w", err)
				}

				output += "==> Owner: %s\n==> Amount: %d\n"

				fmt.Printf(output,
					"Minted",
					log.Address,
					event.Owner,
					event.Amount,
				)
			}
		}
	}
}
