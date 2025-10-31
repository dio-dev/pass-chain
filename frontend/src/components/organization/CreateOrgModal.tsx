'use client';

import React, { useState } from 'react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import { useOrganization } from '@/contexts/OrganizationContext';
import { Building, X } from 'lucide-react';
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
import { toast } from 'sonner';

interface CreateOrgModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function CreateOrgModal({ open, onOpenChange }: CreateOrgModalProps) {
  const { address } = useAccount();
  const router = useRouter();
  const { refreshOrganizations, switchOrg } = useOrganization();
  const [name, setName] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!name.trim()) {
      toast.error('Organization name is required');
      return;
    }

    if (!address) {
      toast.error('Please connect your wallet');
      return;
    }

    setIsCreating(true);

    try {
      const response = await fetch('http://localhost:8080/api/v1/organizations', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Wallet-Address': address,
        },
        body: JSON.stringify({ name: name.trim() }),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create organization');
      }

      const org = await response.json();

      // Success!
      toast.success(`🎉 ${org.name} created!`, {
        description: 'Your organization is ready. Invite your team to get started.',
      });

      // Refresh organizations list
      await refreshOrganizations();

      // Switch to new org
      switchOrg(org.id);

      // Close modal
      onOpenChange(false);
      setName('');

      // Navigate to org dashboard
      router.push(`/organizations/${org.id}`);
    } catch (error) {
      console.error('Failed to create organization:', error);
      toast.error('Failed to create organization', {
        description: error instanceof Error ? error.message : 'Please try again',
      });
    } finally {
      setIsCreating(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Building className="h-5 w-5 text-blue-500" />
            Create Organization
          </DialogTitle>
          <DialogDescription>
            Start collaborating with your team. Invite members and share passwords securely.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleCreate} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="org-name">Organization Name</Label>
            <Input
              id="org-name"
              placeholder="My Startup"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
              disabled={isCreating}
            />
            <p className="text-xs text-gray-500">
              💡 Choose a name your team will recognize
            </p>
          </div>

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 space-y-2">
            <p className="text-sm font-medium text-blue-900">
              ✨ What happens next:
            </p>
            <ul className="text-xs text-blue-700 space-y-1">
              <li>• Creates organization with default vault</li>
              <li>• Assigns you as Owner (full control)</li>
              <li>• Sets up 6 default roles for your team</li>
              <li>• Ready to invite members immediately</li>
            </ul>
          </div>

          <div className="flex gap-2 justify-end">
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isCreating}
            >
              Cancel
            </Button>
            <Button type="submit" disabled={isCreating || !name.trim()}>
              {isCreating ? (
                <>
                  <div className="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2" />
                  Creating...
                </>
              ) : (
                'Create Organization'
              )}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

