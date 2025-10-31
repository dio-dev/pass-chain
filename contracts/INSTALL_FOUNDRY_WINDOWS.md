# Installing Foundry on Windows - Complete Guide

## ✅ Option 1: Native Windows Installation (Easiest)

Foundry now has native Windows support! Here's how:

### Step 1: Download Foundryup for Windows

Open PowerShell and run:

```powershell
# Download the Windows installer
Invoke-WebRequest -Uri "https://github.com/foundry-rs/foundry/releases/latest/download/foundry_nightly_windows_amd64.tar.gz" -OutFile "foundry.tar.gz"

# Extract it
tar -xzf foundry.tar.gz

# Move to a permanent location
New-Item -ItemType Directory -Force -Path "$env:USERPROFILE\.foundry\bin"
Move-Item -Force forge.exe "$env:USERPROFILE\.foundry\bin\"
Move-Item -Force cast.exe "$env:USERPROFILE\.foundry\bin\"
Move-Item -Force anvil.exe "$env:USERPROFILE\.foundry\bin\"
Move-Item -Force chisel.exe "$env:USERPROFILE\.foundry\bin\"

# Add to PATH (permanent)
$oldPath = [Environment]::GetEnvironmentVariable("Path", "User")
$newPath = "$oldPath;$env:USERPROFILE\.foundry\bin"
[Environment]::SetEnvironmentVariable("Path", $newPath, "User")

# Reload PATH for current session
$env:Path = [Environment]::GetEnvironmentVariable("Path", "User")
```

### Step 2: Verify Installation

```powershell
# Restart PowerShell, then check:
forge --version
cast --version
anvil --version
```

---

## ✅ Option 2: Using Scoop Package Manager

If you have [Scoop](https://scoop.sh/) installed:

```powershell
scoop install foundry
```

To install Scoop first:
```powershell
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex
```

---

## ✅ Option 3: Using WSL2 (If you prefer Linux)

```powershell
# Open WSL
wsl

# In WSL, run:
curl -L https://foundry.paradigm.xyz | bash
source ~/.bashrc
foundryup

# Your project is accessible at:
cd /mnt/c/Users/alex/projects/pass-chain/contracts
```

---

## 🚀 Quick Start After Installation

```powershell
cd C:\Users\alex\projects\pass-chain\contracts

# Install OpenZeppelin dependencies
forge install OpenZeppelin/openzeppelin-contracts

# Run tests
forge test -vvv

# Start local blockchain (in separate terminal)
anvil

# Deploy to local blockchain
forge script script/Deploy.s.sol:DeployScript --rpc-url http://localhost:8545 --broadcast
```

---

## 🐛 Troubleshooting

### "forge: command not found"
- Restart PowerShell after installation
- Check PATH: `$env:Path`
- Manually add to PATH if needed

### "Missing git"
```powershell
# Install Git for Windows
winget install Git.Git
# Or download from: https://git-scm.com/download/win
```

### "OpenSSL error"
Download OpenSSL for Windows: https://slproweb.com/products/Win32OpenSSL.html

---

## 📝 What You'll Be Able to Do

Once Foundry is installed:

1. **Test Smart Contracts Locally**
```powershell
forge test -vvv
```

2. **Deploy to Local Blockchain**
```powershell
# Terminal 1: Start Anvil
anvil

# Terminal 2: Deploy
forge script script/Deploy.s.sol:DeployScript --rpc-url http://localhost:8545 --broadcast
```

3. **Deploy to Sepolia Testnet**
```powershell
forge script script/Deploy.s.sol:DeployScript --rpc-url $env:SEPOLIA_RPC_URL --broadcast --verify
```

4. **Interact with Contracts**
```powershell
cast call 0x... "balanceOf(address)" 0xYourAddress
```

---

**Choose Option 1 (Native Windows) for best experience!**

**AUUUUFFFF!** 🔥

