const { ethers } = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  console.log("🚀 Starting Pass Chain Contract Deployment...\n");

  // Get network info
  const network = await ethers.provider.getNetwork();
  const [deployer] = await ethers.getSigners();
  const balance = await ethers.provider.getBalance(deployer.address);

  console.log("📋 Deployment Info:");
  console.log("  Network:", network.name, `(Chain ID: ${network.chainId})`);
  console.log("  Deployer:", deployer.address);
  console.log("  Balance:", ethers.formatEther(balance), "ETH");
  
  // Get config from environment
  const treasuryAddress = process.env.TREASURY_ADDRESS || deployer.address;
  const backendAddress = process.env.BACKEND_ADDRESS || deployer.address;
  
  console.log("  Treasury:", treasuryAddress);
  console.log("  Backend Minter:", backendAddress);
  console.log("\n");

  // 1. Deploy PCT Token
  console.log("1️⃣  Deploying PassChain Token (PCT)...");
  const PassChainToken = await ethers.getContractFactory("PassChainToken");
  const token = await PassChainToken.deploy(treasuryAddress);
  await token.waitForDeployment();
  const tokenAddress = await token.getAddress();
  
  console.log("   ✅ PCT Token deployed at:", tokenAddress);
  
  const totalSupply = await token.totalSupply();
  const storageFee = await token.storageFeePerCredential();
  const accessFee = await token.accessFeePerRetrieval();
  
  console.log("   📊 Initial supply:", ethers.formatEther(totalSupply), "PCT");
  console.log("   💰 Storage fee:", ethers.formatEther(storageFee), "PCT");
  console.log("   💰 Access fee:", ethers.formatEther(accessFee), "PCT");
  console.log("\n");

  // 2. Deploy Shard NFT
  console.log("2️⃣  Deploying PassChain Shard NFT (PCS)...");
  const PassChainShardNFT = await ethers.getContractFactory("PassChainShardNFT");
  const nft = await PassChainShardNFT.deploy(backendAddress);
  await nft.waitForDeployment();
  const nftAddress = await nft.getAddress();
  
  console.log("   ✅ PCS NFT deployed at:", nftAddress);
  console.log("   🔐 Authorized minter:", await nft.minter());
  console.log("\n");

  // 3. Transfer initial tokens to treasury
  console.log("3️⃣  Distributing tokens...");
  const treasuryAllocation = (totalSupply * 20n) / 100n; // 20% to treasury
  
  if (treasuryAddress !== deployer.address) {
    const tx = await token.transfer(treasuryAddress, treasuryAllocation);
    await tx.wait();
    console.log("   ✅ Transferred", ethers.formatEther(treasuryAllocation), "PCT to treasury");
  } else {
    console.log("   ℹ️  Treasury is deployer, no transfer needed");
  }
  console.log("\n");

  // 4. Save deployment info
  const deploymentInfo = {
    network: {
      name: network.name,
      chainId: Number(network.chainId),
      timestamp: new Date().toISOString(),
    },
    deployer: deployer.address,
    treasury: treasuryAddress,
    backend: backendAddress,
    contracts: {
      PassChainToken: {
        address: tokenAddress,
        totalSupply: ethers.formatEther(totalSupply),
        storageFee: ethers.formatEther(storageFee),
        accessFee: ethers.formatEther(accessFee),
      },
      PassChainShardNFT: {
        address: nftAddress,
        minter: backendAddress,
      },
    },
  };

  const deploymentsDir = path.join(__dirname, "../deployments");
  if (!fs.existsSync(deploymentsDir)) {
    fs.mkdirSync(deploymentsDir, { recursive: true });
  }

  const filename = path.join(
    deploymentsDir,
    `${network.name}-${network.chainId}-${Date.now()}.json`
  );
  
  fs.writeFileSync(filename, JSON.stringify(deploymentInfo, null, 2));

  // Print summary
  console.log("═════════════════════════════════════════════");
  console.log("✅ DEPLOYMENT COMPLETE");
  console.log("═════════════════════════════════════════════");
  console.log("Network:", network.name);
  console.log("Chain ID:", network.chainId);
  console.log("\n📄 Contracts:");
  console.log("  PCT Token:", tokenAddress);
  console.log("  PCS NFT:", nftAddress);
  console.log("\n💾 Deployment info saved to:");
  console.log("  ", filename);
  
  // Verification instructions
  if (network.chainId !== 31337n) {
    console.log("\n🔍 Verify contracts with:");
    console.log(`  npx hardhat verify --network ${network.name} ${tokenAddress} "${treasuryAddress}"`);
    console.log(`  npx hardhat verify --network ${network.name} ${nftAddress} "${backendAddress}"`);
  }
  
  console.log("\n📝 Update your .env with:");
  console.log(`NEXT_PUBLIC_PCT_ADDRESS=${tokenAddress}`);
  console.log(`NEXT_PUBLIC_NFT_ADDRESS=${nftAddress}`);
  
  console.log("\n═════════════════════════════════════════════\n");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error("❌ Deployment failed:", error);
    process.exit(1);
  });

