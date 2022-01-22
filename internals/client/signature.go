package client

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

func FromPrivateKey(enc string) (privkey *ecdsa.PrivateKey, pubaddr *ethcommon.Address, err error) {
	privateKey, err := ethcrypto.HexToECDSA(enc)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot parse private key to ECDSA: %w", err)
	}

	pubkey := privateKey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("expected pubkey was type *ecdsa.PublicKey")
	}

	addr := ethcrypto.PubkeyToAddress(*pubkeyECDSA)
	return privateKey, &addr, nil
}
