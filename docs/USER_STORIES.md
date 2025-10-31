# 📖 Pass Chain Enterprise - User Stories

## 🎯 Overview
This document contains user stories for Pass Chain Enterprise, organized by persona and feature area.

---

## 👤 Personas

### 1. **Alice** - Startup Founder / Owner
- Needs to securely store company API keys
- Wants to share credentials with team
- Requires audit trail for compliance

### 2. **Bob** - DevOps Engineer / Admin
- Manages infrastructure secrets
- Onboards new developers
- Rotates keys regularly

### 3. **Carol** - Developer / Member
- Needs access to staging/prod credentials
- Works on multiple projects
- Only needs credentials for assigned projects

### 4. **Dave** - Security Officer / Auditor
- Reviews access logs
- Monitors suspicious activity
- Generates compliance reports

### 5. **Eve** - Contractor / Guest
- Temporary access to specific credentials
- Limited permissions
- Time-bound access

---

## 🏢 Organization Management

### Epic 1: Organization Setup

#### US-1.1: Create Organization
**As Alice (Founder)**, I want to create a new organization so that I can invite my team and manage shared credentials.

**Acceptance Criteria:**
- [ ] I can connect my wallet
- [ ] I can enter organization name
- [ ] System creates default vault
- [ ] System assigns me as Owner
- [ ] I see organization dashboard

**Technical:**
- POST `/api/v1/organizations`
- Creates org, default roles, default vault
- Creator gets Owner role

---

#### US-1.2: View My Organizations
**As Carol (Member)**, I want to see all organizations I belong to so that I can switch between them.

**Acceptance Criteria:**
- [ ] I see a list of all my organizations
- [ ] Each shows: name, role, member count
- [ ] I can click to enter org dashboard
- [ ] I can filter by role or name

**Technical:**
- GET `/api/v1/organizations`
- Returns orgs where user is active member

---

### Epic 2: Team Management

#### US-2.1: Invite Team Member
**As Bob (Admin)**, I want to invite a new developer to my organization so they can access project credentials.

**Acceptance Criteria:**
- [ ] I can enter email OR wallet address
- [ ] I can select a role (Member, Admin, etc.)
- [ ] System generates invite link
- [ ] Invite expires after 7 days
- [ ] I receive confirmation

**Technical:**
- POST `/api/v1/organizations/:id/invite`
- Generates unique token
- Optional: Send email notification

---

#### US-2.2: Accept Invitation
**As Carol (New Member)**, I want to accept an invitation to join an organization.

**Acceptance Criteria:**
- [ ] I receive invite link
- [ ] I connect my wallet
- [ ] I see org details (name, inviter)
- [ ] I can accept or decline
- [ ] On accept, I'm added to org

**Technical:**
- GET `/api/v1/invitations/:token` (preview)
- POST `/api/v1/invitations/:token/accept`

---

#### US-2.3: Manage Member Roles
**As Bob (Admin)**, I want to change a member's role when they get promoted.

**Acceptance Criteria:**
- [ ] I see list of org members
- [ ] I can change anyone's role (except last Owner)
- [ ] Member immediately gets new permissions
- [ ] Action is logged in audit

**Technical:**
- PUT `/api/v1/organizations/:id/members/:userId/role`
- Permission: member:update

---

#### US-2.4: Remove Team Member
**As Alice (Owner)**, I want to remove a contractor when their contract ends.

**Acceptance Criteria:**
- [ ] I see list of members
- [ ] I can remove any member (except myself if last Owner)
- [ ] Removed member loses all org access immediately
- [ ] Their personal vault remains intact
- [ ] Action is logged

**Technical:**
- DELETE `/api/v1/organizations/:id/members/:userId`
- Status changed to 'removed'

---

## 🗂️ Vault & Project Management

### Epic 3: Vault Organization

#### US-3.1: Create Project
**As Bob (Admin)**, I want to create a project for the mobile app team so their credentials are isolated.

**Acceptance Criteria:**
- [ ] I can enter project name & description
- [ ] System creates project vault automatically
- [ ] I can see project in projects list
- [ ] Only assigned members can access

**Technical:**
- POST `/api/v1/organizations/:id/projects`
- Creates project + project vault

---

#### US-3.2: View Available Vaults
**As Carol (Member)**, I want to see all vaults I have access to.

**Acceptance Criteria:**
- [ ] I see: Personal vault, Org vaults, Project vaults
- [ ] Each shows: name, type, credential count
- [ ] I can filter by type
- [ ] I can switch between vaults

**Technical:**
- GET `/api/v1/vaults`
- Returns vaults based on:
  - Personal: owner_user_id = current_user
  - Org/Project: Check permissions + vault_access

---

#### US-3.3: Share Vault Access
**As Bob (Admin)**, I want to grant Carol access to the staging vault.

**Acceptance Criteria:**
- [ ] I select a vault
- [ ] I select user(s) to share with
- [ ] Optionally select specific role
- [ ] Users immediately get access
- [ ] Action is logged

**Technical:**
- POST `/api/v1/vaults/:id/access`
- Creates vault_access record

---

## 🔐 Credential Management

### Epic 4: Store & Retrieve Credentials

#### US-4.1: Store Credential in Org Vault
**As Carol (Member)**, I want to store the staging database password in the project vault.

**Acceptance Criteria:**
- [ ] I select target vault
- [ ] I enter: name, username, password, URL, tags
- [ ] System encrypts client-side
- [ ] System splits key into shards
- [ ] Shard 1 → Vault, Shard 2 → DB (or Fabric future), Shard 3 → Browser
- [ ] I see success confirmation

**Technical:**
- POST `/api/v1/credentials`
- Body includes `vaultId`
- Check vault write permission

---

#### US-4.2: Retrieve Credential from Shared Vault
**As Carol (Member)**, I want to retrieve the staging DB password from the project vault.

**Acceptance Criteria:**
- [ ] I browse project vault
- [ ] I click "Reveal Password"
- [ ] System prompts wallet signature
- [ ] System fetches shards from Vault + DB
- [ ] Client decrypts password
- [ ] Access is logged to audit

**Technical:**
- GET `/api/v1/credentials/:id`
- Check credential access (via vault)
- Retrieve shards
- Log audit entry

---

#### US-4.3: Filter Credentials by Vault
**As Bob (Admin)**, I want to see only credentials in the production vault.

**Acceptance Criteria:**
- [ ] I select vault from dropdown
- [ ] Credentials list filters to that vault
- [ ] I can search within vault
- [ ] I can filter by tags

**Technical:**
- GET `/api/v1/credentials?vaultId=xxx`

---

## 📊 Audit & Security

### Epic 5: Audit Logs

#### US-5.1: View Organization Audit Log
**As Dave (Security Officer)**, I want to see all credential access in my organization.

**Acceptance Criteria:**
- [ ] I navigate to org audit page
- [ ] I see table: timestamp, user, action, resource, IP
- [ ] I can filter by: user, action type, date range
- [ ] I can export to CSV
- [ ] Logs show Fabric TX hash (if anchored)

**Technical:**
- GET `/api/v1/organizations/:id/audit`
- Permission: org:audit
- Returns audit_logs filtered by org_id

---

#### US-5.2: View Credential History
**As Carol (Member)**, I want to see who accessed the production API key.

**Acceptance Criteria:**
- [ ] I open credential detail
- [ ] I click "History" tab
- [ ] I see: timestamp, action, user
- [ ] I see hashed IP address
- [ ] I see Fabric TX hash (future)

**Technical:**
- GET `/api/v1/credentials/:id/audit`
- Returns audit_logs for credential_id

---

#### US-5.3: View Vault Audit Log
**As Bob (Admin)**, I want to audit all actions on the production vault.

**Acceptance Criteria:**
- [ ] I open vault settings
- [ ] I see audit tab
- [ ] Logs show: credentials created/deleted/accessed in vault
- [ ] I can filter by action

**Technical:**
- GET `/api/v1/vaults/:id/audit`

---

#### US-5.4: Read Audit Logs from Vault (Backup)
**As Dave (Security Officer)**, I want to retrieve audit logs directly from HashiCorp Vault as a backup verification method.

**Acceptance Criteria:**
- [ ] System stores audit log references in Vault KV
- [ ] Each audit entry has: timestamp, user, action, resource_id
- [ ] I can query Vault via API: `GET /v1/secret/data/audit/:orgId/:year/:month`
- [ ] Logs match database records (verification)
- [ ] Logs are immutable once written

**Technical:**
- Write to Vault path: `passchain/audit/:orgId/:credentialId/:timestamp`
- Store: action, user_id, ip_hash, metadata
- Read via Vault API or backend endpoint

---

### Epic 6: Security & Compliance

#### US-6.1: Rotate Credential
**As Bob (Admin)**, I want to rotate the production API key.

**Acceptance Criteria:**
- [ ] I click "Rotate" on credential
- [ ] System generates new encryption key
- [ ] Old shards invalidated
- [ ] New shards created
- [ ] Rotation logged in audit
- [ ] Blockchain records rotation (future)

**Technical:**
- POST `/api/v1/credentials/:id/rotate`
- Permission: credential:rotate

---

#### US-6.2: Generate Compliance Report
**As Dave (Security Officer)**, I want to export last 90 days of audit logs for SOC2 audit.

**Acceptance Criteria:**
- [ ] I select date range
- [ ] I select report type (CSV, PDF)
- [ ] Report includes: all access, changes, member updates
- [ ] Report shows Fabric anchors (future)
- [ ] Report is timestamped and signed

**Technical:**
- GET `/api/v1/organizations/:id/compliance-report`
- Query audit_logs, generate PDF

---

## 🔑 RBAC & Permissions

### Epic 7: Role Management

#### US-7.1: Create Custom Role
**As Alice (Owner)**, I want to create a "DevOps Lead" role with specific permissions.

**Acceptance Criteria:**
- [ ] I click "Create Role"
- [ ] I enter role name & description
- [ ] I select permissions from checklist
- [ ] Role appears in role list
- [ ] I can assign to members

**Technical:**
- POST `/api/v1/organizations/:id/roles`
- Creates role + permissions

---

#### US-7.2: View My Permissions
**As Carol (Member)**, I want to see what I'm allowed to do.

**Acceptance Criteria:**
- [ ] I navigate to "My Permissions"
- [ ] I see my role name
- [ ] I see list of permissions (create credential, read vault, etc.)
- [ ] I see which vaults I can access

**Technical:**
- GET `/api/v1/me/permissions?orgId=xxx`

---

## 🌐 Advanced Features (Future)

### Epic 8: Blockchain Integration

#### US-8.1: Anchor Audit Log to Fabric
**As Dave (Security Officer)**, I want all audit logs anchored to blockchain for immutability.

**Acceptance Criteria:**
- [ ] Every credential access creates Fabric transaction
- [ ] TX hash stored in audit_logs.fabric_tx_hash
- [ ] I can verify log integrity via Fabric
- [ ] Tampered logs are detectable

---

#### US-8.2: View Blockchain Audit Trail
**As Dave**, I want to view the Fabric blockchain for audit verification.

**Acceptance Criteria:**
- [ ] I navigate to "Blockchain Audit"
- [ ] I see list of transactions
- [ ] I can click TX hash to see details
- [ ] I can verify credential access happened

---

### Epic 9: Token Economy (Future)

#### US-9.1: Pay for Storage with Token
**As Alice (Owner)**, I want to pay for credential storage using PassChain tokens.

**Acceptance Criteria:**
- [ ] I connect wallet with PCT tokens
- [ ] I approve token spend
- [ ] Credential is stored after payment
- [ ] Payment recorded on-chain

---

## 📱 Mobile & Extension (Future)

### Epic 10: Browser Extension

#### US-10.1: Autofill Credentials
**As Carol (Member)**, I want my browser extension to autofill stored credentials.

**Acceptance Criteria:**
- [ ] Extension detects login form
- [ ] Shows matching credentials from accessible vaults
- [ ] I select credential
- [ ] Extension fills username & password
- [ ] Access is logged

---

## 🎯 Summary

### By Phase:

**Phase 1 (MVP - Done):**
- ✅ US-4.1, US-4.2 (Personal vault only)
- ✅ US-5.2 (Basic audit)

**Phase 2 (Organizations):**
- US-1.1, US-1.2 (Org setup)
- US-2.1, US-2.2, US-2.3, US-2.4 (Team management)
- US-3.1, US-3.2, US-3.3 (Vault/project management)
- US-4.1, US-4.2, US-4.3 (Vault-scoped credentials)
- US-5.1, US-5.3, US-5.4 (Org audit logs + Vault backup)

**Phase 3 (Enterprise):**
- US-5.4 (Vault audit log storage)
- US-6.1, US-6.2 (Security & compliance)
- US-7.1, US-7.2 (Custom roles)

**Phase 4 (Blockchain):**
- US-8.1, US-8.2 (Fabric integration)
- US-9.1 (Token payments)

**Phase 5 (Extensions):**
- US-10.1 (Browser extension)

---

## 📊 Story Count by Persona

| Persona | Story Count | Priority |
|---------|------------|----------|
| Alice (Owner) | 5 | High |
| Bob (Admin) | 7 | High |
| Carol (Member) | 6 | High |
| Dave (Security) | 5 | Medium |
| Eve (Guest) | 1 | Low |

**Total:** 24 user stories (18 MVP/Enterprise, 6 Future)

---

**AUUUUFFFF!** 🔥

