# Documentation Locations

## Landing Page
- **Location**: `frontend/src/app/page.tsx`
- **URL**: http://localhost:3000 (when frontend running)
- **Features**: Full marketing landing page with animations, features, how it works

## GitBook Documentation
- **Location**: `docs/`
- **Build**: `cd docs && gitbook serve`
- **URL**: http://localhost:4000
- **Auto-deploy**: GitHub Actions to GitHub Pages

## Structure

```
frontend/
└── src/app/page.tsx          # 🏠 Landing page (main entry)

docs/
├── index.html                # Redirect to GitBook
├── README.md                 # GitBook home
├── SUMMARY.md                # Navigation
├── book.json                 # GitBook config
├── getting-started/
├── architecture/
├── deployment/
├── security/
├── api/
├── development/
└── faq/
```

## View Locally

### Landing Page (Frontend)
```bash
cd frontend
npm run dev
# Visit http://localhost:3000
```

### Documentation (GitBook)
```bash
cd docs
gitbook serve
# Visit http://localhost:4000
```

## Deploy

### Landing Page
Automatically deployed with frontend to Vercel/AWS/GKE

### Documentation
GitHub Actions deploys to GitHub Pages on push

---

**AUUUUFFFF!** 🔥

