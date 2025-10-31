'use client';

import { useState, useEffect } from 'react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { Shield, Plus, Key, Eye, EyeOff, Lock, Trash2, Copy, Clock, Download, Activity } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useWalletSigner } from '@/lib/wallet';
import {
  generateEncryptionKey,
  encryptData,
  decryptData,
  splitSecret,
  reconstructSecret,
} from '@/lib/crypto';
import {
  createCredential,
  getCredentials,
  getCredentialById,
  deleteCredential,
  getCredentialAuditLogs,
  type Credential,
  type AuditLog,
} from '@/lib/api';
import { toast } from 'sonner';

interface CredentialWithPassword extends Credential {
  decryptedPassword?: string;
}

export default function DashboardPage() {
  const { address, isConnected } = useAccount();
  const router = useRouter();
  const { signCredentialOperation } = useWalletSigner();
  const [credentials, setCredentials] = useState<CredentialWithPassword[]>([]);
  const [loading, setLoading] = useState(false);
  const [showAddModal, setShowAddModal] = useState(false);
  const [revealedPasswords, setRevealedPasswords] = useState<Set<string>>(new Set());
  const [historyModal, setHistoryModal] = useState<{id: string, name: string} | null>(null);

  // Load credentials on mount
  useEffect(() => {
    if (isConnected && address) {
      loadCredentials();
    }
  }, [isConnected, address]);

  const loadCredentials = async () => {
    if (!address) return;
    
    try {
      setLoading(true);
      const { signature } = await signCredentialOperation('read');
      const creds = await getCredentials(address, signature);
      setCredentials(creds);
      toast.success(`Loaded ${creds.length} credentials`);
    } catch (error) {
      console.error('Failed to load credentials:', error);
      toast.error('Failed to load credentials');
    } finally {
      setLoading(false);
    }
  };

  const togglePasswordVisibility = async (id: string) => {
    if (revealedPasswords.has(id)) {
      // Hide password
      setRevealedPasswords(prev => {
        const newSet = new Set(prev);
        newSet.delete(id);
        return newSet;
      });
      // Clear decrypted password
      setCredentials(prev =>
        prev.map(c => (c.id === id ? { ...c, decryptedPassword: undefined } : c))
      );
    } else {
      // Reveal password - decrypt it
      try {
        if (!address) return;
        
        toast.info('Retrieving encryption keys...');
        const { signature } = await signCredentialOperation('read', id);
        const credentialData = await getCredentialById(id, address, signature);
        
        // Get user's backup share from localStorage
        // Try to get Share3 from localStorage (optional backup)
        const backupShares = JSON.parse(localStorage.getItem('pass-chain-backups') || '{}');
        const share3 = backupShares[id];
        
        // Reconstruct encryption key from available shares (need at least 2 of 3)
        // Share1 from Vault, Share2 from Database/Blockchain, Share3 from localStorage (optional)
        const reconstructedKey = reconstructSecret(
          credentialData.share1,
          credentialData.share2,
          share3 // Optional - works without it
        );
        
        // Decrypt the password
        const decrypted = decryptData(
          credentialData.encryptedData,
          credentialData.nonce,
          reconstructedKey
        );
        
        // Update state
        setCredentials(prev =>
          prev.map(c => (c.id === id ? { ...c, decryptedPassword: decrypted } : c))
        );
        setRevealedPasswords(prev => new Set(prev).add(id));
        toast.success('Password decrypted successfully!');
      } catch (error) {
        console.error('Failed to decrypt:', error);
        toast.error('Failed to decrypt password');
      }
    }
  };

  const handleDeleteCredential = async (id: string, name: string) => {
    if (!address || !confirm(`Delete credential "${name}"?`)) return;
    
    try {
      const { signature } = await signCredentialOperation('delete', id, name);
      await deleteCredential(id, address, signature);
      setCredentials(prev => prev.filter(c => c.id !== id));
      
      // Remove backup from localStorage
      const backupShares = JSON.parse(localStorage.getItem('pass-chain-backups') || '{}');
      delete backupShares[id];
      localStorage.setItem('pass-chain-backups', JSON.stringify(backupShares));
      
      toast.success('Credential deleted');
    } catch (error) {
      console.error('Failed to delete:', error);
      toast.error('Failed to delete credential');
    }
  };

  if (!isConnected) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center">
        <div className="text-center">
          <Shield className="h-20 w-20 text-purple-400 mx-auto mb-6" />
          <h1 className="text-4xl font-bold text-white mb-4">Connect Your Wallet</h1>
          <p className="text-gray-300 mb-8">Connect your Web3 wallet to access your secure credentials</p>
          <ConnectButton />
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Header */}
      <nav className="border-b border-purple-500/20 bg-slate-900/50 backdrop-blur-sm">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-3">
            <Shield className="h-8 w-8 text-purple-400" />
            <div>
              <h1 className="text-xl font-bold text-white">Pass Chain</h1>
              <p className="text-xs text-gray-400">Decentralized Password Manager</p>
            </div>
          </div>
          <div className="flex items-center gap-4">
            <div className="text-right mr-4">
              <p className="text-xs text-gray-400">Connected Wallet</p>
              <p className="text-sm text-white font-mono">
                {address?.slice(0, 6)}...{address?.slice(-4)}
              </p>
            </div>
            <ConnectButton />
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="container mx-auto px-4 py-8">
        {/* Stats */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          <StatCard
            icon={<Key className="h-6 w-6 text-purple-400" />}
            label="Total Credentials"
            value={credentials.length.toString()}
          />
          <StatCard
            icon={<Lock className="h-6 w-6 text-green-400" />}
            label="Encrypted"
            value="XChaCha20"
          />
          <StatCard
            icon={<Shield className="h-6 w-6 text-blue-400" />}
            label="Security"
            value="2-of-3 Split"
          />
          <StatCard
            icon={<Clock className="h-6 w-6 text-yellow-400" />}
            label="Wallet"
            value={`${address?.slice(0, 4)}...${address?.slice(-2)}`}
          />
        </div>

        {/* Actions */}
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-white">My Credentials</h2>
          <div className="flex gap-3">
            <Button
              onClick={() => router.push('/explorer')}
              variant="outline"
              className="border-purple-500/50 text-purple-400 hover:bg-purple-950"
            >
              <Activity className="h-5 w-5 mr-2" />
              Blockchain Explorer
            </Button>
            <Button
              onClick={() => setShowAddModal(true)}
              className="bg-purple-600 hover:bg-purple-700 text-white"
            >
              <Plus className="h-5 w-5 mr-2" />
              Add Credential
            </Button>
          </div>
        </div>

        {/* Loading */}
        {loading && (
          <div className="text-center py-20">
            <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-purple-500 mx-auto"></div>
            <p className="text-gray-400 mt-4">Loading credentials...</p>
          </div>
        )}

        {/* Credentials List */}
        {!loading && credentials.length > 0 && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {credentials.map((cred) => (
              <CredentialCard
                key={cred.id}
                credential={cred}
                isRevealed={revealedPasswords.has(cred.id)}
                onToggleReveal={() => togglePasswordVisibility(cred.id)}
                onDelete={() => handleDeleteCredential(cred.id, cred.name)}
                onHistoryClick={(cred) => setHistoryModal({id: cred.id, name: cred.name})}
              />
            ))}
          </div>
        )}

        {!loading && credentials.length === 0 && (
          <div className="text-center py-20">
            <Lock className="h-16 w-16 text-gray-600 mx-auto mb-4" />
            <h3 className="text-xl font-semibold text-gray-400 mb-2">No credentials yet</h3>
            <p className="text-gray-500 mb-6">Add your first credential to get started</p>
            <Button
              onClick={() => setShowAddModal(true)}
              className="bg-purple-600 hover:bg-purple-700 text-white"
            >
              <Plus className="h-5 w-5 mr-2" />
              Add Your First Credential
            </Button>
          </div>
        )}
      </main>

      {/* Add Credential Modal */}
      {showAddModal && (
        <AddCredentialModal
          walletAddress={address!}
          onClose={() => setShowAddModal(false)}
          onSuccess={loadCredentials}
        />
      )}

      {/* Access History Modal */}
      {historyModal && (
        <AccessHistoryModal
          credentialId={historyModal.id}
          credentialName={historyModal.name}
          walletAddress={address!}
          onClose={() => setHistoryModal(null)}
        />
      )}
    </div>
  );
}

function StatCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: string }) {
  return (
    <div className="bg-slate-800/50 backdrop-blur-sm border border-purple-500/20 rounded-lg p-4">
      <div className="flex items-center justify-between mb-2">
        <span className="text-gray-400 text-sm">{label}</span>
        {icon}
      </div>
      <p className="text-2xl font-bold text-white">{value}</p>
    </div>
  );
}

function CredentialCard({
  credential,
  isRevealed,
  onToggleReveal,
  onDelete,
  onHistoryClick,
}: {
  credential: CredentialWithPassword;
  isRevealed: boolean;
  onToggleReveal: () => void;
  onDelete: () => void;
  onHistoryClick: (cred: Credential) => void;
}) {
  const [copied, setCopied] = useState(false);

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    toast.success('Copied to clipboard!');
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-slate-800/50 backdrop-blur-sm border border-purple-500/20 rounded-lg p-5 hover:border-purple-500/40 transition-all">
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 bg-purple-600/20 rounded-lg flex items-center justify-center">
            <Key className="h-5 w-5 text-purple-400" />
          </div>
          <div>
            <h3 className="text-lg font-semibold text-white">{credential.name}</h3>
            {credential.url && (
              <p className="text-xs text-gray-400">{new URL(credential.url).hostname}</p>
            )}
          </div>
        </div>
      </div>

      <div className="space-y-2 mb-4">
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-400">Username</span>
          <div className="flex items-center gap-2">
            <span className="text-sm text-white font-mono">{credential.username}</span>
            <button
              onClick={() => copyToClipboard(credential.username)}
              className="text-gray-400 hover:text-white transition-colors"
            >
              <Copy className="h-4 w-4" />
            </button>
          </div>
        </div>

        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-400">Password</span>
          <div className="flex items-center gap-2">
            <span className="text-sm text-white font-mono">
              {isRevealed && credential.decryptedPassword
                ? credential.decryptedPassword
                : '••••••••'}
            </span>
            <button
              onClick={onToggleReveal}
              className="text-gray-400 hover:text-white transition-colors"
            >
              {isRevealed ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
            </button>
            {isRevealed && credential.decryptedPassword && (
              <button
                onClick={() => copyToClipboard(credential.decryptedPassword!)}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <Copy className="h-4 w-4" />
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="flex items-center justify-between pt-3 border-t border-slate-700">
        <div className="flex items-center gap-3">
          <span className="text-xs text-gray-500">
            {credential.lastAccessed ? `Accessed ${formatDate(credential.lastAccessed)}` : 'Never accessed'}
          </span>
          <button
            onClick={() => onHistoryClick(credential)}
            className="text-blue-400 hover:text-blue-300 transition-colors text-xs flex items-center gap-1"
          >
            <Clock className="h-3 w-3" />
            History
          </button>
        </div>
        <button
          onClick={onDelete}
          className="text-red-400 hover:text-red-300 transition-colors"
        >
          <Trash2 className="h-4 w-4" />
        </button>
      </div>
    </div>
  );
}

function AddCredentialModal({
  walletAddress,
  onClose,
  onSuccess,
}: {
  walletAddress: string;
  onClose: () => void;
  onSuccess: () => void;
}) {
  const { signCredentialOperation } = useWalletSigner();
  const [formData, setFormData] = useState({
    name: '',
    username: '',
    password: '',
    url: 'https://',
    urlProtocol: 'https://' as 'https://' | 'http://',
  });
  const [saving, setSaving] = useState(false);

  const handleUrlChange = (value: string) => {
    // Extract protocol if user pastes full URL
    if (value.startsWith('http://')) {
      setFormData({ ...formData, url: value, urlProtocol: 'http://' });
    } else if (value.startsWith('https://')) {
      setFormData({ ...formData, url: value, urlProtocol: 'https://' });
    } else {
      // Auto-add protocol
      setFormData({ ...formData, url: formData.urlProtocol + value.replace(/^(https?:\/\/)/, '') });
    }
  };

  const toggleProtocol = () => {
    const newProtocol = formData.urlProtocol === 'https://' ? 'http://' : 'https://';
    const urlWithoutProtocol = formData.url.replace(/^https?:\/\//, '');
    setFormData({ 
      ...formData, 
      urlProtocol: newProtocol,
      url: newProtocol + urlWithoutProtocol
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);

    try {
      toast.info('Generating encryption key...');
      
      // 1. Generate random encryption key
      const encryptionKey = generateEncryptionKey();
      
      // 2. Encrypt the password
      const { nonce, ciphertext } = encryptData(formData.password, encryptionKey);
      
      // 3. Split the key into 3 shares (2-of-3)
      const { share1, share2, share3 } = splitSecret(encryptionKey);
      
      // 4. Sign the transaction with wallet
      toast.info('Sign the transaction with your wallet...');
      const { signature } = await signCredentialOperation('create', undefined, formData.name);
      
      // 5. Send to backend (share1 → Vault, share2 → Blockchain)
      toast.info('Storing encrypted credential...');
      const result = await createCredential({
        name: formData.name,
        username: formData.username,
        url: formData.url || undefined,
        encryptedData: ciphertext,
        nonce,
        share1,
        share2,
        walletAddress,
        signature,
      });
      
      // 6. Store share3 locally as backup
      const backupShares = JSON.parse(localStorage.getItem('pass-chain-backups') || '{}');
      backupShares[result.id] = share3;
      localStorage.setItem('pass-chain-backups', JSON.stringify(backupShares));
      
      toast.success('✅ Credential encrypted and saved!');
      toast.info('💾 Backup key saved locally');
      
      onSuccess();
      onClose();
    } catch (error: any) {
      console.error('Failed to save credential:', error);
      toast.error(error.message || 'Failed to save credential');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-slate-800 border border-purple-500/20 rounded-lg p-6 w-full max-w-md">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-bold text-white">Add New Credential</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white" disabled={saving}>
            ✕
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Name *
            </label>
            <input
              type="text"
              required
              disabled={saving}
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:border-purple-500 focus:outline-none disabled:opacity-50"
              placeholder="e.g., GitHub, AWS, Gmail"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Username / Email *
            </label>
            <input
              type="text"
              required
              disabled={saving}
              value={formData.username}
              onChange={(e) => setFormData({ ...formData, username: e.target.value })}
              className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:border-purple-500 focus:outline-none disabled:opacity-50"
              placeholder="username or email"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Password *
            </label>
            <input
              type="password"
              required
              disabled={saving}
              value={formData.password}
              onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              className="w-full bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:border-purple-500 focus:outline-none disabled:opacity-50"
              placeholder="Enter password"
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-300 mb-2">
              Website URL (optional)
            </label>
            <div className="flex gap-2">
              <button
                type="button"
                onClick={toggleProtocol}
                disabled={saving}
                className="bg-slate-900 border border-slate-700 rounded-lg px-3 py-2 text-white hover:border-purple-500 focus:border-purple-500 focus:outline-none disabled:opacity-50 font-mono text-sm"
              >
                {formData.urlProtocol}
              </button>
              <input
                type="text"
                disabled={saving}
                value={formData.url.replace(/^https?:\/\//, '')}
                onChange={(e) => handleUrlChange(e.target.value)}
                className="flex-1 bg-slate-900 border border-slate-700 rounded-lg px-4 py-2 text-white focus:border-purple-500 focus:outline-none disabled:opacity-50"
                placeholder="example.com or github.com/username"
              />
            </div>
            <p className="text-xs text-gray-500 mt-1">Click protocol button to switch between https:// and http://</p>
          </div>

          <div className="bg-purple-900/20 border border-purple-500/30 rounded-lg p-4 mt-4">
            <p className="text-sm text-purple-300">
              <Lock className="h-4 w-4 inline mr-2" />
              Your password will be encrypted with <strong>XChaCha20-Poly1305</strong> and split into 3 shares:
            </p>
            <ul className="text-xs text-purple-200 mt-2 ml-6 space-y-1">
              <li>• Share 1 → HashiCorp Vault</li>
              <li>• Share 2 → Hyperledger Fabric</li>
              <li>• Share 3 → Your browser (backup)</li>
            </ul>
          </div>

          <div className="flex gap-3 pt-4">
            <Button
              type="button"
              onClick={onClose}
              variant="outline"
              disabled={saving}
              className="flex-1 border-gray-600 text-gray-300"
            >
              Cancel
            </Button>
            <Button
              type="submit"
              disabled={saving}
              className="flex-1 bg-purple-600 hover:bg-purple-700 text-white disabled:opacity-50"
            >
              {saving ? (
                <>Processing...</>
              ) : (
                <>
                  <Lock className="h-4 w-4 mr-2" />
                  Encrypt & Save
                </>
              )}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}

function AccessHistoryModal({
  credentialId,
  credentialName,
  walletAddress,
  onClose,
}: {
  credentialId: string;
  credentialName: string;
  walletAddress: string;
  onClose: () => void;
}) {
  const [logs, setLogs] = useState<AuditLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchHistory = async () => {
      try {
        const history = await getCredentialAuditLogs(credentialId, walletAddress);
        setLogs(history);
      } catch (error) {
        console.error('Failed to fetch history:', error);
        toast.error('Failed to load access history');
      } finally {
        setLoading(false);
      }
    };
    fetchHistory();
  }, [credentialId, walletAddress]);

  const getActionColor = (action: string) => {
    switch (action) {
      case 'create': return 'text-green-400';
      case 'read': return 'text-blue-400';
      case 'delete': return 'text-red-400';
      default: return 'text-gray-400';
    }
  };

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'create': return <Plus className="h-4 w-4" />;
      case 'read': return <Eye className="h-4 w-4" />;
      case 'delete': return <Trash2 className="h-4 w-4" />;
      default: return <Activity className="h-4 w-4" />;
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50">
      <div className="bg-slate-800 border border-purple-500/20 rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-hidden flex flex-col">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h2 className="text-2xl font-bold text-white">Access History</h2>
            <p className="text-sm text-gray-400 mt-1">{credentialName}</p>
          </div>
          <button onClick={onClose} className="text-gray-400 hover:text-white">
            ✕
          </button>
        </div>

        <div className="flex-1 overflow-y-auto">
          {loading && (
            <div className="text-center py-8">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-purple-500 mx-auto"></div>
              <p className="text-gray-400 mt-4">Loading history...</p>
            </div>
          )}

          {!loading && logs.length === 0 && (
            <div className="text-center py-8">
              <Clock className="h-12 w-12 text-gray-600 mx-auto mb-4" />
              <p className="text-gray-400">No access history yet</p>
            </div>
          )}

          {!loading && logs.length > 0 && (
            <div className="space-y-3">
              {logs.map((log) => (
                <div
                  key={log.id}
                  className="bg-slate-900/50 border border-slate-700 rounded-lg p-4 hover:border-purple-500/30 transition-all"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-3">
                      <div className={`mt-1 ${getActionColor(log.action)}`}>
                        {getActionIcon(log.action)}
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <span className={`font-semibold ${getActionColor(log.action)} capitalize`}>
                            {log.action}
                          </span>
                          {log.txHash && (
                            <span className="text-xs text-gray-500 font-mono">
                              tx: {log.txHash.slice(0, 8)}...
                            </span>
                          )}
                        </div>
                        <p className="text-sm text-gray-400 mt-1">
                          {new Date(log.timestamp).toLocaleString()}
                        </p>
                        {log.ipHash && (
                          <p className="text-xs text-gray-600 mt-1 font-mono">
                            IP: {log.ipHash.slice(0, 16)}...
                          </p>
                        )}
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="mt-6 pt-4 border-t border-slate-700">
          <div className="flex items-center gap-2 text-sm text-gray-400">
            <Shield className="h-4 w-4 text-purple-400" />
            <span>All access is logged on the blockchain for security and audit purposes</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function formatDate(dateString: string): string {
  const date = new Date(dateString);
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (hours < 1) return 'just now';
  if (hours < 24) return `${hours}h ago`;
  if (days < 7) return `${days}d ago`;
  return date.toLocaleDateString();
}
