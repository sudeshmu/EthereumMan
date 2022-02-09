pragma solidity ^0.4.24;

contract ERC1155 {
   function balanceOf(address tokenOwner, uint256 id) public constant returns (uint256 balance);
   event TransferSingle(address indexed operator, address indexed from, address indexed to, uint256 id, uint256 tokens);
   event TransferBatch(address indexed operator, address indexed from, address indexed to, uint256[] ids, uint256[] tokens);
}