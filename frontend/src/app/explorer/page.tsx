// Blockchain explorer page for Pass Chain
'use client';

import { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { Button } from '@/components/ui/button';
import { formatDistanceToNow } from 'date-fns';
import { 
  Shield, 
  Clock, 
  Eye, 
  Trash2, 
  Key, 
  Link as LinkIcon,
  Activity,
  TrendingUp,
  Database,
  Lock
} from 'lucide-react';
import { toast } from 'sonner';

interface AuditLog {
  id: string;
  credentialId: string;
  credentialName: string;
  action: string; // "create", "read", "delete"
  timestamp: string;
  txHash?: string;
  ipHash?: string;
  blockNumber?: number;
}

interface BlockchainStats {
  totalCredentials: number;
  totalAccesses: number;
  lastActivity: string;
  vaultShards: number;
  blockchainShards: number;
}

export default function BlockchainExplorerPage() {
  const { address, isConnected } = useAccount();
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [stats, setStats] = useState<BlockchainStats | null>(null);
  const [loading, setLoading] = useState(false);
  const [filter, setFilter] = useState<'all' | 'create' | 'read' | 'delete'>('all');

  useEffect(() => {
    if (isConnected && address) {
      fetchAuditLogs();
      fetchStats();
    }
  }, [isConnected, address]);

  const fetchAuditLogs = async () => {
    if (!address) return;
    
    setLoading(true);
    try {
      // Fetch from backend
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/audit-logs?wallet=${address}`,
        {
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to fetch audit logs');
      }

      const data = await response.json();
      setAuditLogs(data.logs || []);
    } catch (error) {
      console.error('Error fetching audit logs:', error);
      toast.error('Failed to fetch blockchain activity');
    } finally {
      setLoading(false);
    }
  };

  const fetchStats = async () => {
    if (!address) return;

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/api/v1/stats?wallet=${address}`,
        {
          headers: {
            'X-Wallet-Address': address,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to fetch stats');
      }

      const data = await response.json();
      setStats(data);
    } catch (error) {
      console.error('Error fetching stats:', error);
    }
  };

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'create':
        return <Key className="h-4 w-4 text-green-500" />;
      case 'read':
        return <Eye className="h-4 w-4 text-blue-500" />;
      case 'delete':
        return <Trash2 className="h-4 w-4 text-red-500" />;
      default:
        return <Activity className="h-4 w-4 text-gray-500" />;
    }
  };

  const getActionColor = (action: string) => {
    switch (action) {
      case 'create':
        return 'bg-green-50 text-green-700 border-green-200';
      case 'read':
        return 'bg-blue-50 text-blue-700 border-blue-200';
      case 'delete':
        return 'bg-red-50 text-red-700 border-red-200';
      default:
        return 'bg-gray-50 text-gray-700 border-gray-200';
    }
  };

  const filteredLogs = auditLogs.filter(log => 
    filter === 'all' ? true : log.action === filter
  );

  if (!isConnected) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 flex items-center justify-center p-4">
        <div className="text-center">
          <Shield className="h-16 w-16 text-blue-600 mx-auto mb-4" />
          <h1 className="text-2xl font-bold mb-2">Connect Your Wallet</h1>
          <p className="text-gray-600">
            Please connect your wallet to view blockchain activity
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 p-6">
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Blockchain Explorer
          </h1>
          <p className="text-gray-600">
            Immutable audit trail of all your credential activities
          </p>
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <div className="bg-blue-100 p-3 rounded-lg">
                  <Database className="h-6 w-6 text-blue-600" />
                </div>
                <TrendingUp className="h-5 w-5 text-green-500" />
              </div>
              <p className="text-gray-600 text-sm">Total Credentials</p>
              <p className="text-3xl font-bold text-gray-900">{stats.totalCredentials}</p>
            </div>

            <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <div className="bg-purple-100 p-3 rounded-lg">
                  <Activity className="h-6 w-6 text-purple-600" />
                </div>
              </div>
              <p className="text-gray-600 text-sm">Total Accesses</p>
              <p className="text-3xl font-bold text-gray-900">{stats.totalAccesses}</p>
            </div>

            <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <div className="bg-green-100 p-3 rounded-lg">
                  <Shield className="h-6 w-6 text-green-600" />
                </div>
              </div>
              <p className="text-gray-600 text-sm">Vault Shards</p>
              <p className="text-3xl font-bold text-gray-900">{stats.vaultShards}</p>
            </div>

            <div className="bg-white rounded-xl shadow-md p-6 border border-gray-100">
              <div className="flex items-center justify-between mb-2">
                <div className="bg-orange-100 p-3 rounded-lg">
                  <Lock className="h-6 w-6 text-orange-600" />
                </div>
              </div>
              <p className="text-gray-600 text-sm">Blockchain Shards</p>
              <p className="text-3xl font-bold text-gray-900">{stats.blockchainShards}</p>
            </div>
          </div>
        )}

        {/* Filters */}
        <div className="bg-white rounded-xl shadow-md p-4 mb-6">
          <div className="flex gap-2">
            <Button
              variant={filter === 'all' ? 'default' : 'outline'}
              onClick={() => setFilter('all')}
              className="flex-1"
            >
              All Activities
            </Button>
            <Button
              variant={filter === 'create' ? 'default' : 'outline'}
              onClick={() => setFilter('create')}
              className="flex-1"
            >
              <Key className="h-4 w-4 mr-2" />
              Created
            </Button>
            <Button
              variant={filter === 'read' ? 'default' : 'outline'}
              onClick={() => setFilter('read')}
              className="flex-1"
            >
              <Eye className="h-4 w-4 mr-2" />
              Accessed
            </Button>
            <Button
              variant={filter === 'delete' ? 'default' : 'outline'}
              onClick={() => setFilter('delete')}
              className="flex-1"
            >
              <Trash2 className="h-4 w-4 mr-2" />
              Deleted
            </Button>
          </div>
        </div>

        {/* Audit Logs */}
        <div className="bg-white rounded-xl shadow-md overflow-hidden">
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-bold text-gray-900">Audit Trail</h2>
              <Button
                onClick={fetchAuditLogs}
                disabled={loading}
                variant="outline"
                size="sm"
              >
                {loading ? 'Refreshing...' : 'Refresh'}
              </Button>
            </div>
          </div>

          <div className="divide-y divide-gray-200">
            {loading && filteredLogs.length === 0 ? (
              <div className="p-12 text-center">
                <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4 animate-pulse" />
                <p className="text-gray-600">Loading blockchain activity...</p>
              </div>
            ) : filteredLogs.length === 0 ? (
              <div className="p-12 text-center">
                <Activity className="h-12 w-12 text-gray-400 mx-auto mb-4" />
                <p className="text-gray-600">No activity recorded yet</p>
                <p className="text-sm text-gray-500 mt-2">
                  Your credential actions will appear here
                </p>
              </div>
            ) : (
              filteredLogs.map((log) => (
                <div
                  key={log.id}
                  className="p-6 hover:bg-gray-50 transition-colors"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-4 flex-1">
                      <div className="mt-1">{getActionIcon(log.action)}</div>
                      
                      <div className="flex-1">
                        <div className="flex items-center gap-2 mb-1">
                          <span
                            className={`px-3 py-1 rounded-full text-xs font-medium border ${getActionColor(
                              log.action
                            )}`}
                          >
                            {log.action.toUpperCase()}
                          </span>
                          <span className="text-sm font-semibold text-gray-900">
                            {log.credentialName}
                          </span>
                        </div>

                        <div className="flex items-center gap-4 text-sm text-gray-600 mt-2">
                          <div className="flex items-center gap-1">
                            <Clock className="h-3 w-3" />
                            {formatDistanceToNow(new Date(log.timestamp), {
                              addSuffix: true,
                            })}
                          </div>

                          {log.txHash && (
                            <div className="flex items-center gap-1">
                              <LinkIcon className="h-3 w-3" />
                              <span className="font-mono text-xs">
                                {log.txHash.slice(0, 8)}...{log.txHash.slice(-6)}
                              </span>
                            </div>
                          )}

                          {log.blockNumber && (
                            <div className="flex items-center gap-1">
                              <Database className="h-3 w-3" />
                              Block #{log.blockNumber}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="text-right">
                      <p className="text-xs text-gray-500">
                        {new Date(log.timestamp).toLocaleString()}
                      </p>
                    </div>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>

        {/* Info Banner */}
        <div className="mt-8 bg-blue-50 border border-blue-200 rounded-xl p-6">
          <div className="flex items-start gap-4">
            <Shield className="h-6 w-6 text-blue-600 flex-shrink-0 mt-1" />
            <div>
              <h3 className="font-semibold text-blue-900 mb-2">
                Immutable Security
              </h3>
              <p className="text-sm text-blue-800">
                Every action is cryptographically signed and stored on Hyperledger Fabric blockchain.
                This audit trail cannot be tampered with or deleted, ensuring complete transparency and security.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

