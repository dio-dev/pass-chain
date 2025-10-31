# Pass Chain Smart Contracts 🔐

Solidity smart contracts for Pass Chain password manager ecosystem.

## Contracts

### 1. PassChainToken (PCT) - ERC-20
**Address:** TBD (deploy first)

**Purpose:**
- Payment token for credential storage
- Payment token for credential access
- Staking for premium features
- Governance (future)

**Features:**
- Initial supply: 1 billion PCT
- Storage fee: 10 PCT per credential
- Access fee: 1 PCT per retrieval
- Pausable for emergencies
- Burnable
- Mintable by owner (for liquidity/rewards)

### 2. PassChainShardNFT (PCS) - ERC-721
**Address:** TBD (deploy first)

**Purpose:**
- NFT representing encrypted Share3 backup on IPFS
- Permanent, decentralized backup solution
- Transferable (can gift credentials)
- Tradeable (secondary market potential)

**Features:**
- Each NFT = 1 encrypted shard on IPFS
- Metadata includes credential reference
- Only authorized backend can mint
- Owner can burn (when deleting credential)
- Tracks all shards per wallet

## Architecture

```
User Creates Credential:
1. Pay 10 PCT storage fee
   ↓
2. Backend encrypts Share3
   ↓
3. Upload to IPFS → Get hash (QmXyz...)
   ↓
4. Mint NFT with IPFS hash as metadata
   ↓
5. User receives NFT in wallet

User Accesses Credential:
1. Pay 1 PCT access fee
   ↓
2. Backend fetches Share1 (Vault) + Share2 (DB)
   ↓
3. User already has NFT (can get Share3 from IPFS if needed)
   ↓
4. Reconstruct key & decrypt

User Loses Everything:
1. Still has NFT in wallet
   ↓
2. View NFT metadata → Get IPFS hash
   ↓
3. Download from IPFS → Get encrypted Share3
   ↓
4. Decrypt with wallet signature
   ↓
5. Recover credentials!
```

## Deployment

### Prerequisites
```bash
# Install Foundry
curl -L https://foundry.paradigm.xyz | bash
foundryup

# Install dependencies
cd contracts
forge install OpenZeppelin/openzeppelin-contracts
```

### Environment Setup
Create `.env`:
```bash
# Private key (for deployment)
PRIVATE_KEY=your_private_key_here

# RPC URLs
SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/YOUR_KEY
MAINNET_RPC_URL=https://mainnet.infura.io/v3/YOUR_KEY
POLYGON_RPC_URL=https://polygon-rpc.com

# Addresses
TREASURY_ADDRESS=0x...  # Multisig for mainnet
BACKEND_ADDRESS=0x...   # Backend service address

# API Keys (for verification)
ETHERSCAN_API_KEY=your_etherscan_key
POLYGONSCAN_API_KEY=your_polygonscan_key
```

### Deploy to Local (Anvil)
```bash
# Start local node
anvil

# Deploy
forge script script/Deploy.s.sol:DeployScript \
  --rpc-url http://localhost:8545 \
  --broadcast
```

### Deploy to Sepolia Testnet
```bash
forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $SEPOLIA_RPC_URL \
  --broadcast \
  --verify
```

### Deploy to Mainnet
```bash
# CAUTION: Real money!
forge script script/Deploy.s.sol:DeployScript \
  --rpc-url $MAINNET_RPC_URL \
  --broadcast \
  --verify \
  --slow  # Use slow mode for better gas prices
```

## Testing

```bash
# Run all tests
forge test -vvv

# Run specific test
forge test --match-test testPayStorageFee -vvv

# Gas report
forge test --gas-report

# Coverage
forge coverage
```

## Contract Interactions

### For Users (Frontend)

**Buy PCT Tokens:**
```solidity
// Use Uniswap or other DEX
```

**Approve token spending:**
```javascript
const pct = new ethers.Contract(PCT_ADDRESS, ABI, signer);
await pct.approve(PCT_ADDRESS, ethers.utils.parseEther("1000"));
```

**Pay storage fee:**
```javascript
await pct.payStorageFee("credential-id-123");
```

**Pay access fee:**
```javascript
await pct.payAccessFee("credential-id-123");
```

**View your NFTs:**
```javascript
const nft = new ethers.Contract(NFT_ADDRESS, ABI, provider);
const shards = await nft.getWalletShards(userAddress);
```

**Get IPFS hash from NFT:**
```javascript
const tokenURI = await nft.tokenURI(tokenId);
// tokenURI = "ipfs://QmXyz..." 
// Download from IPFS to get encrypted Share3
```

### For Backend (Minting NFTs)

**Mint shard NFT:**
```javascript
const nft = new ethers.Contract(NFT_ADDRESS, ABI, backendSigner);

// 1. Encrypt Share3
const encryptedShare3 = encryptShare3(share3, userWallet);

// 2. Upload to IPFS
const ipfsHash = await uploadToIPFS(encryptedShare3);

// 3. Mint NFT
const tx = await nft.mintShard(
  userAddress,
  credentialId,
  ipfsHash
);

await tx.wait();
```

## Security

### Audits
- [ ] Internal audit complete
- [ ] External audit (CertiK/OpenZeppelin)
- [ ] Bug bounty program

### Best Practices
- ✅ OpenZeppelin contracts (battle-tested)
- ✅ Pausable for emergencies
- ✅ Access control (Ownable, minter role)
- ✅ Reentrancy protection (via OpenZeppelin)
- ✅ Comprehensive tests
- ✅ Gas optimizations

### Admin Functions
Only owner can:
- Update fees
- Update treasury address
- Pause/unpause token
- Mint new tokens
- Update NFT minter address

## Token Economics

### PCT Distribution
- 50% - Public sale / DEX liquidity
- 20% - Treasury (fees collection)
- 15% - Team (4-year vesting)
- 10% - Ecosystem rewards
- 5% - Initial users airdrop

### Fee Structure
- Storage: 10 PCT per credential (~$1 at launch)
- Access: 1 PCT per retrieval (~$0.10)
- Adjustable by governance (future)

### NFT Value
- **Utility:** Backup recovery key
- **Rarity:** Each credential = unique NFT
- **Transferable:** Can sell/gift credentials
- **Market:** Secondary market on OpenSea

## Integrations

### IPFS
- Store encrypted Share3
- Pinned by Pinata/Web3.Storage
- Redundant backups

### Frontend
```javascript
import { PassChainSDK } from '@passchain/sdk';

const sdk = new PassChainSDK({
  pctAddress: '0x...',
  nftAddress: '0x...',
  provider: window.ethereum
});

// Create credential with payment
await sdk.createCredential({
  name: 'GitHub',
  username: 'user@email.com',
  password: 'encrypted...'
});
```

### Backend
- Monitor blockchain events
- Verify payments before storing
- Mint NFTs automatically
- Track usage for analytics

## Roadmap

### Phase 1 (Q1 2025) ✅
- [x] Token contract
- [x] NFT contract
- [x] Deployment scripts
- [x] Tests

### Phase 2 (Q2 2025)
- [ ] Deploy to testnet
- [ ] Frontend integration
- [ ] Backend integration
- [ ] IPFS setup

### Phase 3 (Q3 2025)
- [ ] External audit
- [ ] Mainnet deployment
- [ ] DEX listing (Uniswap)
- [ ] Public launch

### Phase 4 (Q4 2025)
- [ ] Staking mechanism
- [ ] Governance (DAO)
- [ ] Premium features
- [ ] Mobile SDK

## Support

- 📧 Email: contracts@passchain.io
- 💬 Discord: [discord.gg/passchain](https://discord.gg/passchain)
- 📖 Docs: [docs.passchain.io/contracts](https://docs.passchain.io/contracts)

---

**Built with Foundry** ⚒️

**AUUUUFFFF!** 🔥

