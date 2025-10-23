---
hidden: true
---

# Frontend - Pass Chain Web Application

Next.js-based frontend for Pass Chain password management system.

## Tech Stack

* **Next.js 14** - React framework with App Router
* **TypeScript** - Type safety
* **Tailwind CSS** - Styling
* **RainbowKit + Wagmi** - Web3 wallet integration
* **TanStack Query** - Data fetching
* **Zustand** - State management
* **Framer Motion** - Animations
* **Zod** - Validation

## Features

* 🔐 **Wallet Authentication** - Connect with MetaMask, WalletConnect, etc.
* 📝 **Credential Management** - Store and manage passwords
* 🔓 **Secure Reveal** - Decrypt credentials with payment
* 📊 **Usage Dashboard** - View access logs and analytics
* 💰 **Payment Integration** - Pay for storage and access
* 🎨 **Modern UI** - Beautiful, responsive design
* ♿ **Accessible** - WCAG compliant components

## Getting Started

### Prerequisites

* Node.js 18+
* npm or pnpm

### Installation

```bash
cd frontend
pnpm install
```

### Environment Variables

Create `.env.local`:

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_WALLET_CONNECT_PROJECT_ID=your_project_id
```

### Development

```bash
pnpm dev
```

Open http://localhost:3000

### Build

```bash
pnpm build
pnpm start
```

## Project Structure

```
frontend/
├── src/
│   ├── app/              # Next.js App Router pages
│   │   ├── layout.tsx    # Root layout
│   │   ├── page.tsx      # Home page
│   │   └── globals.css   # Global styles
│   ├── components/       # React components
│   │   ├── ui/           # UI components (buttons, cards, etc.)
│   │   └── providers.tsx # Context providers
│   ├── hooks/            # Custom React hooks
│   ├── lib/              # Utility functions
│   ├── store/            # Zustand stores
│   ├── types/            # TypeScript types
│   └── config/           # Configuration files
├── public/               # Static assets
└── package.json
```

## Key Components

### Wallet Connection

Using RainbowKit for wallet connection:

```tsx
import { ConnectButton } from '@rainbow-me/rainbowkit';

<ConnectButton />
```

### API Integration

```typescript
import axios from 'axios';
import { API_BASE_URL } from '@/config/constants';

const api = axios.create({
  baseURL: API_BASE_URL,
});
```

### State Management

```typescript
import { create } from 'zustand';

export const useAuthStore = create((set) => ({
  user: null,
  setUser: (user) => set({ user }),
}));
```

## Available Scripts

* `pnpm dev` - Start development server
* `pnpm build` - Build for production
* `pnpm start` - Start production server
* `pnpm lint` - Run ESLint
* `pnpm type-check` - Run TypeScript checks
* `pnpm test` - Run Jest tests
* `pnpm test:e2e` - Run Playwright E2E tests

## Styling

Using Tailwind CSS with custom design system:

```tsx
<div className="bg-purple-600 hover:bg-purple-700 text-white p-4 rounded-lg">
  Button
</div>
```

## Testing

### Unit Tests

```bash
pnpm test
```

### E2E Tests

```bash
pnpm test:e2e
```

## Deployment

### Vercel (Recommended)

1. Push to GitHub
2. Connect to Vercel
3. Deploy automatically

### Docker

```bash
docker build -t pass-chain-frontend .
docker run -p 3000:3000 pass-chain-frontend
```

## Security

* All encryption happens client-side
* Private keys never leave the browser
* Wallet signatures for authentication
* HTTPS enforced in production

## Browser Support

* Chrome/Edge (Chromium) 90+
* Firefox 90+
* Safari 14+
* Opera 80+

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Write/update tests
5. Submit a pull request

## License

MIT
