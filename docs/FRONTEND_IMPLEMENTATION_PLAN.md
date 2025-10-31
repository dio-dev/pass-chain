# 🎨 Frontend Implementation Plan - Enterprise UI

## 🎯 Goal: Build the most intuitive enterprise password manager UI

**Timeline:** Week 3-4 (10 days)
**Tech Stack:** React, Next.js, TailwindCSS, Shadcn/UI

---

## 📦 Phase 1: Foundation (Days 1-2)

### Task 1.1: Update Project Structure
```
frontend/src/
├── app/
│   ├── layout.tsx                    (update)
│   ├── page.tsx                      (landing - existing)
│   ├── dashboard/
│   │   └── page.tsx                  (update for org context)
│   ├── organizations/
│   │   ├── new/page.tsx              (create org)
│   │   └── [orgId]/
│   │       ├── page.tsx              (org dashboard)
│   │       ├── members/page.tsx      (member management)
│   │       ├── projects/page.tsx     (projects list)
│   │       └── settings/page.tsx     (org settings)
│   ├── projects/
│   │   └── [projectId]/page.tsx      (project detail)
│   ├── vaults/
│   │   └── [vaultId]/page.tsx        (vault view)
│   └── invitations/
│       └── [token]/page.tsx          (accept invite)
│
├── components/
│   ├── layout/
│   │   ├── AppShell.tsx              (main layout)
│   │   ├── Header.tsx                (header with context)
│   │   ├── Sidebar.tsx               (navigation)
│   │   └── Breadcrumbs.tsx           (context trail)
│   │
│   ├── organization/
│   │   ├── OrgSwitcher.tsx           (dropdown switcher)
│   │   ├── CreateOrgModal.tsx        (create org)
│   │   ├── OrgCard.tsx               (org display card)
│   │   └── OrgStats.tsx              (stats widget)
│   │
│   ├── members/
│   │   ├── MemberList.tsx            (member table)
│   │   ├── InviteMemberModal.tsx     (invite form)
│   │   ├── MemberCard.tsx            (member display)
│   │   └── RoleBadge.tsx             (role indicator)
│   │
│   ├── projects/
│   │   ├── ProjectCard.tsx           (project card)
│   │   ├── CreateProjectModal.tsx    (create project)
│   │   └── ProjectList.tsx           (projects grid)
│   │
│   ├── vaults/
│   │   ├── VaultSelector.tsx         (vault dropdown)
│   │   ├── VaultCard.tsx             (vault display)
│   │   └── VaultAccessModal.tsx      (grant access)
│   │
│   ├── credentials/
│   │   ├── CredentialList.tsx        (updated with vault)
│   │   ├── AddCredentialModal.tsx    (vault-aware)
│   │   ├── CredentialCard.tsx        (existing)
│   │   └── CredentialDetail.tsx      (with history)
│   │
│   ├── common/
│   │   ├── EmptyState.tsx            (smart empty states)
│   │   ├── LoadingState.tsx          (loading indicators)
│   │   ├── ErrorState.tsx            (error display)
│   │   ├── QuickSearch.tsx           (Cmd+K search)
│   │   ├── Toast.tsx                 (notifications)
│   │   └── PermissionGate.tsx        (hide based on role)
│   │
│   └── ui/
│       └── (existing shadcn components)
│
├── contexts/
│   ├── OrganizationContext.tsx       (org state)
│   ├── VaultContext.tsx              (vault state)
│   └── PermissionsContext.tsx        (user permissions)
│
├── hooks/
│   ├── useOrganization.ts            (org operations)
│   ├── useVaults.ts                  (vault operations)
│   ├── useMembers.ts                 (member operations)
│   ├── useProjects.ts                (project operations)
│   ├── usePermissions.ts             (permission checks)
│   └── useKeyboard.ts                (keyboard shortcuts)
│
└── lib/
    ├── api.ts                        (update with new endpoints)
    └── utils.ts                      (utility functions)
```

---

## 📋 Day-by-Day Plan

### Day 1: Foundation & Context
**Goal:** Set up app shell and organization context

1. **Create OrganizationContext**
   - Store current org, user's orgs, selected vault
   - Provide org switching functionality
   
2. **Update AppShell Layout**
   - Header with OrgSwitcher
   - Sidebar with vault navigation
   - Main content area
   
3. **Create OrgSwitcher Component**
   - Dropdown with user's orgs
   - "Personal" vs org views
   - "+ Create Organization"

**Files to create:**
- `contexts/OrganizationContext.tsx`
- `components/layout/AppShell.tsx`
- `components/layout/Header.tsx`
- `components/layout/Sidebar.tsx`
- `components/organization/OrgSwitcher.tsx`

---

### Day 2: Organization Management
**Goal:** Create & view organizations

1. **Create Organization Flow**
   - CreateOrgModal component
   - Simple form (just name)
   - Success feedback
   
2. **Organization Dashboard**
   - Org overview page
   - Stats (members, vaults, credentials)
   - Quick actions
   
3. **API Integration**
   - Add org endpoints to api.ts
   - Create useOrganization hook

**Files to create:**
- `components/organization/CreateOrgModal.tsx`
- `app/organizations/new/page.tsx`
- `app/organizations/[orgId]/page.tsx`
- `hooks/useOrganization.ts`

---

### Day 3: Member Management
**Goal:** Invite & manage team members

1. **Member List Component**
   - Table with members
   - Role badges
   - Actions (change role, remove)
   
2. **Invite Member Flow**
   - InviteMemberModal
   - Email or wallet input
   - Role selector
   
3. **Accept Invitation Page**
   - Preview invitation
   - Accept with wallet
   - Redirect to org

**Files to create:**
- `components/members/MemberList.tsx`
- `components/members/InviteMemberModal.tsx`
- `components/members/RoleBadge.tsx`
- `app/organizations/[orgId]/members/page.tsx`
- `app/invitations/[token]/page.tsx`
- `hooks/useMembers.ts`

---

### Day 4: Project Management
**Goal:** Create & manage projects

1. **Project Cards Component**
   - Grid of project cards
   - Show vault, credential count
   - Click to open project
   
2. **Create Project Flow**
   - CreateProjectModal
   - Auto-creates vault
   - Success feedback
   
3. **Project Detail Page**
   - Project info
   - Project vault
   - Credentials list
   - Share with team

**Files to create:**
- `components/projects/ProjectCard.tsx`
- `components/projects/CreateProjectModal.tsx`
- `app/organizations/[orgId]/projects/page.tsx`
- `app/projects/[projectId]/page.tsx`
- `hooks/useProjects.ts`

---

### Day 5: Vault Management
**Goal:** Vault hierarchy & access control

1. **VaultSelector Component**
   - Dropdown with all accessible vaults
   - Group by type (Personal, Org, Project)
   - Show access level
   
2. **Vault Detail Page**
   - Vault info
   - Credentials list
   - Grant/revoke access
   
3. **Vault Access Modal**
   - Select users to share with
   - Optional role assignment
   - Success feedback

**Files to create:**
- `components/vaults/VaultSelector.tsx`
- `components/vaults/VaultAccessModal.tsx`
- `app/vaults/[vaultId]/page.tsx`
- `hooks/useVaults.ts`

---

### Day 6: Update Credential Flow
**Goal:** Vault-aware credential management

1. **Update AddCredentialModal**
   - Show current vault context
   - Allow vault selection
   - Smart defaults
   
2. **Update CredentialList**
   - Filter by vault
   - Show vault badge
   - Vault-specific empty states
   
3. **Credential Detail**
   - Show which vault
   - Show access history
   - Show who has access

**Files to update:**
- `components/credentials/AddCredentialModal.tsx`
- `components/credentials/CredentialList.tsx`
- `components/credentials/CredentialDetail.tsx`

---

### Day 7: Smart Empty States
**Goal:** Guide users with empty states

1. **Personal Vault Empty State**
   ```
   🔐 No passwords yet!
   Let's add your first password.
   [+ Add Your First Password]
   ```

2. **Organization Empty State**
   ```
   👥 Ready to collaborate?
   Create your first organization.
   [+ Create Organization]
   ```

3. **Project Empty State**
   ```
   📁 No projects yet!
   Create workspaces for your team.
   [+ Create Project]
   ```

**Files to create:**
- `components/common/EmptyState.tsx` (with variants)

---

### Day 8: Quick Search (Cmd+K)
**Goal:** Fast fuzzy search across vaults

1. **QuickSearch Component**
   - Modal triggered by Cmd+K
   - Fuzzy search all credentials
   - Show vault context
   - Keyboard navigation
   
2. **Search Logic**
   - Search across all accessible vaults
   - Match name, URL, tags
   - Group results by vault

**Files to create:**
- `components/common/QuickSearch.tsx`
- `hooks/useKeyboard.ts`

---

### Day 9: Polish & Feedback
**Goal:** Smooth UX with feedback

1. **Toast Notifications**
   - Success messages
   - Error handling
   - Loading states
   
2. **Loading States**
   - Skeleton screens
   - Optimistic UI
   - Progress indicators
   
3. **Animations**
   - Page transitions
   - Modal animations
   - Hover effects

**Files to create:**
- `components/common/Toast.tsx`
- `components/common/LoadingState.tsx`

---

### Day 10: Permission Gates & Final Polish
**Goal:** Role-based UI hiding

1. **PermissionGate Component**
   - Hide UI based on role
   - Disable buttons
   - Show tooltips
   
2. **Permissions Context**
   - Load user permissions
   - Check permission helper
   - Cache permissions
   
3. **Final Testing**
   - Test all flows
   - Fix bugs
   - Polish UI

**Files to create:**
- `components/common/PermissionGate.tsx`
- `contexts/PermissionsContext.tsx`
- `hooks/usePermissions.ts`

---

## 🎨 Component Examples

### Example 1: OrgSwitcher
```typescript
export function OrgSwitcher() {
  const { currentOrg, organizations, switchOrg } = useOrganization();
  const router = useRouter();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost">
          {currentOrg ? (
            <>
              <Building className="mr-2 h-4 w-4" />
              {currentOrg.name}
            </>
          ) : (
            <>
              <User className="mr-2 h-4 w-4" />
              Personal
            </>
          )}
          <ChevronDown className="ml-2 h-4 w-4" />
        </Button>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent>
        <DropdownMenuItem onClick={() => switchOrg(null)}>
          <User className="mr-2 h-4 w-4" />
          Personal Vault
        </DropdownMenuItem>
        
        <DropdownMenuSeparator />
        
        {organizations.map(org => (
          <DropdownMenuItem 
            key={org.id}
            onClick={() => switchOrg(org.id)}
          >
            <Building className="mr-2 h-4 w-4" />
            {org.name}
            <RoleBadge role={org.role} className="ml-auto" />
          </DropdownMenuItem>
        ))}
        
        <DropdownMenuSeparator />
        
        <DropdownMenuItem onClick={() => router.push('/organizations/new')}>
          <Plus className="mr-2 h-4 w-4" />
          Create Organization
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
```

### Example 2: Empty State
```typescript
export function EmptyState({ 
  variant = 'credentials',
  onAction 
}: EmptyStateProps) {
  const variants = {
    credentials: {
      icon: <Lock className="h-16 w-16 text-purple-400" />,
      title: "No passwords yet!",
      description: "This is your personal vault. Let's add your first password.",
      action: "Add Your First Password",
      tip: "💡 Tip: You can organize passwords by creating projects later."
    },
    organization: {
      icon: <Building className="h-16 w-16 text-blue-400" />,
      title: "Ready to collaborate?",
      description: "Create your first organization to work with your team.",
      action: "Create Organization",
      tip: "✨ Organizations let you share passwords securely with your team."
    },
    // ... more variants
  };

  const config = variants[variant];

  return (
    <div className="flex flex-col items-center justify-center py-16 px-4 text-center">
      {config.icon}
      <h3 className="mt-4 text-2xl font-bold">{config.title}</h3>
      <p className="mt-2 text-gray-600 max-w-md">{config.description}</p>
      
      <Button 
        size="lg" 
        onClick={onAction}
        className="mt-6"
      >
        <Plus className="mr-2 h-5 w-5" />
        {config.action}
      </Button>
      
      <p className="mt-6 text-sm text-gray-500">{config.tip}</p>
    </div>
  );
}
```

### Example 3: PermissionGate
```typescript
export function PermissionGate({ 
  action, 
  resource,
  children,
  fallback = null 
}: PermissionGateProps) {
  const { hasPermission } = usePermissions();
  const canAccess = hasPermission(resource, action);

  if (!canAccess) {
    return fallback;
  }

  return <>{children}</>;
}

// Usage:
<PermissionGate action="invite" resource="member">
  <Button onClick={handleInvite}>
    Invite Member
  </Button>
</PermissionGate>

// Or with disabled state:
<PermissionGate 
  action="invite" 
  resource="member"
  fallback={
    <Tooltip content="You need Admin role to invite members">
      <Button disabled>Invite Member</Button>
    </Tooltip>
  }
>
  <Button onClick={handleInvite}>Invite Member</Button>
</PermissionGate>
```

---

## 🎯 Success Metrics

### User Flow Times:
- Connect wallet → See vault: **< 3 seconds**
- Add first password: **< 30 seconds**
- Create organization: **< 15 seconds**
- Invite team member: **< 10 seconds**
- Create project: **< 10 seconds**
- Find password (Cmd+K): **< 3 seconds**

### Code Quality:
- TypeScript strict mode: **100%**
- Component reusability: **High**
- Performance (Lighthouse): **> 90**
- Accessibility (a11y): **WCAG AA**

---

## 📦 Dependencies to Add

```bash
npm install @radix-ui/react-dropdown-menu
npm install @radix-ui/react-dialog
npm install @radix-ui/react-toast
npm install cmdk  # For Cmd+K search
npm install fuse.js  # Fuzzy search
npm install react-hot-toast  # Toast notifications
```

---

## 🚀 Let's Build!

Ready to implement the most intuitive password manager UI ever built! 🎨

**AUUUUFFFF!** 🔥

