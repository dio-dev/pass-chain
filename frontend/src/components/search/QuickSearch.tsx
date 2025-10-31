'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import {
  Search,
  Lock,
  FolderKanban,
  Users,
  Building,
  ArrowRight,
  Clock,
} from 'lucide-react';
import {
  Dialog,
  DialogContent,
} from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';

interface SearchResult {
  type: 'credential' | 'project' | 'member' | 'organization';
  id: string;
  title: string;
  subtitle: string;
  icon: React.ReactNode;
  url: string;
}

interface QuickSearchProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export function QuickSearch({ open, onOpenChange }: QuickSearchProps) {
  const { address } = useAccount();
  const router = useRouter();
  const [query, setQuery] = useState('');
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedIndex, setSelectedIndex] = useState(0);

  // Search function
  const performSearch = useCallback(async (searchQuery: string) => {
    if (!searchQuery.trim() || !address) {
      setResults([]);
      return;
    }

    setIsSearching(true);

    try {
      // In production, this would be a single API call
      // For now, we'll simulate searching multiple endpoints
      const allResults: SearchResult[] = [];

      // Search credentials
      try {
        const credsResponse = await fetch(
          `http://localhost:8080/api/v1/credentials?wallet=${address}`,
          {
            headers: { 'X-Wallet-Address': address },
          }
        );
        if (credsResponse.ok) {
          const data = await credsResponse.json();
          const filtered = data.credentials?.filter((c: any) =>
            c.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            c.username?.toLowerCase().includes(searchQuery.toLowerCase()) ||
            c.url?.toLowerCase().includes(searchQuery.toLowerCase())
          ) || [];
          
          allResults.push(...filtered.slice(0, 5).map((c: any) => ({
            type: 'credential' as const,
            id: c.id,
            title: c.name,
            subtitle: c.username || c.url || 'Credential',
            icon: <Lock className="h-4 w-4 text-purple-500" />,
            url: `/dashboard?credential=${c.id}`,
          })));
        }
      } catch (e) {
        console.error('Failed to search credentials:', e);
      }

      setResults(allResults);
      setSelectedIndex(0);
    } catch (error) {
      console.error('Search failed:', error);
    } finally {
      setIsSearching(false);
    }
  }, [address]);

  // Debounced search
  useEffect(() => {
    const timer = setTimeout(() => {
      performSearch(query);
    }, 300);

    return () => clearTimeout(timer);
  }, [query, performSearch]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!open) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) => Math.min(prev + 1, results.length - 1));
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) => Math.max(prev - 1, 0));
          break;
        case 'Enter':
          e.preventDefault();
          if (results[selectedIndex]) {
            router.push(results[selectedIndex].url);
            onOpenChange(false);
            setQuery('');
          }
          break;
        case 'Escape':
          e.preventDefault();
          onOpenChange(false);
          setQuery('');
          break;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [open, results, selectedIndex, router, onOpenChange]);

  // Reset on close
  useEffect(() => {
    if (!open) {
      setQuery('');
      setResults([]);
      setSelectedIndex(0);
    }
  }, [open]);

  const handleResultClick = (result: SearchResult) => {
    router.push(result.url);
    onOpenChange(false);
    setQuery('');
  };

  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'credential':
        return <Lock className="h-4 w-4 text-purple-500" />;
      case 'project':
        return <FolderKanban className="h-4 w-4 text-green-500" />;
      case 'member':
        return <Users className="h-4 w-4 text-blue-500" />;
      case 'organization':
        return <Building className="h-4 w-4 text-blue-500" />;
      default:
        return <Search className="h-4 w-4 text-gray-500" />;
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-2xl p-0 gap-0 overflow-hidden">
        {/* Search Input */}
        <div className="flex items-center gap-3 px-4 py-3 border-b">
          <Search className="h-5 w-5 text-gray-400" />
          <Input
            placeholder="Search credentials, projects, members..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="border-0 focus-visible:ring-0 focus-visible:ring-offset-0 px-0"
            autoFocus
          />
          {isSearching && (
            <div className="animate-spin rounded-full h-4 w-4 border-2 border-gray-300 border-t-blue-500" />
          )}
        </div>

        {/* Results */}
        <div className="max-h-[400px] overflow-y-auto">
          {query.trim() === '' ? (
            <div className="p-8 text-center text-gray-500">
              <Search className="h-12 w-12 mx-auto mb-3 text-gray-300" />
              <p className="text-sm">Start typing to search...</p>
              <p className="text-xs text-gray-400 mt-2">
                Search across credentials, projects, and members
              </p>
            </div>
          ) : results.length === 0 && !isSearching ? (
            <div className="p-8 text-center text-gray-500">
              <Search className="h-12 w-12 mx-auto mb-3 text-gray-300" />
              <p className="text-sm">No results found for &quot;{query}&quot;</p>
              <p className="text-xs text-gray-400 mt-2">
                Try a different search term
              </p>
            </div>
          ) : (
            <div className="py-2">
              {results.map((result, index) => (
                <button
                  key={`${result.type}-${result.id}`}
                  onClick={() => handleResultClick(result)}
                  className={`w-full flex items-center gap-3 px-4 py-3 hover:bg-gray-50 transition ${
                    index === selectedIndex ? 'bg-blue-50' : ''
                  }`}
                >
                  <div className="flex-shrink-0">{result.icon}</div>
                  <div className="flex-1 text-left min-w-0">
                    <p className="font-medium text-sm truncate">{result.title}</p>
                    <p className="text-xs text-gray-500 truncate">{result.subtitle}</p>
                  </div>
                  <ArrowRight className="h-4 w-4 text-gray-400 flex-shrink-0" />
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="px-4 py-2 border-t bg-gray-50 flex items-center justify-between text-xs text-gray-500">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-white border rounded text-[10px]">↑</kbd>
              <kbd className="px-1.5 py-0.5 bg-white border rounded text-[10px]">↓</kbd>
              <span>Navigate</span>
            </div>
            <div className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-white border rounded text-[10px]">↵</kbd>
              <span>Select</span>
            </div>
            <div className="flex items-center gap-1">
              <kbd className="px-1.5 py-0.5 bg-white border rounded text-[10px]">ESC</kbd>
              <span>Close</span>
            </div>
          </div>
          {results.length > 0 && (
            <span>{results.length} result{results.length !== 1 ? 's' : ''}</span>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}

