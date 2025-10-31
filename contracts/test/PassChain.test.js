const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("Pass Chain Contracts", function () {
  let token, nft;
  let owner, treasury, backend, user1, user2;
  
  const INITIAL_BALANCE = ethers.parseEther("1000");

  beforeEach(async function () {
    [owner, treasury, backend, user1, user2] = await ethers.getSigners();

    // Deploy PCT Token
    const PassChainToken = await ethers.getContractFactory("PassChainToken");
    token = await PassChainToken.deploy(treasury.address);
    await token.waitForDeployment();

    // Deploy Shard NFT
    const PassChainShardNFT = await ethers.getContractFactory("PassChainShardNFT");
    nft = await PassChainShardNFT.deploy(backend.address);
    await nft.waitForDeployment();

    // Give users tokens
    await token.transfer(user1.address, INITIAL_BALANCE);
    await token.transfer(user2.address, INITIAL_BALANCE);
  });

  describe("PassChainToken", function () {
    it("Should have correct initial values", async function () {
      expect(await token.name()).to.equal("PassChain Token");
      expect(await token.symbol()).to.equal("PCT");
      expect(await token.treasury()).to.equal(treasury.address);
      expect(await token.storageFeePerCredential()).to.equal(ethers.parseEther("10"));
      expect(await token.accessFeePerRetrieval()).to.equal(ethers.parseEther("1"));
    });

    it("Should allow paying storage fee", async function () {
      const initialBalance = await token.balanceOf(user1.address);
      const initialTreasury = await token.balanceOf(treasury.address);

      await token.connect(user1).payStorageFee("cred-123");

      expect(await token.balanceOf(user1.address)).to.equal(
        initialBalance - ethers.parseEther("10")
      );
      expect(await token.balanceOf(treasury.address)).to.equal(
        initialTreasury + ethers.parseEther("10")
      );
    });

    it("Should allow paying access fee", async function () {
      const initialBalance = await token.balanceOf(user1.address);
      const initialTreasury = await token.balanceOf(treasury.address);

      await token.connect(user1).payAccessFee("cred-123");

      expect(await token.balanceOf(user1.address)).to.equal(
        initialBalance - ethers.parseEther("1")
      );
      expect(await token.balanceOf(treasury.address)).to.equal(
        initialTreasury + ethers.parseEther("1")
      );
    });

    it("Should reject payment with insufficient balance", async function () {
      const [,,,,, poorUser] = await ethers.getSigners();
      await expect(
        token.connect(poorUser).payStorageFee("cred-123")
      ).to.be.revertedWith("Insufficient balance");
    });

    it("Should allow owner to update fees", async function () {
      await token.updateStorageFee(ethers.parseEther("20"));
      expect(await token.storageFeePerCredential()).to.equal(ethers.parseEther("20"));
    });

    it("Should allow pausing and unpausing", async function () {
      await token.pause();
      await expect(
        token.connect(user1).transfer(user2.address, 100)
      ).to.be.revertedWithCustomError(token, "EnforcedPause");

      await token.unpause();
      await token.connect(user1).transfer(user2.address, 100);
    });
  });

  describe("PassChainShardNFT", function () {
    it("Should have correct initial values", async function () {
      expect(await nft.name()).to.equal("PassChain Shard NFT");
      expect(await nft.symbol()).to.equal("PCS");
      expect(await nft.minter()).to.equal(backend.address);
    });

    it("Should allow minter to mint shard", async function () {
      await nft.connect(backend).mintShard(
        user1.address,
        "cred-123",
        "QmXyz123"
      );

      const tokenId = await nft.getTokenId("cred-123");
      expect(await nft.ownerOf(tokenId)).to.equal(user1.address);
      expect(await nft.getCredentialId(tokenId)).to.equal("cred-123");
      expect(await nft.tokenURI(tokenId)).to.equal("ipfs://QmXyz123");
    });

    it("Should reject duplicate shard minting", async function () {
      await nft.connect(backend).mintShard(user1.address, "cred-123", "QmXyz123");
      
      await expect(
        nft.connect(backend).mintShard(user1.address, "cred-123", "QmXyz456")
      ).to.be.revertedWith("Shard already exists");
    });

    it("Should reject non-minter minting", async function () {
      await expect(
        nft.connect(user1).mintShard(user1.address, "cred-123", "QmXyz123")
      ).to.be.revertedWith("Only minter can call this");
    });

    it("Should allow burning shard", async function () {
      await nft.connect(backend).mintShard(user1.address, "cred-123", "QmXyz123");
      const tokenId = await nft.getTokenId("cred-123");

      await nft.connect(user1).burnShard(tokenId);

      await expect(nft.ownerOf(tokenId)).to.be.revertedWithCustomError(
        nft,
        "ERC721NonexistentToken"
      );
    });

    it("Should track wallet shards", async function () {
      await nft.connect(backend).mintShard(user1.address, "cred-1", "QmHash1");
      await nft.connect(backend).mintShard(user1.address, "cred-2", "QmHash2");
      await nft.connect(backend).mintShard(user1.address, "cred-3", "QmHash3");

      const shards = await nft.getWalletShards(user1.address);
      expect(shards.length).to.equal(3);
    });

    it("Should allow NFT transfer", async function () {
      await nft.connect(backend).mintShard(user1.address, "cred-123", "QmXyz123");
      const tokenId = await nft.getTokenId("cred-123");

      await nft.connect(user1).transferFrom(user1.address, user2.address, tokenId);

      expect(await nft.ownerOf(tokenId)).to.equal(user2.address);
      expect((await nft.getWalletShards(user1.address)).length).to.equal(0);
      expect((await nft.getWalletShards(user2.address)).length).to.equal(1);
    });
  });

  describe("Integration", function () {
    it("Should complete full credential lifecycle", async function () {
      // 1. User pays storage fee
      await token.connect(user1).payStorageFee("cred-123");

      // 2. Backend mints NFT
      await nft.connect(backend).mintShard(user1.address, "cred-123", "QmEncrypted");

      // 3. User pays access fee
      await token.connect(user1).payAccessFee("cred-123");

      // 4. Verify NFT
      const tokenId = await nft.getTokenId("cred-123");
      expect(await nft.ownerOf(tokenId)).to.equal(user1.address);
      expect(await nft.tokenURI(tokenId)).to.equal("ipfs://QmEncrypted");

      // 5. Delete credential
      await nft.connect(user1).burnShard(tokenId);
    });
  });
});

