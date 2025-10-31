'use client';

import React, { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { ConnectButton } from '@/components/ConnectButton';
import { OrgSwitcher } from '@/components/organization/OrgSwitcher';
import { QuickSearch } from '@/components/search/QuickSearch';
import { useOrganization } from '@/contexts/OrganizationContext';
import { 
  Search, 
  Bell, 
  Settings,
  Lock,
  Building,
  FolderKanban,
  Users,
  Shield
} from 'lucide-react';
import { Button } from '@/components/ui/button';

export function Header() {
  const { address, isConnected } = useAccount();
  const pathname = usePathname();
  const { currentOrg, isPersonalMode } = useOrganization();
  const [showSearch, setShowSearch] = useState(false);

  // Global keyboard shortcut for search (Cmd+K / Ctrl+K)
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault();
        setShowSearch(true);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/60">
      <div className="container flex h-16 items-center px-4">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2 mr-6">
          <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-purple-600 to-blue-600 flex items-center justify-center">
            <Lock className="h-5 w-5 text-white" />
          </div>
          <span className="font-bold text-xl hidden sm:inline">Pass Chain</span>
        </Link>

        {/* Organization Switcher */}
        {isConnected && (
          <div className="mr-6">
            <OrgSwitcher />
          </div>
        )}

        {/* Navigation */}
        {isConnected && (
          <nav className="hidden md:flex items-center gap-1 flex-1">
            <Link href="/dashboard">
              <Button
                variant={pathname === '/dashboard' ? 'secondary' : 'ghost'}
                size="sm"
                className="gap-2"
              >
                <Lock className="h-4 w-4" />
                Passwords
              </Button>
            </Link>

            {!isPersonalMode && currentOrg && (
              <>
                <Link href={`/organizations/${currentOrg.id}/projects`}>
                  <Button
                    variant={pathname.includes('/projects') ? 'secondary' : 'ghost'}
                    size="sm"
                    className="gap-2"
                  >
                    <FolderKanban className="h-4 w-4" />
                    Projects
                  </Button>
                </Link>

                <Link href={`/organizations/${currentOrg.id}/members`}>
                  <Button
                    variant={pathname.includes('/members') ? 'secondary' : 'ghost'}
                    size="sm"
                    className="gap-2"
                  >
                    <Users className="h-4 w-4" />
                    Team
                  </Button>
                </Link>

                <Link href={`/organizations/${currentOrg.id}/audit`}>
                  <Button
                    variant={pathname.includes('/audit') ? 'secondary' : 'ghost'}
                    size="sm"
                    className="gap-2"
                  >
                    <Shield className="h-4 w-4" />
                    Audit
                  </Button>
                </Link>
              </>
            )}
          </nav>
        )}

        {/* Right Side Actions */}
        <div className="flex items-center gap-2 ml-auto">
          {isConnected && (
            <>
              {/* Quick Search - Cmd+K */}
              <Button
                variant="outline"
                size="sm"
                className="hidden sm:flex gap-2"
                onClick={() => setShowSearch(true)}
              >
                <Search className="h-4 w-4" />
                <span className="hidden lg:inline">Search</span>
                <kbd className="hidden lg:inline pointer-events-none h-5 select-none items-center gap-1 rounded border bg-muted px-1.5 font-mono text-[10px] font-medium opacity-100">
                  <span className="text-xs">⌘</span>K
                </kbd>
              </Button>

              {/* Notifications */}
              <Button variant="ghost" size="icon">
                <Bell className="h-5 w-5" />
              </Button>

              {/* Settings */}
              <Link href="/settings">
                <Button variant="ghost" size="icon">
                  <Settings className="h-5 w-5" />
                </Button>
              </Link>
            </>
          )}

          {/* Connect Wallet Button */}
          <ConnectButton />
        </div>
      </div>

      {/* Quick Search Modal */}
      <QuickSearch open={showSearch} onOpenChange={setShowSearch} />
    </header>
  );
}

