// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

contract Counter {
    int private count = 1;  // starts with 1 to make all increments cost the same

    function incrementCounter() public {
        count += 1;
    }

    function getCount() public view returns (int) {
        return count-1;
    }
}
