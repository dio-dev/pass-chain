// API client for Pass Chain backend

const API_BASE_URL = typeof window !== 'undefined' 
  ? (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080')
  : 'http://passchain-backend:8080'; // Internal cluster URL for SSR

export interface CreateCredentialRequest {
  name: string;
  username: string;
  url?: string;
  encryptedData: string;
  nonce: string;
  share1: string; // Will be stored in Vault
  share2: string; // Will be stored in Blockchain
  walletAddress: string;
  signature: string;
}

export interface Credential {
  id: string;
  name: string;
  username: string;
  url?: string;
  encryptedData: string;
  nonce: string;
  walletAddress: string;
  createdAt: string;
  lastAccessed?: string;
}

export interface RetrieveCredentialResponse extends Credential {
  share1: string; // Retrieved from Vault
  share2: string; // Retrieved from Blockchain
}

/**
 * Create a new credential
 */
export async function createCredential(
  data: CreateCredentialRequest
): Promise<{ id: string; share3: string }> {
  const response = await fetch(`${API_BASE_URL}/api/v1/credentials`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Wallet-Address': data.walletAddress,
      'X-Signature': data.signature,
    },
    body: JSON.stringify(data),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.message || 'Failed to create credential');
  }

  return response.json();
}

/**
 * Get all credentials for a wallet
 */
export async function getCredentials(
  walletAddress: string,
  signature: string
): Promise<Credential[]> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/credentials?wallet=${walletAddress}`,
    {
      headers: {
        'X-Wallet-Address': walletAddress,
        'X-Signature': signature,
      },
    }
  );

  if (!response.ok) {
    throw new Error('Failed to fetch credentials');
  }

  return response.json();
}

/**
 * Get a specific credential with decryption keys
 */
export async function getCredentialById(
  id: string,
  walletAddress: string,
  signature: string
): Promise<RetrieveCredentialResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/credentials/${id}`, {
    headers: {
      'X-Wallet-Address': walletAddress,
      'X-Signature': signature,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to retrieve credential');
  }

  return response.json();
}

/**
 * Delete a credential
 */
export async function deleteCredential(
  id: string,
  walletAddress: string,
  signature: string
): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/api/v1/credentials/${id}`, {
    method: 'DELETE',
    headers: {
      'X-Wallet-Address': walletAddress,
      'X-Signature': signature,
    },
  });

  if (!response.ok) {
    throw new Error('Failed to delete credential');
  }
}

/**
 * Health check
 */
export async function healthCheck(): Promise<{ status: string; service: string }> {
  const response = await fetch(`${API_BASE_URL}/health`);
  return response.json();
}

