pragma solidity >=0.4.22 <0.9.0;

contract Coin {
    address public minter;
    uint public totalSupply;

    mapping(address => uint) public balances;

    event Sent(address from, address to, uint amount);
    event Minted(address owner, uint amount);

    constructor() {
        minter = msg.sender;
    }

    function mint(address receiver, uint amount) public {
        require(msg.sender == minter);
        balances[receiver] += amount;
        totalSupply += amount;

        emit Minted(receiver, amount);
    }

    error InsufficientBalance(uint requested, uint available);

    function send(address receiver, uint amount) public {
        if (amount > balances[msg.sender]) {
            revert InsufficientBalance({
                requested: amount,
                available: balances[msg.sender]
            });
        }

        balances[msg.sender] -= amount;
        balances[receiver] += amount;
        emit Sent(msg.sender, receiver, amount);
    }
}