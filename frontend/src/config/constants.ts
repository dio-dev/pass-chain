export const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export const STORAGE_FEE = '0.001'; // ETH
export const USAGE_FEE = '0.0001'; // ETH

export const SUPPORTED_CHAINS = {
  mainnet: 1,
  polygon: 137,
  sepolia: 11155111,
} as const;




