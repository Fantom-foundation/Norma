// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

contract Helper {

    // distribute the send in amount evenly among the given receivers
    function distribute(address payable[] calldata receivers) public payable {
        uint256 amount = msg.value / receivers.length;
        for (uint256 i = 0; i < receivers.length; i++) {
            receivers[i].transfer(amount);
        }
    }

}
