// SPDX-License-Identifier: GPL-3.0-only
pragma solidity ^0.8.4;

interface IUniswapV2Pair {
    function swap(uint amount0Out, uint amount1Out, address to, bytes calldata data) external;
    function getReserves() external view returns (uint112 reserve0, uint112 reserve1, uint32 blockTimestampLast);
}

// a library for performing overflow-safe math, courtesy of DappHub (https://github.com/dapphub/ds-math)
library SafeMath {
    function add(uint x, uint y) internal pure returns (uint z) {
        require((z = x + y) >= x, 'ds-math-add-overflow');
    }

    function sub(uint x, uint y) internal pure returns (uint z) {
        require((z = x - y) <= x, 'ds-math-sub-underflow');
    }

    function mul(uint x, uint y) internal pure returns (uint z) {
        require(y == 0 || (z = x * y) / y == x, 'ds-math-mul-overflow');
    }
}

// Based on https://github.com/Uniswap/v2-periphery/blob/master/contracts/UniswapV2Router02.sol
contract UniswapRouter {
    using SafeMath for uint;
    int private count = 0;

    // **** SWAP ****
    // requires the initial amount to have already been sent to the first pair
    function _swap(uint[] memory amounts, address[] memory path, address[] memory pairsPath, address _to) internal virtual {
        for (uint i; i < pairsPath.length; i++) {
            (address input, address output) = (path[i], path[i + 1]);
            uint amountOut = amounts[i + 1];
            (uint amount0Out, uint amount1Out) = input < output ? (uint(0), amountOut) : (amountOut, uint(0));
            address to = i+1 < pairsPath.length ? pairsPath[i+1] : _to;
            IUniswapV2Pair(pairsPath[i]).swap(
                amount0Out, amount1Out, to, new bytes(0)
            );
        }
    }

    // simplified version of swapping function from Uniswap router - does not require amountOutMin,
    // swapped tokens are always sent to the tx sender
    function swapExactTokensForTokens(
        uint amountIn,
        address[] calldata path,
        address[] calldata pairsPath
    ) external returns (uint[] memory amounts) {
        amounts = getAmountsOut(amountIn, path, pairsPath);
        safeTransferFrom(
            path[0], msg.sender, pairsPath[0], amounts[0]
        );
        _swap(amounts, path, pairsPath, msg.sender);
        count++;
    }

    // performs chained getAmountOut calculations on any number of pairs
    function getAmountsOut(uint amountIn, address[] memory path, address[] memory pairsPath) internal view returns (uint[] memory amounts) {
        require(path.length >= 2, 'UniswapV2Library: INVALID_PATH');
        require(path.length == pairsPath.length+1, 'invalid length of pairsPath param');
        amounts = new uint[](path.length);
        amounts[0] = amountIn;
        for (uint i; i < path.length - 1; i++) {
            (uint reserveIn, uint reserveOut) = getReserves(path[i], path[i + 1], pairsPath[i]);
            amounts[i + 1] = getAmountOut(amounts[i], reserveIn, reserveOut);
        }
    }

    // given an input amount of an asset and pair reserves, returns the maximum output amount of the other asset
    function getAmountOut(uint amountIn, uint reserveIn, uint reserveOut) internal pure returns (uint amountOut) {
        require(amountIn > 0, 'UniswapV2Library: INSUFFICIENT_INPUT_AMOUNT');
        require(reserveIn > 0 && reserveOut > 0, 'UniswapV2Library: INSUFFICIENT_LIQUIDITY');
        uint amountInWithFee = amountIn.mul(997);
        uint numerator = amountInWithFee.mul(reserveOut);
        uint denominator = reserveIn.mul(1000).add(amountInWithFee);
        amountOut = numerator / denominator;
    }

    // fetches and sorts the reserves for a pair
    function getReserves(address tokenA, address tokenB, address pair) internal view returns (uint reserveA, uint reserveB) {
        (uint reserve0, uint reserve1,) = IUniswapV2Pair(pair).getReserves();
        (reserveA, reserveB) = tokenA < tokenB ? (reserve0, reserve1) : (reserve1, reserve0);
    }

    // helper method for interacting with ERC20 tokens that do not consistently return true/false
    function safeTransferFrom(address token, address from, address to, uint value) internal {
        // bytes4(keccak256(bytes('transferFrom(address,address,uint256)')));
        (bool success, bytes memory data) = token.call(abi.encodeWithSelector(0x23b872dd, from, to, value));
        require(success && (data.length == 0 || abi.decode(data, (bool))), 'TransferHelper: TRANSFER_FROM_FAILED');
    }

    function getCount() public view returns (int) {
        return count;
    }
}
