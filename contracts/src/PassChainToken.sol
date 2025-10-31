// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/**
 * @title PassChainToken (PCT)
 * @dev ERC-20 token for Pass Chain ecosystem
 * Used for:
 * - Paying storage fees (10 PCT per credential)
 * - Paying access fees (1 PCT per retrieval)
 * - Future: Staking, governance, premium features
 */
contract PassChainToken is ERC20, Ownable, Pausable {
    uint256 private constant DECIMALS = 18;
    
    // Initial supply: 1 billion tokens
    uint256 private constant INITIAL_SUPPLY = 1_000_000_000 * 10**DECIMALS;
    
    // Fee structure
    uint256 public storageFeePerCredential = 10 * 10**DECIMALS; // 10 PCT per credential
    uint256 public accessFeePerRetrieval = 1 * 10**DECIMALS;     // 1 PCT per access
    
    // Fee recipient (treasury)
    address public treasury;
    
    // Events
    event StorageFeePaid(address indexed payer, uint256 amount, string credentialId);
    event AccessFeePaid(address indexed payer, uint256 amount, string credentialId);
    event FeeUpdated(string feeType, uint256 newAmount);
    event TreasuryUpdated(address indexed oldTreasury, address indexed newTreasury);
    
    constructor(address _treasury) ERC20("PassChain Token", "PCT") Ownable(msg.sender) {
        require(_treasury != address(0), "Invalid treasury address");
        treasury = _treasury;
        _mint(msg.sender, INITIAL_SUPPLY);
    }
    
    /**
     * @dev Pay storage fee for creating a new credential
     * @param credentialId Unique identifier for the credential
     */
    function payStorageFee(string memory credentialId) external whenNotPaused {
        require(bytes(credentialId).length > 0, "Invalid credential ID");
        require(balanceOf(msg.sender) >= storageFeePerCredential, "Insufficient balance");
        
        _transfer(msg.sender, treasury, storageFeePerCredential);
        emit StorageFeePaid(msg.sender, storageFeePerCredential, credentialId);
    }
    
    /**
     * @dev Pay access fee for retrieving a credential
     * @param credentialId Unique identifier for the credential
     */
    function payAccessFee(string memory credentialId) external whenNotPaused {
        require(bytes(credentialId).length > 0, "Invalid credential ID");
        require(balanceOf(msg.sender) >= accessFeePerRetrieval, "Insufficient balance");
        
        _transfer(msg.sender, treasury, accessFeePerRetrieval);
        emit AccessFeePaid(msg.sender, accessFeePerRetrieval, credentialId);
    }
    
    /**
     * @dev Update storage fee (only owner)
     * @param newFee New storage fee amount
     */
    function updateStorageFee(uint256 newFee) external onlyOwner {
        storageFeePerCredential = newFee;
        emit FeeUpdated("storage", newFee);
    }
    
    /**
     * @dev Update access fee (only owner)
     * @param newFee New access fee amount
     */
    function updateAccessFee(uint256 newFee) external onlyOwner {
        accessFeePerRetrieval = newFee;
        emit FeeUpdated("access", newFee);
    }
    
    /**
     * @dev Update treasury address (only owner)
     * @param newTreasury New treasury address
     */
    function updateTreasury(address newTreasury) external onlyOwner {
        require(newTreasury != address(0), "Invalid treasury address");
        address oldTreasury = treasury;
        treasury = newTreasury;
        emit TreasuryUpdated(oldTreasury, newTreasury);
    }
    
    /**
     * @dev Mint new tokens (only owner)
     * @param to Address to receive tokens
     * @param amount Amount to mint
     */
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }
    
    /**
     * @dev Burn tokens from caller's balance
     * @param amount Amount to burn
     */
    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }
    
    /**
     * @dev Pause token transfers (only owner, emergency use)
     */
    function pause() external onlyOwner {
        _pause();
    }
    
    /**
     * @dev Unpause token transfers (only owner)
     */
    function unpause() external onlyOwner {
        _unpause();
    }
    
    /**
     * @dev Override _update to add pause functionality
     */
    function _update(address from, address to, uint256 value)
        internal
        override
        whenNotPaused
    {
        super._update(from, to, value);
    }
}
