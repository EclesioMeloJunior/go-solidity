package client

import "github.com/ethereum/go-ethereum/ethclient"

type ETH struct {
	Address string
	Client  *ethclient.Client
}

func NewClient(addr string) (ethcli *ETH, err error) {
	// TODO: use a explicit context
	c, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}

	return &ETH{
		Address: addr,
		Client:  c,
	}, nil
}
