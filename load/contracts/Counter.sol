// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

contract Counter {
    // counter is an internal counter tracking the number of increment calls.
    // The counter is initialized to 1 to make all increment-counter calls
    // equally expensive. Otherwise, the first call incrementing the counter
    // from 0 to 1 would have to pay extra gas for the storage allocation.
    int private count = 1;

    function incrementCounter() public {
        count += 1;
    }

    function getCount() public view returns (int) {
        return count-1;
    }
}
