// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract CrossChainBridge {
    event TokensLocked(address indexed sender, uint256 amount, string hyperledgerAddress);

    mapping(address => uint256) public lockedBalances;

    function lockTokens(string memory hyperledgerAddress) public payable {
        require(msg.value > 0, "Must send ETH");
        lockedBalances[msg.sender] += msg.value;

        emit TokensLocked(msg.sender, msg.value, hyperledgerAddress);
    }

    function getLockedBalance(address user) public view returns (uint256) {
        return lockedBalances[user];
    }
}
