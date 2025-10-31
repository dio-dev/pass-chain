# Pass Chain Enterprise - Database Schema

## Core Entities

### Users (Wallet-based Identity)
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_address VARCHAR(42) UNIQUE NOT NULL,
    ens_name VARCHAR(255),
    display_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    device_fingerprint TEXT,
    
    INDEX idx_wallet (wallet_address)
);
```

### Organizations
```sql
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    plan VARCHAR(50) DEFAULT 'free', -- free, pro, enterprise
    max_members INT DEFAULT 5,
    max_vaults INT DEFAULT 10,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_slug (slug)
);
```

### Organization Members
```sql
CREATE TABLE organization_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id),
    status VARCHAR(20) DEFAULT 'invited', -- invited, active, suspended, removed
    invited_by UUID REFERENCES users(id),
    joined_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(org_id, user_id),
    INDEX idx_org_user (org_id, user_id),
    INDEX idx_status (org_id, status)
);
```

### Roles (RBAC)
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN DEFAULT false, -- true for default roles (Owner, Admin, etc.)
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(org_id, name),
    INDEX idx_org_role (org_id, name)
);

-- Default system roles for every org
INSERT INTO roles (name, is_system, description) VALUES
    ('Owner', true, 'Full control over organization'),
    ('Admin', true, 'Manage members, roles, and vaults'),
    ('Security Officer', true, 'View audit logs and security settings'),
    ('Member', true, 'Read/write assigned vaults'),
    ('Auditor', true, 'Read-only access + audit logs'),
    ('Guest', true, 'Limited read access to shared credentials');
```

### Permissions
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    resource_type VARCHAR(50) NOT NULL, -- org, project, vault, credential
    action VARCHAR(50) NOT NULL, -- create, read, update, delete, share, rotate, invite, audit
    
    UNIQUE(role_id, resource_type, action),
    INDEX idx_role_permissions (role_id)
);
```

### Projects (Team Workspaces)
```sql
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_org_project (org_id)
);
```

### Vaults (Credential Containers)
```sql
CREATE TABLE vaults (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE, -- NULL for personal vaults
    owner_user_id UUID REFERENCES users(id), -- Set for personal vaults
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL, -- Optional project link
    vault_type VARCHAR(20) NOT NULL, -- personal, org, project
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_org_vault (org_id),
    INDEX idx_owner (owner_user_id),
    INDEX idx_project (project_id),
    
    CONSTRAINT chk_vault_type CHECK (
        (vault_type = 'personal' AND owner_user_id IS NOT NULL AND org_id IS NULL) OR
        (vault_type = 'org' AND org_id IS NOT NULL AND owner_user_id IS NULL) OR
        (vault_type = 'project' AND project_id IS NOT NULL AND org_id IS NOT NULL)
    )
);
```

### Vault Access (Who can access which vault)
```sql
CREATE TABLE vault_access (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id UUID REFERENCES vaults(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id),
    granted_by UUID REFERENCES users(id),
    granted_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(vault_id, user_id),
    INDEX idx_vault_user (vault_id, user_id)
);
```

### Credentials (Updated for vault-scoped)
```sql
CREATE TABLE credentials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vault_id UUID REFERENCES vaults(id) ON DELETE CASCADE,
    credential_name VARCHAR(255) NOT NULL,
    username VARCHAR(255),
    url TEXT,
    encrypted_data TEXT NOT NULL,
    nonce VARCHAR(255) NOT NULL,
    storage_ref TEXT, -- IPFS or S3 ref for large data
    fabric_commit_hash VARCHAR(255),
    share2 TEXT, -- Fallback DB storage for share2
    tags TEXT[], -- Array of tags for filtering
    credential_type VARCHAR(50) DEFAULT 'password', -- password, api_key, ssh_key, cert, note
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    last_accessed TIMESTAMP,
    deleted_at TIMESTAMP,
    
    INDEX idx_vault_credentials (vault_id),
    INDEX idx_created_by (created_by),
    INDEX idx_tags USING GIN(tags)
);
```

### Key Shards (Reference to shard locations)
```sql
CREATE TABLE key_shards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    credential_id UUID REFERENCES credentials(id) ON DELETE CASCADE,
    shard_index INT NOT NULL, -- 1, 2, or 3
    location VARCHAR(50) NOT NULL, -- vault, fabric, client
    storage_ref TEXT NOT NULL, -- Vault path or Fabric transaction ID
    created_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(credential_id, shard_index),
    INDEX idx_credential_shards (credential_id)
);
```

### Audit Logs (Updated for multi-tenant)
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    vault_id UUID REFERENCES vaults(id) ON DELETE SET NULL,
    credential_id UUID REFERENCES credentials(id) ON DELETE SET NULL,
    action VARCHAR(50) NOT NULL, -- create, read, update, delete, share, rotate, invite, remove_member
    resource_type VARCHAR(50), -- credential, vault, org, member
    metadata JSONB, -- Additional context
    ip_address VARCHAR(255), -- Hashed
    device_fingerprint TEXT,
    timestamp TIMESTAMP DEFAULT NOW(),
    fabric_tx_hash VARCHAR(255),
    
    INDEX idx_org_audit (org_id, timestamp DESC),
    INDEX idx_user_audit (user_id, timestamp DESC),
    INDEX idx_credential_audit (credential_id, timestamp DESC),
    INDEX idx_action (action),
    INDEX idx_timestamp (timestamp DESC)
);
```

### Invitations
```sql
CREATE TABLE invitations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    org_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(255), -- Optional
    wallet_address VARCHAR(42), -- Optional - if inviting known wallet
    role_id UUID REFERENCES roles(id),
    invited_by UUID REFERENCES users(id),
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    accepted_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    INDEX idx_org_invitations (org_id),
    INDEX idx_token (token),
    INDEX idx_wallet (wallet_address)
);
```

---

## Key Relationships

```
User 1:N OrganizationMember N:1 Organization
User 1:N Vault (personal)
Organization 1:N Vault (org/project)
Organization 1:N Project
Project 1:1 Vault (optional)
Vault 1:N Credential
Credential 1:3 KeyShard
Organization 1:N Role
Role 1:N Permission
User N:M Vault (via vault_access)
```

---

## Access Control Logic

### Check if user can access credential:
```sql
SELECT c.* 
FROM credentials c
JOIN vaults v ON c.vault_id = v.id
LEFT JOIN vault_access va ON va.vault_id = v.id AND va.user_id = $userId
LEFT JOIN organization_members om ON om.org_id = v.org_id AND om.user_id = $userId
LEFT JOIN permissions p ON p.role_id = om.role_id
WHERE c.id = $credentialId
  AND (
    -- Personal vault owner
    (v.vault_type = 'personal' AND v.owner_user_id = $userId)
    OR
    -- Org member with permission
    (v.vault_type IN ('org', 'project') AND p.resource_type = 'credential' AND p.action = 'read')
    OR
    -- Explicit vault access
    (va.user_id = $userId)
  );
```

---

## Migration Strategy

### Phase 1: Add new tables (non-breaking)
1. Create users, organizations, organization_members
2. Create roles, permissions
3. Create projects, vaults
4. Update credentials to reference vaults

### Phase 2: Data migration
1. Migrate existing credentials to personal vaults
2. Create default user for each wallet_address in credentials
3. Create default "Personal" organization for each user

### Phase 3: Enforce new logic
1. Update API handlers to check vault access
2. Add organization context to all queries
3. Implement RBAC middleware

---

**This schema supports:**
- ✅ Multi-tenant organizations
- ✅ Flexible RBAC
- ✅ Personal + Team vaults
- ✅ Project workspaces
- ✅ Audit trail with org context
- ✅ Future: Per-org clusters

