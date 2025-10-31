'use client';

import React from 'react';

interface RoleBadgeProps {
  role: string;
  size?: 'sm' | 'md' | 'lg';
}

const roleColors: Record<string, { bg: string; text: string; emoji: string }> = {
  Owner: { bg: 'bg-purple-100', text: 'text-purple-700', emoji: '👑' },
  Admin: { bg: 'bg-blue-100', text: 'text-blue-700', emoji: '⚡' },
  'Security Officer': { bg: 'bg-red-100', text: 'text-red-700', emoji: '🛡️' },
  Member: { bg: 'bg-green-100', text: 'text-green-700', emoji: '👤' },
  Auditor: { bg: 'bg-yellow-100', text: 'text-yellow-700', emoji: '🔍' },
  Guest: { bg: 'bg-gray-100', text: 'text-gray-700', emoji: '👁️' },
};

const sizeClasses = {
  sm: 'text-xs px-2 py-0.5',
  md: 'text-sm px-2.5 py-1',
  lg: 'text-base px-3 py-1.5',
};

export function RoleBadge({ role, size = 'md' }: RoleBadgeProps) {
  const colors = roleColors[role] || roleColors.Member;
  
  return (
    <span
      className={`inline-flex items-center gap-1 rounded-full font-medium ${colors.bg} ${colors.text} ${sizeClasses[size]}`}
    >
      <span>{colors.emoji}</span>
      <span>{role}</span>
    </span>
  );
}

