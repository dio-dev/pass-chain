# Pass Chain Contracts - Hardhat Deployment Guide

## 🚀 Quick Start

### 1. Install Dependencies
```bash
cd contracts
npm install
```

### 2. Setup Environment
```bash
cp .env.example .env
# Edit .env with your keys
```

### 3. Compile Contracts
```bash
npm run compile
```

### 4. Run Tests
```bash
npm test
```

---

## 🌐 Deployment

### Local Development (Hardhat Node)

**Terminal 1:** Start local blockchain
```bash
npm run node
```

**Terminal 2:** Deploy
```bash
npm run deploy:local
```

---

### Sepolia Testnet

**Setup `.env`:**
```env
PRIVATE_KEY=your_private_key_here
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
ETHERSCAN_API_KEY=your_etherscan_key
TREASURY_ADDRESS=0x...  # Your treasury wallet
BACKEND_ADDRESS=0x...   # Backend service address
```

**Deploy:**
```bash
npm run deploy:sepolia
```

**Verify Contracts:**
```bash
npx hardhat verify --network sepolia <PCT_ADDRESS> "<TREASURY_ADDRESS>"
npx hardhat verify --network sepolia <NFT_ADDRESS> "<BACKEND_ADDRESS>"
```

---

### Mainnet

**⚠️ CAUTION: Real money!**

```bash
# Make sure your .env has:
# - MAINNET_RPC_URL
# - PRIVATE_KEY (with sufficient ETH for gas)
# - TREASURY_ADDRESS (preferably a multisig)
# - BACKEND_ADDRESS

npm run deploy:mainnet
```

**Recommended:**
- Use hardware wallet (Ledger/Trezor)
- Test on Sepolia first
- Use a multisig for treasury
- Start with low gas prices

---

### Polygon

```bash
# Setup in .env:
# POLYGON_RPC_URL=https://polygon-rpc.com
# POLYGONSCAN_API_KEY=your_key

npm run deploy:polygon
```

---

## 📁 Deployment Output

After deployment, you'll get:

```
deployments/
└── sepolia-11155111-1234567890.json
```

**Example JSON:**
```json
{
  "network": {
    "name": "sepolia",
    "chainId": 11155111,
    "timestamp": "2024-10-24T12:00:00.000Z"
  },
  "deployer": "0x...",
  "treasury": "0x...",
  "backend": "0x...",
  "contracts": {
    "PassChainToken": {
      "address": "0x...",
      "totalSupply": "1000000000.0",
      "storageFee": "10.0",
      "accessFee": "1.0"
    },
    "PassChainShardNFT": {
      "address": "0x...",
      "minter": "0x..."
    }
  }
}
```

---

## 🔗 Frontend Integration

**Update `.env` in frontend:**
```env
NEXT_PUBLIC_PCT_ADDRESS=0x...
NEXT_PUBLIC_NFT_ADDRESS=0x...
```

**Use in app:**
```javascript
import { ethers } from 'ethers';
import PCT_ABI from './abis/PassChainToken.json';
import NFT_ABI from './abis/PassChainShardNFT.json';

const pct = new ethers.Contract(
  process.env.NEXT_PUBLIC_PCT_ADDRESS,
  PCT_ABI,
  signer
);

// Pay storage fee
await pct.payStorageFee(credentialId);
```

---

## 🔍 Verification

**Manual verification:**
```bash
npx hardhat verify --network sepolia \
  0xYourTokenAddress \
  "0xTreasuryAddress"

npx hardhat verify --network sepolia \
  0xYourNFTAddress \
  "0xBackendAddress"
```

**Check on Etherscan:**
- Sepolia: https://sepolia.etherscan.io/
- Mainnet: https://etherscan.io/
- Polygon: https://polygonscan.com/

---

## 🧪 Testing

```bash
# Run all tests
npm test

# Run with gas reporting
REPORT_GAS=true npm test

# Run specific test
npx hardhat test test/PassChain.test.js

# Coverage report
npx hardhat coverage
```

---

## 🛠️ Useful Commands

```bash
# Compile contracts
npm run compile

# Clean artifacts
npx hardhat clean

# Console (interact with contracts)
npx hardhat console --network sepolia

# Run a task
npx hardhat accounts
npx hardhat balances

# Get contract size
npx hardhat size-contracts
```

---

## 📊 Gas Estimates

| Action | Gas Used | Cost @ 50 gwei |
|--------|----------|----------------|
| Deploy PCT | ~2,000,000 | ~$5 |
| Deploy NFT | ~2,500,000 | ~$6 |
| Mint NFT | ~150,000 | ~$0.40 |
| Pay Storage Fee | ~50,000 | ~$0.13 |
| Pay Access Fee | ~45,000 | ~$0.12 |

---

## 🔐 Security

**Before Mainnet:**
- [ ] External audit (CertiK/OpenZeppelin)
- [ ] Test on Sepolia extensively
- [ ] Use hardware wallet for deployment
- [ ] Set up multisig for treasury
- [ ] Enable time locks for admin functions
- [ ] Set up monitoring/alerts

**Post-Deployment:**
- Monitor contract events
- Set up Tenderly alerts
- Regular security checks
- Bug bounty program

---

## 📝 Notes

- **Hardhat** is more popular than Foundry for Node.js projects
- Fully compatible with your existing Solidity contracts
- Better TypeScript support
- Easier to integrate with existing Node.js backend

---

**No need for Rust/Cargo! Just Node.js!**

**AUUUUFFFF!** 🔥

