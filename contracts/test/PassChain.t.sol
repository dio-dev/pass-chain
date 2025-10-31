// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Test.sol";
import "../src/PassChainToken.sol";
import "../src/PassChainShardNFT.sol";

contract PassChainTest is Test {
    PassChainToken public token;
    PassChainShardNFT public nft;
    
    address public owner = address(this);
    address public treasury = address(0x1);
    address public backend = address(0x2);
    address public user1 = address(0x3);
    address public user2 = address(0x4);
    
    uint256 constant INITIAL_BALANCE = 1000 * 10**18; // 1000 PCT
    
    function setUp() public {
        // Deploy contracts
        token = new PassChainToken(treasury);
        nft = new PassChainShardNFT(backend);
        
        // Give users some tokens
        token.transfer(user1, INITIAL_BALANCE);
        token.transfer(user2, INITIAL_BALANCE);
    }
    
    // ========================================
    // Token Tests
    // ========================================
    
    function testTokenInitialization() public {
        assertEq(token.name(), "PassChain Token");
        assertEq(token.symbol(), "PCT");
        assertEq(token.treasury(), treasury);
        assertEq(token.storageFeePerCredential(), 10 * 10**18);
        assertEq(token.accessFeePerRetrieval(), 1 * 10**18);
    }
    
    function testPayStorageFee() public {
        vm.startPrank(user1);
        
        uint256 initialBalance = token.balanceOf(user1);
        uint256 initialTreasury = token.balanceOf(treasury);
        
        token.payStorageFee("credential-123");
        
        assertEq(token.balanceOf(user1), initialBalance - 10 * 10**18);
        assertEq(token.balanceOf(treasury), initialTreasury + 10 * 10**18);
        
        vm.stopPrank();
    }
    
    function testPayAccessFee() public {
        vm.startPrank(user1);
        
        uint256 initialBalance = token.balanceOf(user1);
        uint256 initialTreasury = token.balanceOf(treasury);
        
        token.payAccessFee("credential-123");
        
        assertEq(token.balanceOf(user1), initialBalance - 1 * 10**18);
        assertEq(token.balanceOf(treasury), initialTreasury + 1 * 10**18);
        
        vm.stopPrank();
    }
    
    function testCannotPayWithInsufficientBalance() public {
        address poorUser = address(0x5);
        
        vm.startPrank(poorUser);
        vm.expectRevert("Insufficient balance");
        token.payStorageFee("credential-123");
        vm.stopPrank();
    }
    
    function testUpdateFees() public {
        token.updateStorageFee(20 * 10**18);
        assertEq(token.storageFeePerCredential(), 20 * 10**18);
        
        token.updateAccessFee(2 * 10**18);
        assertEq(token.accessFeePerRetrieval(), 2 * 10**18);
    }
    
    function testUpdateTreasury() public {
        address newTreasury = address(0x10);
        token.updateTreasury(newTreasury);
        assertEq(token.treasury(), newTreasury);
    }
    
    function testPauseUnpause() public {
        token.pause();
        
        vm.startPrank(user1);
        vm.expectRevert("Pausable: paused");
        token.transfer(user2, 100);
        vm.stopPrank();
        
        token.unpause();
        
        vm.startPrank(user1);
        token.transfer(user2, 100);
        vm.stopPrank();
    }
    
    // ========================================
    // NFT Tests
    // ========================================
    
    function testNFTInitialization() public {
        assertEq(nft.name(), "PassChain Shard NFT");
        assertEq(nft.symbol(), "PCS");
        assertEq(nft.minter(), backend);
    }
    
    function testMintShard() public {
        vm.startPrank(backend);
        
        uint256 tokenId = nft.mintShard(
            user1,
            "credential-123",
            "QmXyz123abc"
        );
        
        assertEq(tokenId, 1);
        assertEq(nft.ownerOf(tokenId), user1);
        assertEq(nft.getCredentialId(tokenId), "credential-123");
        assertEq(nft.getTokenId("credential-123"), tokenId);
        
        string memory uri = nft.tokenURI(tokenId);
        assertEq(uri, "ipfs://QmXyz123abc");
        
        vm.stopPrank();
    }
    
    function testCannotMintDuplicate() public {
        vm.startPrank(backend);
        
        nft.mintShard(user1, "credential-123", "QmXyz123");
        
        vm.expectRevert("Shard already exists");
        nft.mintShard(user1, "credential-123", "QmXyz456");
        
        vm.stopPrank();
    }
    
    function testOnlyMinterCanMint() public {
        vm.startPrank(user1);
        
        vm.expectRevert("Only minter can call this");
        nft.mintShard(user1, "credential-123", "QmXyz123");
        
        vm.stopPrank();
    }
    
    function testBurnShard() public {
        // Mint first
        vm.prank(backend);
        uint256 tokenId = nft.mintShard(user1, "credential-123", "QmXyz123");
        
        // Burn
        vm.prank(user1);
        nft.burnShard(tokenId);
        
        // Check it's gone
        vm.expectRevert();
        nft.ownerOf(tokenId);
    }
    
    function testGetWalletShards() public {
        vm.startPrank(backend);
        
        nft.mintShard(user1, "cred-1", "QmHash1");
        nft.mintShard(user1, "cred-2", "QmHash2");
        nft.mintShard(user1, "cred-3", "QmHash3");
        
        vm.stopPrank();
        
        uint256[] memory shards = nft.getWalletShards(user1);
        assertEq(shards.length, 3);
    }
    
    function testTransferNFT() public {
        // Mint
        vm.prank(backend);
        uint256 tokenId = nft.mintShard(user1, "credential-123", "QmXyz123");
        
        // Transfer
        vm.prank(user1);
        nft.transferFrom(user1, user2, tokenId);
        
        // Check new owner
        assertEq(nft.ownerOf(tokenId), user2);
        
        // Check wallet shards updated
        assertEq(nft.getWalletShards(user1).length, 0);
        assertEq(nft.getWalletShards(user2).length, 1);
    }
    
    // ========================================
    // Integration Tests
    // ========================================
    
    function testCompleteFlow() public {
        // 1. User pays storage fee
        vm.prank(user1);
        token.payStorageFee("credential-123");
        
        // 2. Backend mints NFT with shard
        vm.prank(backend);
        uint256 tokenId = nft.mintShard(user1, "credential-123", "QmEncryptedShard");
        
        // 3. User accesses credential (pays fee)
        vm.prank(user1);
        token.payAccessFee("credential-123");
        
        // 4. User can view their NFT
        assertEq(nft.ownerOf(tokenId), user1);
        string memory uri = nft.tokenURI(tokenId);
        assertEq(uri, "ipfs://QmEncryptedShard");
        
        // 5. User deletes credential
        vm.prank(user1);
        nft.burnShard(tokenId);
    }
}

