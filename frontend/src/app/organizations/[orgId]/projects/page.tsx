'use client';

import React, { useEffect, useState } from 'react';
import { useAccount } from 'wagmi';
import { useParams, useRouter } from 'next/navigation';
import { useOrganization } from '@/contexts/OrganizationContext';
import { FolderKanban, Plus, FolderOpen } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { ProjectCard } from '@/components/project/ProjectCard';
import { CreateProjectModal } from '@/components/project/CreateProjectModal';
import { toast } from 'sonner';

interface Project {
  id: string;
  name: string;
  description?: string;
  memberCount?: number;
  credentialCount?: number;
  createdAt: string;
}

export default function ProjectsPage() {
  const params = useParams();
  const router = useRouter();
  const { address } = useAccount();
  const { currentOrg } = useOrganization();
  const [projects, setProjects] = useState<Project[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showCreateModal, setShowCreateModal] = useState(false);

  const orgId = params.orgId as string;

  useEffect(() => {
    if (address && orgId) {
      loadProjects();
    }
  }, [address, orgId]);

  const loadProjects = async () => {
    if (!address) return;

    setIsLoading(true);
    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/organizations/${orgId}/projects`,
        {
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (response.ok) {
        const data = await response.json();
        setProjects(data.projects || []);
      }
    } catch (error) {
      console.error('Failed to load projects:', error);
      toast.error('Failed to load projects');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDeleteProject = async (projectId: string, projectName: string) => {
    if (!address) return;

    if (!confirm(`Delete project "${projectName}"? This will remove all associated data.`)) {
      return;
    }

    try {
      const response = await fetch(
        `http://localhost:8080/api/v1/projects/${projectId}`,
        {
          method: 'DELETE',
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (response.ok) {
        toast.success(`Project "${projectName}" deleted`);
        loadProjects();
      } else {
        throw new Error('Failed to delete project');
      }
    } catch (error) {
      console.error('Failed to delete project:', error);
      toast.error('Failed to delete project');
    }
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-96">
        <div className="animate-spin rounded-full h-8 w-8 border-2 border-green-500 border-t-transparent" />
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-3">
            <FolderKanban className="h-8 w-8 text-green-500" />
            Projects
          </h1>
          <p className="text-gray-600 mt-1">
            Organize credentials by team or application
          </p>
        </div>
        <Button onClick={() => setShowCreateModal(true)}>
          <Plus className="mr-2 h-4 w-4" />
          New Project
        </Button>
      </div>

      {/* Projects Grid */}
      {projects.length === 0 ? (
        <div className="text-center py-16 bg-gray-50 rounded-lg border-2 border-dashed">
          <FolderOpen className="h-16 w-16 text-gray-400 mx-auto mb-4" />
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            No projects yet
          </h3>
          <p className="text-gray-600 mb-6 max-w-md mx-auto">
            Projects help you organize credentials for different teams or applications.
            Create your first project to get started.
          </p>
          <Button onClick={() => setShowCreateModal(true)} size="lg">
            <Plus className="mr-2 h-5 w-5" />
            Create Your First Project
          </Button>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {projects.map((project) => (
            <ProjectCard
              key={project.id}
              project={project}
              onDelete={() => handleDeleteProject(project.id, project.name)}
            />
          ))}
        </div>
      )}

      {/* Create Project Modal */}
      <CreateProjectModal
        open={showCreateModal}
        onOpenChange={setShowCreateModal}
        organizationId={orgId}
        onSuccess={loadProjects}
      />
    </div>
  );
}

