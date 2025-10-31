'use client';

import React from 'react';
import { useAccount } from 'wagmi';
import { useOrganization } from '@/contexts/OrganizationContext';
import { usePathname } from 'next/navigation';
import Link from 'next/link';
import { 
  Lock,
  Building,
  FolderKanban,
  Users,
  Plus,
  ChevronRight,
  Home
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';

export function Sidebar() {
  const { isConnected } = useAccount();
  const { currentOrg, isPersonalMode } = useOrganization();
  const pathname = usePathname();

  if (!isConnected) {
    return null;
  }

  return (
    <aside className="hidden md:flex md:flex-col w-64 border-r bg-gray-50/50 min-h-[calc(100vh-4rem)]">
      <div className="p-4 space-y-4">
        {/* Personal Vault Section */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
              Personal
            </h3>
          </div>
          <Link href="/dashboard">
            <Button
              variant={pathname === '/dashboard' && isPersonalMode ? 'secondary' : 'ghost'}
              className="w-full justify-start gap-2"
              size="sm"
            >
              <Lock className="h-4 w-4 text-purple-500" />
              <span>My Vault</span>
            </Button>
          </Link>
        </div>

        {/* Organization Section */}
        {!isPersonalMode && currentOrg && (
          <div>
            <div className="flex items-center justify-between mb-2">
              <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
                Organization
              </h3>
            </div>
            
            <div className="space-y-1">
              {/* Org Overview */}
              <Link href={`/organizations/${currentOrg.id}`}>
                <Button
                  variant={pathname === `/organizations/${currentOrg.id}` ? 'secondary' : 'ghost'}
                  className="w-full justify-start gap-2"
                  size="sm"
                >
                  <Home className="h-4 w-4" />
                  <span>Overview</span>
                </Button>
              </Link>

              {/* Default Vault */}
              <Link href={`/organizations/${currentOrg.id}/vault`}>
                <Button
                  variant={pathname.includes('/vault') ? 'secondary' : 'ghost'}
                  className="w-full justify-start gap-2"
                  size="sm"
                >
                  <Lock className="h-4 w-4 text-blue-500" />
                  <span>Default Vault</span>
                </Button>
              </Link>

              {/* Projects */}
              <Link href={`/organizations/${currentOrg.id}/projects`}>
                <Button
                  variant={pathname.includes('/projects') ? 'secondary' : 'ghost'}
                  className="w-full justify-start gap-2"
                  size="sm"
                >
                  <FolderKanban className="h-4 w-4 text-green-500" />
                  <span>Projects</span>
                </Button>
              </Link>

              {/* Team Members */}
              <Link href={`/organizations/${currentOrg.id}/members`}>
                <Button
                  variant={pathname.includes('/members') ? 'secondary' : 'ghost'}
                  className="w-full justify-start gap-2"
                  size="sm"
                >
                  <Users className="h-4 w-4" />
                  <span>Team</span>
                </Button>
              </Link>
            </div>
          </div>
        )}

        {/* Quick Actions */}
        <div className="pt-4 border-t">
          <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">
            Quick Actions
          </h3>
          <div className="space-y-1">
            {isPersonalMode ? (
              <Link href="/organizations/new">
                <Button
                  variant="outline"
                  className="w-full justify-start gap-2 text-blue-600 border-blue-200 hover:bg-blue-50"
                  size="sm"
                >
                  <Plus className="h-4 w-4" />
                  <span>Create Organization</span>
                </Button>
              </Link>
            ) : (
              <Link href={`/organizations/${currentOrg?.id}/projects/new`}>
                <Button
                  variant="outline"
                  className="w-full justify-start gap-2 text-green-600 border-green-200 hover:bg-green-50"
                  size="sm"
                >
                  <Plus className="h-4 w-4" />
                  <span>New Project</span>
                </Button>
              </Link>
            )}
          </div>
        </div>
      </div>
    </aside>
  );
}

