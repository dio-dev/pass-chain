# 🚀 Pass Chain Enterprise - Implementation Plan

## 📋 Current Status

### ✅ Completed (Phase 1 - MVP):
- Wallet-based authentication
- Personal password manager
- Client-side encryption (XChaCha20-Poly1305)
- Shamir Secret Sharing (3 shards)
- HashiCorp Vault integration
- PostgreSQL storage
- Basic audit logs
- Frontend dashboard
- Kubernetes deployment
- Smart contracts (ERC-20 token, ERC-721 NFT)

### ✅ Completed (Phase 2 - Architecture):
- **Database schema** for multi-tenant orgs (`docs/ENTERPRISE_SCHEMA.md`)
- **Go models** for all entities (`backend/internal/models/models_enterprise.go`)
- **RBAC service** with 6 default roles (`backend/internal/services/rbac_service.go`)
- **Organization service** (`backend/internal/services/organization_service.go`)
- **User stories** (24 stories across 10 epics) (`docs/USER_STORIES.md`)
- **Audit + Vault integration** design (`docs/AUDIT_VAULT_INTEGRATION.md`)

---

## 🎯 Phase 3: Backend Implementation (Week 1-2)

### Goal: Implement all organization, vault, and project APIs

### Task 1: Database Migration
**File:** `backend/internal/database/migrations/001_enterprise.go`

**Steps:**
1. Create migration file
2. Add AutoMigrate for new models:
   - User, Organization, OrganizationMember
   - Role, Permission
   - Project, Vault, VaultAccess
   - Invitation, KeyShard
3. Create migration script to backfill existing data:
   - Create User for each existing wallet
   - Create personal Vault for each User
   - Link existing Credentials to personal Vaults

**Code:**
```go
// Run migrations
db.AutoMigrate(
    &models.User{},
    &models.Organization{},
    &models.OrganizationMember{},
    &models.Role{},
    &models.Permission{},
    &models.Project{},
    &models.Vault{},
    &models.VaultAccess{},
    &models.KeyShard{},
    &models.Invitation{},
)

// Backfill existing credentials
migrateExistingCredentials(db)
```

**User Story:** Foundation for all other stories

---

### Task 2: User/Auth Handlers
**File:** `backend/internal/api/handlers/user_handler.go`

**Endpoints:**
- `GET /api/v1/me` - Get current user profile
- `PUT /api/v1/me` - Update profile (display name)
- `GET /api/v1/me/organizations` - List user's orgs

**User Story:** US-1.2 (View My Organizations)

---

### Task 3: Organization Handlers
**File:** `backend/internal/api/handlers/organization_handler.go`

**Endpoints:**
- `POST /api/v1/organizations` - Create org
- `GET /api/v1/organizations` - List user's orgs
- `GET /api/v1/organizations/:id` - Get org details
- `PUT /api/v1/organizations/:id` - Update org
- `DELETE /api/v1/organizations/:id` - Delete org

**User Stories:** US-1.1 (Create Organization), US-1.2 (View Orgs)

**Example:**
```go
func (h *OrganizationHandler) CreateOrganization(c *fiber.Ctx) error {
    var req struct {
        Name string `json:"name" validate:"required"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    userID := c.Locals("userId").(string)
    org, err := h.orgService.CreateOrganization(c.Context(), req.Name, userID)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(201).JSON(org)
}
```

---

### Task 4: Member Management Handlers
**File:** `backend/internal/api/handlers/member_handler.go`

**Endpoints:**
- `GET /api/v1/organizations/:id/members` - List members
- `POST /api/v1/organizations/:id/invite` - Invite member
- `POST /api/v1/invitations/:token/accept` - Accept invite
- `GET /api/v1/invitations/:token` - Preview invite
- `PUT /api/v1/organizations/:id/members/:userId/role` - Update role
- `DELETE /api/v1/organizations/:id/members/:userId` - Remove member

**User Stories:** US-2.1, US-2.2, US-2.3, US-2.4

---

### Task 5: Project Handlers
**File:** `backend/internal/api/handlers/project_handler.go`

**Endpoints:**
- `POST /api/v1/organizations/:id/projects` - Create project
- `GET /api/v1/organizations/:id/projects` - List projects
- `GET /api/v1/projects/:id` - Get project
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project

**User Stories:** US-3.1 (Create Project)

---

### Task 6: Vault Handlers
**File:** `backend/internal/api/handlers/vault_handler.go`

**Endpoints:**
- `POST /api/v1/vaults` - Create vault
- `GET /api/v1/vaults` - List accessible vaults
- `GET /api/v1/vaults/:id` - Get vault
- `PUT /api/v1/vaults/:id` - Update vault
- `DELETE /api/v1/vaults/:id` - Delete vault
- `POST /api/v1/vaults/:id/access` - Grant access
- `DELETE /api/v1/vaults/:id/access/:userId` - Revoke access
- `GET /api/v1/vaults/:id/audit` - Vault audit logs

**User Stories:** US-3.2, US-3.3, US-5.3

---

### Task 7: Update Credential Handlers
**File:** `backend/internal/api/handlers/credentials.go` (update existing)

**Changes:**
1. Add `vaultId` to create/update requests
2. Check vault access permissions (via RBAC)
3. Filter credentials by vault
4. Update audit logging to include vault context

**Updated Endpoints:**
- `POST /api/v1/credentials` - Now requires `vaultId`
- `GET /api/v1/credentials` - Add `?vaultId=xxx` filter
- `GET /api/v1/credentials/:id` - Check vault access

**User Stories:** US-4.1, US-4.2, US-4.3

---

### Task 8: Enhanced Audit Handlers
**File:** `backend/internal/api/handlers/audit_handler.go`

**Endpoints:**
- `GET /api/v1/organizations/:id/audit` - Org audit logs
- `GET /api/v1/vaults/:id/audit` - Vault audit logs
- `GET /api/v1/credentials/:id/audit` - Credential history (existing)
- `GET /api/v1/audit/organizations/:orgId/vault` - Read from Vault backup
- `GET /api/v1/audit/credentials/:id/verify` - Verify integrity

**User Stories:** US-5.1, US-5.3, US-5.4

---

### Task 9: Role Management Handlers
**File:** `backend/internal/api/handlers/role_handler.go`

**Endpoints:**
- `GET /api/v1/organizations/:id/roles` - List roles
- `POST /api/v1/organizations/:id/roles` - Create custom role
- `GET /api/v1/roles/:id` - Get role details
- `PUT /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role
- `GET /api/v1/me/permissions?orgId=xxx` - My permissions

**User Stories:** US-7.1, US-7.2

---

## 🎨 Phase 4: Frontend Implementation (Week 3-4)

### Goal: Build organization management UI

### Task 10: Organization Context & Switcher
**Files:**
- `frontend/src/contexts/OrganizationContext.tsx`
- `frontend/src/components/OrgSwitcher.tsx`

**Features:**
- Global org context (selected org, user's orgs)
- Org switcher dropdown in header
- "Personal" vs org views

**User Story:** US-1.2

---

### Task 11: Organization Dashboard
**File:** `frontend/src/app/organizations/[orgId]/page.tsx`

**Features:**
- Org overview (name, plan, member count)
- Quick stats (vaults, credentials, projects)
- Recent activity
- Member list preview

**User Story:** US-1.1, US-1.2

---

### Task 12: Member Management Page
**File:** `frontend/src/app/organizations/[orgId]/members/page.tsx`

**Features:**
- Member list table (name, wallet, role, status)
- Invite button → InviteModal
- Role dropdown (change role)
- Remove member action
- Filter by role/status

**User Stories:** US-2.1, US-2.3, US-2.4

---

### Task 13: Invite System
**Files:**
- `frontend/src/components/InviteModal.tsx`
- `frontend/src/app/invitations/[token]/page.tsx`

**Features:**
- Invite modal (enter email or wallet, select role)
- Generate invite link
- Invite acceptance page (preview org, accept/decline)

**User Stories:** US-2.1, US-2.2

---

### Task 14: Project Management
**File:** `frontend/src/app/organizations/[orgId]/projects/page.tsx`

**Features:**
- Project cards grid
- Create project button
- Project detail view
- Link to project vault

**User Story:** US-3.1

---

### Task 15: Vault Selector
**File:** `frontend/src/components/VaultSelector.tsx`

**Features:**
- Dropdown/sidebar with vault list
- Group by type (Personal, Org, Project)
- Filter credentials by selected vault
- Vault access indicator (read/write)

**User Stories:** US-3.2, US-4.3

---

### Task 16: Update Credential Forms
**Files:**
- `frontend/src/app/dashboard/page.tsx` (update)
- `frontend/src/components/AddCredentialModal.tsx` (new)

**Changes:**
- Add vault selector to "Add Credential" form
- Default to current selected vault
- Show vault name in credential list

**User Stories:** US-4.1, US-4.2

---

### Task 17: Audit Dashboard
**File:** `frontend/src/app/organizations/[orgId]/audit/page.tsx`

**Features:**
- Audit log table (timestamp, user, action, resource)
- Filters (user, action, date range, resource type)
- Export to CSV button
- "Verify Integrity" button (for credentials)
- Toggle between DB logs and Vault logs

**User Stories:** US-5.1, US-5.3, US-5.4

---

### Task 18: Role Management UI
**File:** `frontend/src/app/organizations/[orgId]/roles/page.tsx`

**Features:**
- Role list (default + custom)
- Create custom role button
- Permission checklist (resource type × action)
- "My Permissions" view

**User Stories:** US-7.1, US-7.2

---

### Task 19: Permission Gates
**File:** `frontend/src/components/PermissionGate.tsx`

**Features:**
- Component that hides UI based on permissions
- Usage: `<PermissionGate action="member:invite">...</PermissionGate>`
- Disable buttons for unauthorized actions

**User Stories:** All (security)

---

## 🧪 Phase 5: Testing & Integration (Week 5)

### Task 20: Backend Unit Tests
**Files:** `backend/internal/services/*_test.go`

**Coverage:**
- RBAC service (permission checks)
- Organization service (create, invite, etc.)
- Audit service (DB + Vault writes)

---

### Task 21: Integration Tests
**Files:** `backend/tests/integration_test.go`

**Scenarios:**
- Create org → Invite member → Accept → Access vault
- Create project → Add credential → Share → Retrieve
- Audit log end-to-end (action → DB → Vault → verify)

---

### Task 22: Frontend E2E Tests (Optional)
**Files:** `frontend/e2e/*.spec.ts`

**Scenarios:**
- User creates org
- User invites team member
- Member accepts and accesses shared credential

---

## 📦 Phase 6: Deployment & Documentation (Week 6)

### Task 23: Update Kubernetes Manifests
**Files:** `infrastructure/k8s/backend/deployment.yaml`

**Changes:**
- Add environment variables for new features
- Ensure database has new tables
- Update health checks

---

### Task 24: Update Helm Chart
**Files:** `infrastructure/k8s/helm/passchain/values.yaml`

**Add:**
- Enterprise feature flags
- Audit Vault integration settings

---

### Task 25: Documentation Updates

**Files to Create/Update:**
- `docs/GETTING_STARTED.md` - Add org setup
- `docs/API.md` - Document all new endpoints
- `docs/RBAC.md` - Explain roles & permissions
- `README.md` - Update with enterprise features

---

## 📊 Implementation Timeline

### Week 1: Backend Core (Tasks 1-5)
- **Mon-Tue:** Database migration + User handlers
- **Wed-Thu:** Organization handlers + Member handlers
- **Fri:** Project handlers

### Week 2: Backend Advanced (Tasks 6-9)
- **Mon-Tue:** Vault handlers + Update credentials
- **Wed-Thu:** Enhanced audit + Vault integration
- **Fri:** Role management handlers + testing

### Week 3: Frontend Foundation (Tasks 10-13)
- **Mon-Tue:** Org context + switcher + dashboard
- **Wed-Thu:** Member management + Invite system
- **Fri:** Testing & polish

### Week 4: Frontend Features (Tasks 14-19)
- **Mon-Tue:** Projects + Vault selector
- **Wed:** Update credential forms
- **Thu:** Audit dashboard
- **Fri:** Role management + Permission gates

### Week 5: Testing & Integration (Tasks 20-22)
- **Mon-Wed:** Unit tests + Integration tests
- **Thu-Fri:** Bug fixes + E2E testing

### Week 6: Deployment & Docs (Tasks 23-25)
- **Mon-Tue:** K8s updates + Helm chart
- **Wed-Thu:** Documentation
- **Fri:** Production deploy + validation

---

## 🎯 Success Criteria

### Phase 3 (Backend) Complete When:
- ✅ All 25+ API endpoints implemented
- ✅ RBAC enforced on all endpoints
- ✅ Database migration tested
- ✅ Audit logs write to DB + Vault
- ✅ Unit tests >80% coverage

### Phase 4 (Frontend) Complete When:
- ✅ Users can create orgs
- ✅ Users can invite/manage members
- ✅ Users can create projects
- ✅ Users can create/access vaults
- ✅ Credentials are vault-scoped
- ✅ Audit dashboard functional
- ✅ Permission gates working

### Phase 5 (Testing) Complete When:
- ✅ All integration tests pass
- ✅ E2E flows validated
- ✅ No critical bugs

### Phase 6 (Deploy) Complete When:
- ✅ Deployed to staging
- ✅ Documentation complete
- ✅ Production-ready

---

## 🚀 Quick Start (Development)

### Backend Setup:
```bash
cd backend

# Run migrations
go run cmd/server/main.go migrate

# Start server
go run cmd/server/main.go
```

### Frontend Setup:
```bash
cd frontend

# Start dev server
npm run dev
```

### Test Organization Flow:
1. Connect wallet
2. Create organization "My Startup"
3. Invite member (email or wallet)
4. Create project "Mobile App"
5. Add credential to project vault
6. Share with team member
7. View audit logs

---

## 📞 Next Steps

**Immediate:**
1. ✅ Start with Task 1 (Database migration)
2. ✅ Implement Tasks 2-3 (User + Org handlers)
3. ✅ Test org creation flow

**This Week:**
- Complete Tasks 1-5 (Backend core)
- Test API endpoints with Postman/curl

**Next Week:**
- Complete Tasks 6-9 (Backend advanced)
- Start frontend implementation

---

**AUUUUFFFF!** 🔥

Let's build this enterprise platform! 🚀

