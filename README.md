## Go Solidity

This repository contains a sample Golang application which interacts with a deployed solidity smart contract 

#### How to Use

- Make sure you have golang installed at your machine.
- Make sure you have a ganache installed and running.
- Update the variables at `Makefile` and execute:

```sh
make it-storage
```

The above command will execute the following actions

- Stablish a connection with the Ethereum Ganache client
- Get the current state of `SimpleStorage.sol` deployed smart contract
- Set a new value to the variable `storedValue uint` at the deployed smart contract