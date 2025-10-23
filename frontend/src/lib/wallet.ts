// Wallet signing utilities

import { useSignMessage } from 'wagmi';
import { hashData } from './crypto';

/**
 * Generate a message to sign for authentication
 */
export function generateAuthMessage(action: string, timestamp?: number): string {
  const ts = timestamp || Date.now();
  return `Pass Chain Authentication\n\nAction: ${action}\nTimestamp: ${ts}\n\nThis signature proves you own this wallet.`;
}

/**
 * Generate a message to sign for credential operations
 */
export function generateCredentialMessage(
  action: 'create' | 'read' | 'delete',
  credentialId?: string,
  credentialName?: string
): string {
  const timestamp = Date.now();
  let message = `Pass Chain - ${action.toUpperCase()} Credential\n\n`;
  
  if (credentialName) {
    message += `Credential: ${credentialName}\n`;
  }
  if (credentialId) {
    message += `ID: ${credentialId}\n`;
  }
  
  message += `Timestamp: ${timestamp}\n\n`;
  message += `By signing, you authorize this ${action} operation.`;
  
  return message;
}

/**
 * Hook to sign messages with wallet
 */
export function useWalletSigner() {
  const { signMessageAsync } = useSignMessage();
  
  const signAuth = async (action: string) => {
    const message = generateAuthMessage(action);
    const signature = await signMessageAsync({ message });
    return { message, signature };
  };
  
  const signCredentialOperation = async (
    action: 'create' | 'read' | 'delete',
    credentialId?: string,
    credentialName?: string
  ) => {
    const message = generateCredentialMessage(action, credentialId, credentialName);
    const signature = await signMessageAsync({ message });
    return { message, signature };
  };
  
  return {
    signAuth,
    signCredentialOperation,
  };
}

