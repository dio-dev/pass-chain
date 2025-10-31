'use client';

import React, { useEffect, useState } from 'react';
import { useAccount } from 'wagmi';
import { useParams, useRouter } from 'next/navigation';
import { useOrganization } from '@/contexts/OrganizationContext';
import {
  Users,
  Mail,
  MoreVertical,
  Trash2,
  Shield,
  Clock,
  UserCheck,
  UserX,
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { RoleBadge } from '@/components/organization/RoleBadge';
import { InviteMemberModal } from '@/components/organization/InviteMemberModal';
import { toast } from 'sonner';

interface Member {
  id: string;
  user: {
    id: string;
    walletAddress: string;
    displayName?: string;
    ensName?: string;
  };
  role: {
    id: string;
    name: string;
    description: string;
  };
  status: string;
  joinedAt?: string;
  createdAt: string;
}

interface Role {
  id: string;
  name: string;
  description: string;
}

export default function MembersPage() {
  const params = useParams();
  const router = useRouter();
  const { address } = useAccount();
  const { currentOrg } = useOrganization();
  const [members, setMembers] = useState<Member[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showInviteModal, setShowInviteModal] = useState(false);

  const orgId = params.orgId as string;

  useEffect(() => {
    if (address && orgId) {
      loadMembers();
      loadRoles();
    }
  }, [address, orgId]);

  const loadMembers = async () => {
    if (!address) return;

    setIsLoading(true);
    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${orgId}/members`,
        {
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setMembers(data.members || []);
      }
    } catch (error) {
      console.error('Failed to load members:', error);
      toast.error('Failed to load team members');
    } finally {
      setIsLoading(false);
    }
  };

  const loadRoles = async () => {
    if (!address) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${orgId}`,
        {
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        // Roles will be in the org response
        // For now, we'll use default roles
        setRoles([
          { id: '1', name: 'Owner', description: 'Full control over organization' },
          { id: '2', name: 'Admin', description: 'Manage members and settings' },
          { id: '3', name: 'Security Officer', description: 'Audit and security' },
          { id: '4', name: 'Member', description: 'Access shared vaults' },
          { id: '5', name: 'Auditor', description: 'View-only access' },
          { id: '6', name: 'Guest', description: 'Limited access' },
        ]);
      }
    } catch (error) {
      console.error('Failed to load roles:', error);
    }
  };

  const handleRemoveMember = async (memberId: string, userName: string) => {
    if (!address) return;

    if (!confirm(`Remove ${userName} from the organization?`)) {
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${orgId}/members/${memberId}`,
        {
          method: 'DELETE',
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (response.ok) {
        toast.success(`${userName} removed from organization`);
        loadMembers();
      } else {
        throw new Error('Failed to remove member');
      }
    } catch (error) {
      console.error('Failed to remove member:', error);
      toast.error('Failed to remove member');
    }
  };

  const handleChangeRole = async (memberId: string, newRoleId: string) => {
    if (!address) return;

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${orgId}/members/${memberId}/role`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            'X-Wallet-Address': address,
          },
          body: JSON.stringify({ roleId: newRoleId }),
        }
      );

      if (response.ok) {
        toast.success('Member role updated');
        loadMembers();
      } else {
        throw new Error('Failed to update role');
      }
    } catch (error) {
      console.error('Failed to update role:', error);
      toast.error('Failed to update member role');
    }
  };

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const formatDate = (date?: string) => {
    if (!date) return 'N/A';
    return new Date(date).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  };

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
            <Users className="h-8 w-8 text-blue-500" />
            Team Members
          </h1>
          <p className="text-gray-600 mt-1">
            Manage who has access to {currentOrg?.name || 'your organization'}
          </p>
        </div>
        <Button onClick={() => setShowInviteModal(true)}>
          <Mail className="mr-2 h-4 w-4" />
          Invite Member
        </Button>
      </div>

      {/* Stats */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Members</CardTitle>
            <UserCheck className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{members.filter(m => m.status === 'active').length}</div>
            <p className="text-xs text-gray-500 mt-1">Active team members</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Pending</CardTitle>
            <Clock className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{members.filter(m => m.status === 'pending').length}</div>
            <p className="text-xs text-gray-500 mt-1">Awaiting acceptance</p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Roles</CardTitle>
            <Shield className="h-4 w-4 text-gray-500" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{roles.length}</div>
            <p className="text-xs text-gray-500 mt-1">Available roles</p>
          </CardContent>
        </Card>
      </div>

      {/* Members List */}
      <Card>
        <CardHeader>
          <CardTitle>All Members</CardTitle>
          <CardDescription>
            {members.length} member{members.length !== 1 ? 's' : ''} in your organization
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {members.length === 0 ? (
              <div className="text-center py-12">
                <UserX className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-600 mb-4">No team members yet</p>
                <Button onClick={() => setShowInviteModal(true)}>
                  <Mail className="mr-2 h-4 w-4" />
                  Invite Your First Member
                </Button>
              </div>
            ) : (
              members.map((member) => (
                <div
                  key={member.id}
                  className="flex items-center justify-between p-4 rounded-lg border hover:bg-gray-50 transition"
                >
                  <div className="flex items-center gap-4 flex-1">
                    {/* Avatar */}
                    <div className="h-10 w-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white font-bold">
                      {member.user.displayName?.[0] || member.user.walletAddress[2]}
                    </div>

                    {/* Info */}
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <p className="font-medium">
                          {member.user.displayName || member.user.ensName || formatAddress(member.user.walletAddress)}
                        </p>
                        <RoleBadge role={member.role.name} size="sm" />
                        {member.status === 'pending' && (
                          <span className="text-xs px-2 py-0.5 rounded-full bg-yellow-100 text-yellow-700">
                            Pending
                          </span>
                        )}
                      </div>
                      <p className="text-sm text-gray-500">
                        {formatAddress(member.user.walletAddress)} • Joined {formatDate(member.joinedAt)}
                      </p>
                    </div>
                  </div>

                  {/* Actions */}
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="ghost" size="icon">
                        <MoreVertical className="h-4 w-4" />
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem
                        onClick={() => {
                          // Change role modal would go here
                          toast.info('Role change feature coming soon');
                        }}
                      >
                        <Shield className="mr-2 h-4 w-4" />
                        Change Role
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        className="text-red-600"
                        onClick={() =>
                          handleRemoveMember(
                            member.id,
                            member.user.displayName || formatAddress(member.user.walletAddress)
                          )
                        }
                      >
                        <Trash2 className="mr-2 h-4 w-4" />
                        Remove Member
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              ))
            )}
          </div>
        </CardContent>
      </Card>

      {/* Invite Modal */}
      <InviteMemberModal
        open={showInviteModal}
        onOpenChange={setShowInviteModal}
        organizationId={orgId}
        roles={roles}
        onSuccess={loadMembers}
      />
    </div>
  );
}

