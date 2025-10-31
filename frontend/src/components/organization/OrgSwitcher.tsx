'use client';

import React from 'react';
import { useAccount } from 'wagmi';
import { useOrganization } from '@/contexts/OrganizationContext';
import { Building, User, ChevronDown, Plus } from 'lucide-react';
import { useRouter } from 'next/navigation';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Button } from '@/components/ui/button';

export function OrgSwitcher() {
  const { address, isConnected } = useAccount();
  const { currentOrg, organizations, switchOrg, isLoading } = useOrganization();
  const router = useRouter();

  if (!isConnected) {
    return null;
  }

  if (isLoading) {
    return (
      <Button variant="ghost" disabled>
        <div className="animate-pulse">Loading...</div>
      </Button>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="gap-2">
          {currentOrg ? (
            <>
              <Building className="h-4 w-4 text-blue-500" />
              <span>{currentOrg.name}</span>
              {currentOrg.role && (
                <span className="ml-2 px-2 py-0.5 text-xs bg-blue-100 text-blue-700 rounded-full">
                  {currentOrg.role.name}
                </span>
              )}
            </>
          ) : (
            <>
              <User className="h-4 w-4 text-purple-500" />
              <span>Personal</span>
            </>
          )}
          <ChevronDown className="h-4 w-4 opacity-50" />
        </Button>
      </DropdownMenuTrigger>

      <DropdownMenuContent align="start" className="w-64">
        {/* Personal Mode */}
        <DropdownMenuItem
          onClick={() => {
            switchOrg(null);
            router.push('/dashboard');
          }}
          className={!currentOrg ? 'bg-purple-50' : ''}
        >
          <User className="mr-2 h-4 w-4 text-purple-500" />
          <div className="flex flex-col">
            <span className="font-medium">Personal Vault</span>
            <span className="text-xs text-gray-500">Your private passwords</span>
          </div>
        </DropdownMenuItem>

        {organizations.length > 0 && <DropdownMenuSeparator />}

        {/* Organizations */}
        {organizations.map((org) => (
          <DropdownMenuItem
            key={org.id}
            onClick={() => {
              switchOrg(org.id);
              router.push(`/organizations/${org.id}`);
            }}
            className={currentOrg?.id === org.id ? 'bg-blue-50' : ''}
          >
            <Building className="mr-2 h-4 w-4 text-blue-500" />
            <div className="flex-1 flex flex-col">
              <span className="font-medium">{org.name}</span>
              <span className="text-xs text-gray-500">
                {org.role?.name} • {org.memberCount || 0} members
              </span>
            </div>
          </DropdownMenuItem>
        ))}

        <DropdownMenuSeparator />

        {/* Create New Organization */}
        <DropdownMenuItem
          onClick={() => router.push('/organizations/new')}
          className="text-blue-600"
        >
          <Plus className="mr-2 h-4 w-4" />
          <span className="font-medium">Create Organization</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}

