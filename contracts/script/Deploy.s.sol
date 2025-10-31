// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/PassChainToken.sol";
import "../src/PassChainShardNFT.sol";

/**
 * @title DeployScript
 * @dev Deployment script for Pass Chain contracts
 * 
 * Usage:
 * 1. Local (Anvil):
 *    forge script script/Deploy.s.sol:DeployScript --rpc-url http://localhost:8545 --broadcast
 * 
 * 2. Testnet (Sepolia):
 *    forge script script/Deploy.s.sol:DeployScript --rpc-url $SEPOLIA_RPC_URL --broadcast --verify
 * 
 * 3. Mainnet:
 *    forge script script/Deploy.s.sol:DeployScript --rpc-url $MAINNET_RPC_URL --broadcast --verify --slow
 */
contract DeployScript is Script {
    // Deployment addresses (will be set during deployment)
    PassChainToken public token;
    PassChainShardNFT public nft;
    
    // Configuration
    address public treasury;
    address public minter;
    
    function setUp() public {
        // Set treasury and minter addresses based on network
        uint256 chainId = block.chainid;
        
        if (chainId == 1) {
            // Mainnet - use multisig
            treasury = vm.envAddress("TREASURY_ADDRESS");
            minter = vm.envAddress("BACKEND_ADDRESS");
        } else if (chainId == 11155111) {
            // Sepolia testnet
            treasury = vm.envAddress("TREASURY_ADDRESS");
            minter = vm.envAddress("BACKEND_ADDRESS");
        } else {
            // Local / Anvil - use deployer
            treasury = msg.sender;
            minter = msg.sender;
        }
    }
    
    function run() public {
        // Get deployer private key
        uint256 deployerPrivateKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerPrivateKey);
        
        console.log("Deploying contracts with account:", deployer);
        console.log("Account balance:", deployer.balance);
        console.log("Treasury:", treasury);
        console.log("Minter:", minter);
        
        vm.startBroadcast(deployerPrivateKey);
        
        // 1. Deploy PCT Token
        console.log("\n1. Deploying PassChain Token (PCT)...");
        token = new PassChainToken(treasury);
        console.log("  PCT Token deployed at:", address(token));
        console.log("  Initial supply:", token.totalSupply() / 10**18, "PCT");
        console.log("  Storage fee:", token.storageFeePerCredential() / 10**18, "PCT");
        console.log("  Access fee:", token.accessFeePerRetrieval() / 10**18, "PCT");
        
        // 2. Deploy Shard NFT
        console.log("\n2. Deploying PassChain Shard NFT (PCS)...");
        nft = new PassChainShardNFT(minter);
        console.log("  PCS NFT deployed at:", address(nft));
        console.log("  Authorized minter:", nft.minter());
        
        // 3. Transfer initial tokens to treasury
        uint256 treasuryAllocation = (token.totalSupply() * 20) / 100; // 20% to treasury
        token.transfer(treasury, treasuryAllocation);
        console.log("\n3. Transferred", treasuryAllocation / 10**18, "PCT to treasury");
        
        vm.stopBroadcast();
        
        // Print deployment summary
        console.log("\n========================================");
        console.log("DEPLOYMENT COMPLETE");
        console.log("========================================");
        console.log("Network:", getNetworkName());
        console.log("Deployer:", deployer);
        console.log("Treasury:", treasury);
        console.log("Backend Minter:", minter);
        console.log("");
        console.log("PassChain Token (PCT):", address(token));
        console.log("PassChain Shard NFT (PCS):", address(nft));
        console.log("========================================");
        
        // Save deployment addresses
        saveDeployment();
    }
    
    function getNetworkName() internal view returns (string memory) {
        uint256 chainId = block.chainid;
        if (chainId == 1) return "Ethereum Mainnet";
        if (chainId == 11155111) return "Sepolia Testnet";
        if (chainId == 137) return "Polygon Mainnet";
        if (chainId == 42161) return "Arbitrum One";
        if (chainId == 31337) return "Local Anvil";
        return "Unknown Network";
    }
    
    function saveDeployment() internal {
        string memory json = string(
            abi.encodePacked(
                '{\n',
                '  "network": "', getNetworkName(), '",\n',
                '  "chainId": ', vm.toString(block.chainid), ',\n',
                '  "deployer": "', vm.toString(msg.sender), '",\n',
                '  "treasury": "', vm.toString(treasury), '",\n',
                '  "minter": "', vm.toString(minter), '",\n',
                '  "contracts": {\n',
                '    "PassChainToken": "', vm.toString(address(token)), '",\n',
                '    "PassChainShardNFT": "', vm.toString(address(nft)), '"\n',
                '  },\n',
                '  "timestamp": ', vm.toString(block.timestamp), '\n',
                '}'
            )
        );
        
        string memory filename = string(
            abi.encodePacked(
                "deployments/",
                vm.toString(block.chainid),
                "-",
                vm.toString(block.timestamp),
                ".json"
            )
        );
        
        vm.writeFile(filename, json);
        console.log("\nDeployment info saved to:", filename);
    }
}

