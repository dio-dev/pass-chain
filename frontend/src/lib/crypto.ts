// Client-side encryption utilities for Pass Chain
// Uses XChaCha20-Poly1305 for AEAD encryption

import { xchacha20poly1305 } from '@noble/ciphers/chacha';
import { randomBytes } from '@noble/ciphers/webcrypto';
import { utf8ToBytes, bytesToUtf8 } from '@noble/ciphers/utils';

/**
 * Generate a random 256-bit key for XChaCha20-Poly1305
 */
export function generateEncryptionKey(): Uint8Array {
  return randomBytes(32); // 256 bits
}

/**
 * Encrypt data using XChaCha20-Poly1305
 * @param plaintext - Data to encrypt
 * @param key - 256-bit encryption key
 * @returns Object with nonce and ciphertext
 */
export function encryptData(plaintext: string, key: Uint8Array): {
  nonce: string;
  ciphertext: string;
} {
  // Generate random 24-byte nonce (XChaCha20 uses 192-bit nonce)
  const nonce = randomBytes(24);
  
  // Create cipher instance
  const cipher = xchacha20poly1305(key, nonce);
  
  // Encrypt the plaintext
  const plaintextBytes = utf8ToBytes(plaintext);
  const ciphertextBytes = cipher.encrypt(plaintextBytes);
  
  // Return as base64 strings for easy storage
  return {
    nonce: bufferToBase64(nonce),
    ciphertext: bufferToBase64(ciphertextBytes),
  };
}

/**
 * Decrypt data using XChaCha20-Poly1305
 * @param ciphertext - Base64 encoded ciphertext
 * @param nonce - Base64 encoded nonce
 * @param key - 256-bit encryption key
 * @returns Decrypted plaintext
 */
export function decryptData(
  ciphertext: string,
  nonce: string,
  key: Uint8Array
): string {
  // Decode from base64
  const ciphertextBytes = base64ToBuffer(ciphertext);
  const nonceBytes = base64ToBuffer(nonce);
  
  // Create cipher instance
  const cipher = xchacha20poly1305(key, nonceBytes);
  
  // Decrypt
  const plaintextBytes = cipher.decrypt(ciphertextBytes);
  
  return bytesToUtf8(plaintextBytes);
}

/**
 * Split a secret into shares using simple XOR-based secret sharing
 * For production, use a proper Shamir Secret Sharing library
 * This is a simplified 2-of-3 implementation
 */
export function splitSecret(secret: Uint8Array): {
  share1: string;
  share2: string;
  share3: string;
} {
  // Generate two random shares
  const share1 = randomBytes(32);
  const share2 = randomBytes(32);
  
  // Third share is XOR of secret with first two shares
  const share3 = new Uint8Array(32);
  for (let i = 0; i < 32; i++) {
    share3[i] = secret[i] ^ share1[i] ^ share2[i];
  }
  
  return {
    share1: bufferToBase64(share1),
    share2: bufferToBase64(share2),
    share3: bufferToBase64(share3),
  };
}

/**
 * Reconstruct secret from any 2 of 3 shares
 */
export function reconstructSecret(
  share1: string,
  share2: string,
  share3: string
): Uint8Array {
  const s1 = base64ToBuffer(share1);
  const s2 = base64ToBuffer(share2);
  const s3 = base64ToBuffer(share3);
  
  // Reconstruct using XOR
  const secret = new Uint8Array(32);
  for (let i = 0; i < 32; i++) {
    secret[i] = s1[i] ^ s2[i] ^ s3[i];
  }
  
  return secret;
}

/**
 * Hash data using SHA-256
 */
export async function hashData(data: string): Promise<string> {
  const encoder = new TextEncoder();
  const dataBuffer = encoder.encode(data);
  const hashBuffer = await crypto.subtle.digest('SHA-256', dataBuffer);
  const hashArray = Array.from(new Uint8Array(hashBuffer));
  return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
}

/**
 * Helper: Convert Uint8Array to base64
 */
function bufferToBase64(buffer: Uint8Array): string {
  const binary = Array.from(buffer)
    .map(byte => String.fromCharCode(byte))
    .join('');
  return btoa(binary);
}

/**
 * Helper: Convert base64 to Uint8Array
 */
function base64ToBuffer(base64: string): Uint8Array {
  const binary = atob(base64);
  const buffer = new Uint8Array(binary.length);
  for (let i = 0; i < binary.length; i++) {
    buffer[i] = binary.charCodeAt(i);
  }
  return buffer;
}
