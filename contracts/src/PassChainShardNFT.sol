// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title PassChainShardNFT (PCS)
 * @dev ERC-721 NFT representing encrypted Share3 backup on IPFS
 * 
 * Each NFT = One encrypted credential shard backup
 * - Permanent decentralized backup
 * - Transferable (can gift credentials)
 * - Tradeable (secondary market)
 * - IPFS-based metadata
 */
contract PassChainShardNFT is ERC721URIStorage, Ownable {
    // Auto-incrementing token ID
    uint256 private _nextTokenId = 1;
    
    // Authorized minter (backend service)
    address public minter;
    
    // Mapping from credential ID to token ID
    mapping(string => uint256) public credentialToToken;
    
    // Mapping from token ID to credential ID
    mapping(uint256 => string) public tokenToCredential;
    
    // Mapping from wallet to their shard token IDs
    mapping(address => uint256[]) private walletShards;
    
    // Events
    event ShardNFTMinted(address indexed owner, uint256 indexed tokenId, string credentialId, string ipfsCID);
    event ShardBurned(address indexed owner, uint256 indexed tokenId, string credentialId);
    event MinterUpdated(address indexed oldMinter, address indexed newMinter);
    
    // Modifiers
    modifier onlyMinter() {
        require(msg.sender == minter, "Only minter can call this");
        _;
    }
    
    constructor(address _minter) ERC721("PassChain Shard NFT", "PCS") Ownable(msg.sender) {
        require(_minter != address(0), "Invalid minter address");
        minter = _minter;
    }
    
    /**
     * @dev Mint a new shard NFT for a credential backup
     * @param to Wallet address to receive the NFT
     * @param credentialId Unique credential identifier
     * @param ipfsCID IPFS hash of the encrypted shard
     */
    function mintShard(
        address to,
        string memory credentialId,
        string memory ipfsCID
    ) public onlyMinter returns (uint256) {
        require(credentialToToken[credentialId] == 0, "Shard already exists");
        
        uint256 newTokenId = _nextTokenId++;
        _safeMint(to, newTokenId);
        _setTokenURI(newTokenId, string(abi.encodePacked("ipfs://", ipfsCID)));
        
        // Map credential ID to token ID
        credentialToToken[credentialId] = newTokenId;
        tokenToCredential[newTokenId] = credentialId;
        
        // Track user's shards
        walletShards[to].push(newTokenId);
        
        emit ShardNFTMinted(to, newTokenId, credentialId, ipfsCID);
        
        return newTokenId;
    }
    
    /**
     * @dev Burn (delete) a shard NFT
     * @param tokenId Token ID to burn
     */
    function burnShard(uint256 tokenId) public {
        require(ownerOf(tokenId) == msg.sender, "Not the owner");
        
        string memory credentialId = tokenToCredential[tokenId];
        
        // Remove from mappings
        delete credentialToToken[credentialId];
        delete tokenToCredential[tokenId];
        
        // Remove from wallet shards
        _removeFromWalletShards(msg.sender, tokenId);
        
        // Burn the NFT
        _burn(tokenId);
        
        emit ShardBurned(msg.sender, tokenId, credentialId);
    }
    
    /**
     * @dev Get all shard token IDs owned by a wallet
     * @param wallet Wallet address to query
     * @return Array of token IDs
     */
    function getWalletShards(address wallet) public view returns (uint256[] memory) {
        return walletShards[wallet];
    }
    
    /**
     * @dev Get token ID for a credential ID
     * @param credentialId Credential identifier
     * @return Token ID (0 if not found)
     */
    function getTokenId(string memory credentialId) public view returns (uint256) {
        return credentialToToken[credentialId];
    }
    
    /**
     * @dev Get credential ID for a token ID
     * @param tokenId Token identifier
     * @return Credential ID
     */
    function getCredentialId(uint256 tokenId) public view returns (string memory) {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        return tokenToCredential[tokenId];
    }
    
    /**
     * @dev Update the authorized minter address
     * @param newMinter New minter address
     */
    function updateMinter(address newMinter) public onlyOwner {
        require(newMinter != address(0), "Invalid minter address");
        address oldMinter = minter;
        minter = newMinter;
        emit MinterUpdated(oldMinter, newMinter);
    }
    
    /**
     * @dev Remove a token ID from wallet's shard array
     */
    function _removeFromWalletShards(address wallet, uint256 tokenId) private {
        uint256[] storage shards = walletShards[wallet];
        for (uint256 i = 0; i < shards.length; i++) {
            if (shards[i] == tokenId) {
                shards[i] = shards[shards.length - 1];
                shards.pop();
                break;
            }
        }
    }
    
    /**
     * @dev Override transfer to update wallet shards tracking
     */
    function _update(address to, uint256 tokenId, address auth)
        internal
        override
        returns (address)
    {
        address from = _ownerOf(tokenId);
        address result = super._update(to, tokenId, auth);
        
        // Update wallet shards on transfer (but not on mint)
        if (from != address(0) && to != address(0) && from != to) {
            _removeFromWalletShards(from, tokenId);
            walletShards[to].push(tokenId);
        }
        
        return result;
    }
    
    /**
     * @dev Override required for ERC721URIStorage
     */
    function tokenURI(uint256 tokenId)
        public
        view
        override(ERC721URIStorage)
        returns (string memory)
    {
        return super.tokenURI(tokenId);
    }
    
    /**
     * @dev Override required for ERC721URIStorage
     */
    function supportsInterface(bytes4 interfaceId)
        public
        view
        override(ERC721URIStorage)
        returns (bool)
    {
        return super.supportsInterface(interfaceId);
    }
}
