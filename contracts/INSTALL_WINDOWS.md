# Install Foundry on Windows

## Option 1: Using Foundryup Installer (Recommended)

Open PowerShell and run:

```powershell
# Download and run the installer
curl.exe -L https://foundry.paradigm.xyz | bash

# Then restart your terminal and run:
foundryup
```

## Option 2: Using Scoop (Package Manager)

```powershell
# Install Scoop if you don't have it
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Install Foundry
scoop install foundry
```

## Option 3: Direct Binary Download

1. Go to https://github.com/foundry-rs/foundry/releases
2. Download `foundry_nightly_windows_amd64.zip`
3. Extract to `C:\foundry\`
4. Add to PATH:
```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\foundry\bin", "User")
```

## Verify Installation

```powershell
forge --version
cast --version
anvil --version
```

## Quick Start

```powershell
cd C:\Users\alex\projects\pass-chain\contracts

# Install OpenZeppelin
forge install OpenZeppelin/openzeppelin-contracts

# Run tests
forge test -vvv

# Deploy locally
# Terminal 1:
anvil

# Terminal 2:
forge script script/Deploy.s.sol:DeployScript --rpc-url http://localhost:8545 --broadcast
```

---

**Try Option 1 with `curl.exe` (not `curl`)** - PowerShell's `curl` is different!

**AUUUUFFFF!** 🔥

