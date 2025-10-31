'use client';

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useAccount } from 'wagmi';

interface Organization {
  id: string;
  name: string;
  slug: string;
  plan: string;
  role?: {
    id: string;
    name: string;
  };
  memberCount?: number;
  vaultCount?: number;
}

interface OrganizationContextType {
  currentOrg: Organization | null;
  organizations: Organization[];
  isLoading: boolean;
  error: string | null;
  switchOrg: (orgId: string | null) => void;
  refreshOrganizations: () => Promise<void>;
  isPersonalMode: boolean;
}

const OrganizationContext = createContext<OrganizationContextType | undefined>(undefined);

export function OrganizationProvider({ children }: { children: ReactNode }) {
  const { address, isConnected } = useAccount();
  const [currentOrg, setCurrentOrg] = useState<Organization | null>(null);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Load organizations when wallet connects
  useEffect(() => {
    if (isConnected && address) {
      loadOrganizations();
    } else {
      setOrganizations([]);
      setCurrentOrg(null);
    }
  }, [isConnected, address]);

  // Load user's organizations
  const loadOrganizations = async () => {
    if (!address) return;

    setIsLoading(true);
    setError(null);

    try {
      const response = await fetch('http://localhost:8080/api/v1/me/organizations', {
        headers: {
          'X-Wallet-Address': address,
        },
      });

      if (!response.ok) {
        throw new Error('Failed to load organizations');
      }

      const data = await response.json();
      setOrganizations(data || []);

      // Restore last selected org from localStorage
      const savedOrgId = localStorage.getItem('currentOrgId');
      if (savedOrgId && data?.find((org: Organization) => org.id === savedOrgId)) {
        const org = data.find((org: Organization) => org.id === savedOrgId);
        setCurrentOrg(org);
      }
    } catch (err) {
      console.error('Failed to load organizations:', err);
      setError(err instanceof Error ? err.message : 'Failed to load organizations');
    } finally {
      setIsLoading(false);
    }
  };

  // Switch organization
  const switchOrg = (orgId: string | null) => {
    if (orgId === null) {
      // Switch to personal mode
      setCurrentOrg(null);
      localStorage.removeItem('currentOrgId');
    } else {
      const org = organizations.find(o => o.id === orgId);
      if (org) {
        setCurrentOrg(org);
        localStorage.setItem('currentOrgId', orgId);
      }
    }
  };

  // Refresh organizations list
  const refreshOrganizations = async () => {
    await loadOrganizations();
  };

  const value: OrganizationContextType = {
    currentOrg,
    organizations,
    isLoading,
    error,
    switchOrg,
    refreshOrganizations,
    isPersonalMode: currentOrg === null,
  };

  return (
    <OrganizationContext.Provider value={value}>
      {children}
    </OrganizationContext.Provider>
  );
}

export function useOrganization() {
  const context = useContext(OrganizationContext);
  if (context === undefined) {
    throw new Error('useOrganization must be used within an OrganizationProvider');
  }
  return context;
}

