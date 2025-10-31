# Contributing to Pass Chain

Thank you for your interest in contributing to Pass Chain! 🔐

## 🚀 Getting Started

1. **Fork the repository**
2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/pass-chain.git
   cd pass-chain
   ```
3. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/ORIGINAL_OWNER/pass-chain.git
   ```
4. **Set up development environment** (see README.md)

## 📋 Development Process

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test improvements

### 2. Make Your Changes

Follow our coding standards:
- **Go**: Follow [Effective Go](https://golang.org/doc/effective_go.html)
- **TypeScript**: Follow project ESLint config
- **Commits**: Use [Conventional Commits](https://www.conventionalcommits.org/)

### 3. Write Tests

**All new features must include tests!**

```bash
# Unit tests for new handlers/services
# Integration tests for API endpoints
# E2E tests for complete flows

# Run tests locally
make test-all
```

### 4. Run Linters

```bash
# Backend
cd backend
golangci-lint run

# Frontend
cd frontend
npm run lint
```

### 5. Commit Your Changes

Use conventional commit format:

```bash
git commit -m "feat: add user invitation flow"
git commit -m "fix: resolve UUID generation bug in org service"
git commit -m "docs: update API documentation"
git commit -m "test: add E2E tests for RBAC"
```

Commit types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

### 6. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## ✅ Pull Request Guidelines

### PR Title

Use conventional commit format:
```
feat: add organization invitation system
fix: resolve credential decryption error
docs: update testing guide
```

### PR Description Template

```markdown
## Description
Brief description of what this PR does

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] E2E tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] Documentation updated
- [ ] No new warnings
- [ ] Commit messages follow conventions

## Related Issues
Closes #123
```

### Code Review Process

1. **Automated Checks**
   - All tests must pass
   - Linting must pass
   - Coverage must not decrease
   - Security scan must pass

2. **Manual Review**
   - At least 1 approving review required
   - Reviewers check code quality, logic, tests

3. **Merge**
   - Squash and merge (for clean history)
   - Delete branch after merge

## 🧪 Testing Requirements

### Backend (Go)

**Required for all PRs:**
```bash
# Must pass
make test-all

# Coverage should not decrease
make test-coverage
```

**Test Structure:**
- `internal/api/handlers/*_test.go` - Handler unit tests
- `internal/services/*_test.go` - Service unit tests
- `test/integration_test.go` - API integration tests
- `test/e2e_flow_test.go` - End-to-end flows

**Example Test:**
```go
func TestCreateOrganization(t *testing.T) {
    // Setup
    db, cleanup := setupTestDB(t)
    defer cleanup()
    
    // Test
    org, err := service.CreateOrganization(ctx, "Test Org", userID)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, org.ID)
}
```

### Frontend (TypeScript)

**Required:**
```bash
npm run lint
npm run build  # Must build without errors
```

## 📝 Documentation

Update documentation for:
- **New features**: Add to README.md and relevant docs/
- **API changes**: Update API documentation
- **Breaking changes**: Add migration guide
- **Configuration**: Update example configs

## 🐛 Reporting Bugs

Use GitHub Issues with this template:

```markdown
## Bug Description
Clear description of the bug

## Steps to Reproduce
1. Step 1
2. Step 2
3. See error

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g. Windows 11, macOS 13, Ubuntu 22.04]
- Browser: [e.g. Chrome 120]
- Go version: [e.g. 1.23]
- Node version: [e.g. 18.17]

## Logs/Screenshots
```

## 💡 Requesting Features

Use GitHub Issues with this template:

```markdown
## Feature Description
What feature would you like?

## Use Case
Why is this feature needed?

## Proposed Solution
How would you implement it?

## Alternatives Considered
Other approaches you've thought about
```

## 🏗️ Architecture Guidelines

### Backend Principles

1. **Clean Architecture**
   - `handlers/` - HTTP request handling only
   - `services/` - Business logic
   - `models/` - Data structures
   - `database/` - DB operations

2. **Error Handling**
   ```go
   if err != nil {
       logger.Error("descriptive message", "context", value)
       return fmt.Errorf("user-facing message: %w", err)
   }
   ```

3. **Logging**
   ```go
   logger.Info("action completed", "userId", user.ID, "result", result)
   ```

### Frontend Principles

1. **Component Structure**
   - Small, focused components
   - Reusable UI components in `components/ui/`
   - Feature components in `components/[feature]/`

2. **API Calls**
   - Use centralized `lib/api.ts`
   - Handle errors gracefully
   - Show loading states

3. **State Management**
   - Use React hooks
   - Minimal global state
   - Keep data close to where it's used

## 🔒 Security

### Reporting Security Issues

**Do NOT create public GitHub issues for security vulnerabilities!**

Instead, email: security@pass-chain.com

Include:
- Description of vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Security Guidelines

1. **Never log sensitive data**
   - No passwords, keys, or credentials in logs
   - Hash IP addresses in audit logs

2. **Always validate input**
   - Sanitize user input
   - Validate wallet signatures
   - Check permissions

3. **Use secure defaults**
   - Fail closed, not open
   - Require authentication
   - Minimal permissions

## 📚 Resources

- [Go Style Guide](https://google.github.io/styleguide/go/)
- [TypeScript Style Guide](https://google.github.io/styleguide/tsguide.html)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [React Best Practices](https://react.dev/learn)
- [Testing Best Practices](backend/test/README.md)

## 🎯 Good First Issues

Look for issues labeled `good first issue` - these are beginner-friendly!

## 💬 Questions?

- GitHub Discussions: Ask questions
- Discord: Join our community
- Email: contribute@pass-chain.com

---

Thank you for contributing! **AUUUUFFFF! 🔥**

