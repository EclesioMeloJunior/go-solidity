## Go Solidity

This repository contains a sample Golang application which interacts with a deployed solidity smart contract 

### First steps

- Install Ganache here https://trufflesuite.com/ganache/
- Run your local ganache and then update the `Makefile` file with the right values of your running ganache instance
- Once you have a running Ganache locally, in the gear icon at `Truffle Projects` you should add the `truffle-config.js` file to enable a better visualization of what's going on with the deployed contract (eg. get the contract address).
- Install truffle with npm `npm install -g truffle`
- At project root execute `truffle migrate`, you should notice new blocks created at your running ganache instance


#### simple-storage (example)

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

#### coin (example)


In this example, I built a simple smart contract which keeps a mapping of `account -> amount of tokens (uint)`, and the smart contract has 2 public functions `sent` and `mint`, where the `mint` function is only called by the smart contract upload account (minter)

(... same requirements from past example)

In one terminal execute:

> obs: the `make it-coin-listen` at Makefile should contains `ws://` instead of `http://` at `GANACHE_HOST` env var
```
make it-coin-listen
```

This will spins up a filter logs subscriber which will output 2 kinds of events: `Sent` and `Minted`, these two events are emmited when a token are sent from an account to another and when the minter creates more tokens (increase the supply), respectively.

You can use the `truffle develop` command and interact with the deployed smart contract and call the `mint` and `sent` function to emmit the events and check the output of the listener!

![logs output](https://github.com/EclesioMeloJunior/go-solidity/blob/main/assets/listen-logs.png?raw=true)

obs: I'm currently working on the `make it-coin-transfer` and `make it-coin-mint` commands!

