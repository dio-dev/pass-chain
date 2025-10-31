'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Building, ArrowRight, Sparkles } from 'lucide-react';
import { CreateOrgModal } from '@/components/organization/CreateOrgModal';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

export default function NewOrganizationPage() {
  const router = useRouter();
  const [showModal, setShowModal] = useState(false);

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      {/* Hero Section */}
      <div className="text-center space-y-4">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 mb-4">
          <Building className="h-8 w-8 text-white" />
        </div>
        <h1 className="text-4xl font-bold">Create Your Organization</h1>
        <p className="text-xl text-gray-600 max-w-2xl mx-auto">
          Start collaborating with your team. Share passwords securely and manage access with powerful roles.
        </p>
      </div>

      {/* Benefits Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Card>
          <CardHeader>
            <div className="h-10 w-10 rounded-lg bg-blue-100 text-blue-600 flex items-center justify-center mb-2">
              🔐
            </div>
            <CardTitle className="text-lg">Team Vaults</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription>
              Share passwords securely with your entire team or specific members
            </CardDescription>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="h-10 w-10 rounded-lg bg-green-100 text-green-600 flex items-center justify-center mb-2">
              📁
            </div>
            <CardTitle className="text-lg">Projects</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription>
              Organize credentials into workspaces for different teams or apps
            </CardDescription>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <div className="h-10 w-10 rounded-lg bg-purple-100 text-purple-600 flex items-center justify-center mb-2">
              👥
            </div>
            <CardTitle className="text-lg">Role Management</CardTitle>
          </CardHeader>
          <CardContent>
            <CardDescription>
              Control access with 6 built-in roles: Owner, Admin, Member, and more
            </CardDescription>
          </CardContent>
        </Card>
      </div>

      {/* What Happens */}
      <Card className="border-2 border-blue-200 bg-blue-50/50">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Sparkles className="h-5 w-5 text-blue-500" />
            What Happens When You Create
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold">
              1
            </div>
            <div>
              <p className="font-medium">Organization Created</p>
              <p className="text-sm text-gray-600">Your organization is instantly ready</p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold">
              2
            </div>
            <div>
              <p className="font-medium">Default Vault Created</p>
              <p className="text-sm text-gray-600">A shared vault for your entire team</p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold">
              3
            </div>
            <div>
              <p className="font-medium">Roles Set Up</p>
              <p className="text-sm text-gray-600">6 default roles ready for your team</p>
            </div>
          </div>

          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-blue-500 text-white flex items-center justify-center text-sm font-bold">
              4
            </div>
            <div>
              <p className="font-medium">You&apos;re the Owner</p>
              <p className="text-sm text-gray-600">Full control over your organization</p>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* CTA Buttons */}
      <div className="flex gap-4 justify-center">
        <Button
          variant="outline"
          size="lg"
          onClick={() => router.push('/dashboard')}
        >
          Maybe Later
        </Button>
        <Button
          size="lg"
          onClick={() => setShowModal(true)}
          className="gap-2"
        >
          Create Organization
          <ArrowRight className="h-4 w-4" />
        </Button>
      </div>

      {/* Create Org Modal */}
      <CreateOrgModal open={showModal} onOpenChange={setShowModal} />
    </div>
  );
}

