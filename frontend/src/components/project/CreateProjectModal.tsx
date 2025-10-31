'use client';

import React, { useState } from 'react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import { FolderKanban, X } from 'lucide-react';
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
import { Textarea } from '@/components/ui/textarea';
import { toast } from 'sonner';

interface CreateProjectModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  organizationId: string;
  onSuccess?: (project: any) => void;
}

export function CreateProjectModal({
  open,
  onOpenChange,
  organizationId,
  onSuccess,
}: CreateProjectModalProps) {
  const { address } = useAccount();
  const router = useRouter();
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [isCreating, setIsCreating] = useState(false);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim()) {
      toast.error('Project name is required');
      return;
    }

    if (!address) {
      toast.error('Please connect your wallet');
      return;
    }

    setIsCreating(true);

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${organizationId}/projects`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Wallet-Address': address,
          },
          body: JSON.stringify({
            name: name.trim(),
            description: description.trim() || undefined,
          }),
        }
      );

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Failed to create project');
      }

      const project = await response.json();

      // Success!
      toast.success(`🎉 ${project.name} created!`, {
        description: 'Project vault and permissions are ready.',
      });

      // Reset form
      setName('');
      setDescription('');

      // Close modal
      onOpenChange(false);

      // Callback
      if (onSuccess) {
        onSuccess(project);
      } else {
        // Navigate to project
        router.push(`/projects/${project.id}`);
      }
    } catch (error) {
      console.error('Failed to create project:', error);
      toast.error('Failed to create project', {
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
            <FolderKanban className="h-5 w-5 text-green-500" />
            Create Project
          </DialogTitle>
          <DialogDescription>
            Projects help organize credentials for specific teams or applications.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleCreate} className="space-y-4">
          {/* Project Name */}
          <div className="space-y-2">
            <Label htmlFor="project-name">Project Name</Label>
            <Input
              id="project-name"
              placeholder="Mobile App"
              value={name}
              onChange={(e) => setName(e.target.value)}
              autoFocus
              disabled={isCreating}
            />
            <p className="text-xs text-gray-500">
              💡 Choose a clear name like &quot;Mobile App&quot; or &quot;Production API&quot;
            </p>
          </div>

          {/* Description */}
          <div className="space-y-2">
            <Label htmlFor="description">Description (optional)</Label>
            <Textarea
              id="description"
              placeholder="Credentials for our mobile app team..."
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              rows={3}
              disabled={isCreating}
            />
          </div>

          {/* Info Box */}
          <div className="bg-green-50 border border-green-200 rounded-lg p-3 space-y-2">
            <p className="text-sm font-medium text-green-900">
              ✨ What&apos;s included:
            </p>
            <ul className="text-xs text-green-700 space-y-1">
              <li>• Dedicated project vault</li>
              <li>• Team member access control</li>
              <li>• Isolated credential storage</li>
              <li>• Project-specific audit logs</li>
            </ul>
          </div>

          {/* Actions */}
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
                'Create Project'
              )}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

