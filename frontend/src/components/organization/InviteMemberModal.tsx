'use client';

import React, { useState } from 'react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import { Mail, X } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { toast } from 'sonner';

interface Role {
  id: string;
  name: string;
  description: string;
}

interface InviteMemberModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  organizationId: string;
  roles: Role[];
  onSuccess?: () => void;
}

export function InviteMemberModal({
  open,
  onOpenChange,
  organizationId,
  roles,
  onSuccess,
}: InviteMemberModalProps) {
  const { address } = useAccount();
  const [email, setEmail] = useState('');
  const [selectedRoleId, setSelectedRoleId] = useState<string>('');
  const [isInviting, setIsInviting] = useState(false);

  const handleInvite = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email.trim()) {
      toast.error('Email is required');
      return;
    }

    if (!selectedRoleId) {
      toast.error('Please select a role');
      return;
    }

    if (!address) {
      toast.error('Please connect your wallet');
      return;
    }

    setIsInviting(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${organizationId}/invite`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Wallet-Address': address,
          },
          body: JSON.stringify({
            email: email.trim(),
            roleId: selectedRoleId,
          }),
        }
      );

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to send invitation');
      }

      const data = await response.json();

      // Success!
      toast.success(`✉️ Invitation sent to ${email}`, {
        description: 'They will receive an email with instructions to join.',
      });

      // Reset form
      setEmail('');
      setSelectedRoleId('');

      // Close modal
      onOpenChange(false);

      // Callback
      if (onSuccess) {
        onSuccess();
      }
    } catch (error) {
      console.error('Failed to invite member:', error);
      toast.error('Failed to send invitation', {
        description: error instanceof Error ? error.message : 'Please try again',
      });
    } finally {
      setIsInviting(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Mail className="h-5 w-5 text-blue-500" />
            Invite Team Member
          </DialogTitle>
          <DialogDescription>
            Send an invitation to join your organization. They&apos;ll receive an email with a secure link.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleInvite} className="space-y-4">
          {/* Email Input */}
          <div className="space-y-2">
            <Label htmlFor="email">Email Address</Label>
            <Input
              id="email"
              type="email"
              placeholder="colleague@company.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoFocus
              disabled={isInviting}
            />
          </div>

          {/* Role Select */}
          <div className="space-y-2">
            <Label htmlFor="role">Role</Label>
            <Select value={selectedRoleId} onValueChange={setSelectedRoleId} disabled={isInviting}>
              <SelectTrigger>
                <SelectValue placeholder="Select a role" />
              </SelectTrigger>
              <SelectContent>
                {roles.map((role) => (
                  <SelectItem key={role.id} value={role.id}>
                    <div className="flex flex-col">
                      <span className="font-medium">{role.name}</span>
                      <span className="text-xs text-gray-500">{role.description}</span>
                    </div>
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Role Info */}
          {selectedRoleId && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
              <p className="text-sm font-medium text-blue-900 mb-1">
                🔐 Role Permissions
              </p>
              <p className="text-xs text-blue-700">
                {roles.find(r => r.id === selectedRoleId)?.description}
              </p>
            </div>
          )}

          {/* Actions */}
          <div className="flex gap-2 justify-end">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isInviting}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isInviting || !email.trim() || !selectedRoleId}>
              {isInviting ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2" />
                  Sending...
                </>
              ) : (
                <>
                  <Mail className="mr-2 h-4 w-4" />
                  Send Invitation
                </>
              )}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

