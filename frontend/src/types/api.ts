export interface Credential {
  id: string;
  wallet_address: string;
  credential_name: string;
  description: string;
  vault_path: string;
  blockchain_tx_id: string;
  storage_fee: string;
  created_at: string;
  updated_at: string;
}

export interface CredentialAccess {
  id: string;
  credential_id: string;
  wallet_address: string;
  ip_address: string;
  user_agent: string;
  usage_fee: string;
  blockchain_tx_id: string;
  success: boolean;
  accessed_at: string;
}

export interface CreateCredentialRequest {
  credential_name: string;
  description?: string;
  encrypted_data: string;
}

export interface RevealCredentialResponse {
  encrypted_data: string;
  blockchain_tx_id: string;
}




