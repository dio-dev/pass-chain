'use client';

import React, { useEffect, useState } from 'react';
import { useAccount } from 'wagmi';
import { useParams, useRouter } from 'next/navigation';
import { useOrganization } from '@/contexts/OrganizationContext';
import { 
  Building, 
  Users, 
  Lock, 
  FolderKanban,
  Settings,
  Plus,
  TrendingUp,
  Shield
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

interface OrgStats {
  vaultCount: number;
  credentialCount: number;
  projectCount: number;
}

export default function OrganizationPage() {
  const params = useParams();
  const router = useRouter();
  const { address } = useAccount();
  const { currentOrg, switchOrg } = useOrganization();
  const [stats, setStats] = useState<OrgStats>({ vaultCount: 0, credentialCount: 0, projectCount: 0 });
  const [isLoading, setIsLoading] = useState(true);

  const orgId = params.orgId as string;

  // Ensure we're viewing the correct org
  useEffect(() => {
    if (currentOrg?.id !== orgId) {
      switchOrg(orgId);
    }
  }, [orgId, currentOrg, switchOrg]);

  // Load organization stats
  useEffect(() => {
    const loadStats = async () => {
      if (!address || !orgId) return;

      try {
        const response = await fetch(`http://localhost:8080/api/v1/organizations/${orgId}`, {
          headers: {
            'X-Wallet-Address': address,
          },
        });

        if (response.ok) {
          const data = await response.json();
          setStats({
            vaultCount: data.vaultCount || 1,
            credentialCount: data.credentialCount || 0,
            projectCount: data.projectCount || 0,
          });
        }
      } catch (error) {
        console.error('Failed to load org stats:', error);
      } finally {
        setIsLoading(false);
      }
    };

    loadStats();
  }, [address, orgId]);

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-blue-500 border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <Building className="h-8 w-8 text-blue-500" />
            {currentOrg?.name || 'Organization'}
          </h1>
          <p className="text-gray-600 mt-1">
            Organization overview and quick actions
          </p>
        </div>
        <Button
          variant="outline"
          onClick={() => router.push(`/organizations/${orgId}/settings`)}
        >
          <Settings className="mr-2 h-4 w-4" />
          Settings
        </Button>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Members</CardTitle>
            <Users className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{currentOrg?.memberCount || 1}</div>
            <p className="text-xs text-gray-500 mt-1">Active team members</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Vaults</CardTitle>
            <Lock className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.vaultCount}</div>
            <p className="text-xs text-gray-500 mt-1">Shared vaults</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Projects</CardTitle>
            <FolderKanban className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.projectCount}</div>
            <p className="text-xs text-gray-500 mt-1">Team workspaces</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Credentials</CardTitle>
            <Shield className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{stats.credentialCount}</div>
            <p className="text-xs text-gray-500 mt-1">Stored passwords</p>
          </CardContent>
        </Card>
      </div>

      {/* Quick Actions */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card className="hover:shadow-md transition cursor-pointer" onClick={() => router.push(`/organizations/${orgId}/members`)}>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5 text-blue-500" />
              Invite Team Members
            </CardTitle>
            <CardDescription>
              Add people to your organization and assign roles
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className="w-full">
              <Plus className="mr-2 h-4 w-4" />
              Invite Members
            </Button>
          </CardContent>
        </Card>

        <Card className="hover:shadow-md transition cursor-pointer" onClick={() => router.push(`/organizations/${orgId}/projects`)}>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FolderKanban className="h-5 w-5 text-green-500" />
              Create Project
            </CardTitle>
            <CardDescription>
              Organize credentials into team workspaces
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button className="w-full" variant="outline">
              <Plus className="mr-2 h-4 w-4" />
              New Project
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity / Getting Started */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <TrendingUp className="h-5 w-5" />
            Getting Started
          </CardTitle>
          <CardDescription>
            Complete these steps to set up your organization
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-center gap-3 p-3 rounded-lg border">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-green-100 text-green-600 flex items-center justify-center text-sm font-bold">
              ✓
            </div>
            <div className="flex-1">
              <p className="font-medium">Organization created</p>
              <p className="text-sm text-gray-500">Default vault and roles set up</p>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 rounded-lg border border-blue-200 bg-blue-50">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-sm font-bold">
              2
            </div>
            <div className="flex-1">
              <p className="font-medium">Invite your first team member</p>
              <p className="text-sm text-gray-500">Share access with your team</p>
            </div>
            <Button size="sm" onClick={() => router.push(`/organizations/${orgId}/members`)}>
              Invite
            </Button>
          </div>

          <div className="flex items-center gap-3 p-3 rounded-lg border opacity-60">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-gray-100 text-gray-600 flex items-center justify-center text-sm font-bold">
              3
            </div>
            <div className="flex-1">
              <p className="font-medium">Create your first project</p>
              <p className="text-sm text-gray-500">Organize credentials by workspace</p>
            </div>
          </div>

          <div className="flex items-center gap-3 p-3 rounded-lg border opacity-60">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-gray-100 text-gray-600 flex items-center justify-center text-sm font-bold">
              4
            </div>
            <div className="flex-1">
              <p className="font-medium">Add a shared password</p>
              <p className="text-sm text-gray-500">Start collaborating securely</p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

