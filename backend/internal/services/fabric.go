package fabric

import (
	"encoding/json"
	"fmt"

	"pass-chain/backend/pkg/logger"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

// FabricClient wraps Hyperledger Fabric SDK
type FabricClient struct {
	sdk           *fabsdk.FabricSDK
	channelClient *channel.Client
	channelID     string
	chaincode     string
	logger        *logger.Logger
}

// FabricConfig holds Fabric connection config
type FabricConfig struct {
	ConfigPath    string
	ChannelID     string
	ChaincodeID   string
	OrgName       string
	OrgUser       string
	PeerEndpoint  string
}

// NewFabricClient creates a new Fabric SDK client
func NewFabricClient(cfg *FabricConfig, log *logger.Logger) (*FabricClient, error) {
	// Initialize SDK from config file
	sdk, err := fabsdk.New(config.FromFile(cfg.ConfigPath))
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK: %w", err)
	}

	// Create channel client
	clientChannelContext := sdk.ChannelContext(cfg.ChannelID, fabsdk.WithUser(cfg.OrgUser), fabsdk.WithOrg(cfg.OrgName))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		sdk.Close()
		return nil, fmt.Errorf("failed to create channel client: %w", err)
	}

	log.Info("Fabric client initialized", "channel", cfg.ChannelID, "chaincode", cfg.ChaincodeID)

	return &FabricClient{
		sdk:           sdk,
		channelClient: channelClient,
		channelID:     cfg.ChannelID,
		chaincode     cfg.ChaincodeID,
		logger:        log,
	}, nil
}

// Close closes the Fabric SDK
func (fc *FabricClient) Close() {
	if fc.sdk != nil {
		fc.sdk.Close()
	}
}

// StoreShard stores Share2 in Fabric blockchain (Private Data Collection)
func (fc *FabricClient) StoreShard(wallet, credentialID, share2 string) (string, error) {
	req := channel.Request{
		ChaincodeID: fc.chaincode,
		Fcn:         "StoreShard",
		Args:        [][]byte{[]byte(wallet), []byte(credentialID), []byte(share2)},
	}

	resp, err := fc.channelClient.Execute(req)
	if err != nil {
		fc.logger.Error("Failed to store shard in Fabric", "error", err)
		return "", fmt.Errorf("failed to store shard: %w", err)
	}

	fc.logger.Info("Shard stored in Fabric", "txID", resp.TransactionID, "credentialID", credentialID)
	return string(resp.TransactionID), nil
}

// ReadShard retrieves Share2 from Fabric blockchain
func (fc *FabricClient) ReadShard(wallet, credentialID string) (string, error) {
	req := channel.Request{
		ChaincodeID: fc.chaincode,
		Fcn:         "ReadShard",
		Args:        [][]byte{[]byte(wallet), []byte(credentialID)},
	}

	resp, err := fc.channelClient.Query(req)
	if err != nil {
		fc.logger.Error("Failed to read shard from Fabric", "error", err)
		return "", fmt.Errorf("failed to read shard: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &result); err != nil {
		return "", fmt.Errorf("failed to parse shard response: %w", err)
	}

	share2, ok := result["share2"].(string)
	if !ok {
		return "", fmt.Errorf("invalid share2 format")
	}

	fc.logger.Info("Shard retrieved from Fabric", "credentialID", credentialID)
	return share2, nil
}

// LogAccess logs credential access to Fabric blockchain
func (fc *FabricClient) LogAccess(wallet, credentialID, action, ipHash string) (string, error) {
	req := channel.Request{
		ChaincodeID: fc.chaincode,
		Fcn:         "LogAccess",
		Args:        [][]byte{
			[]byte(wallet),
			[]byte(credentialID),
			[]byte(action),
			[]byte(ipHash),
		},
	}

	resp, err := fc.channelClient.Execute(req)
	if err != nil {
		fc.logger.Error("Failed to log access to Fabric", "error", err)
		return "", fmt.Errorf("failed to log access: %w", err)
	}

	fc.logger.Info("Access logged to Fabric", "txID", resp.TransactionID, "action", action)
	return string(resp.TransactionID), nil
}

// GetAuditLogs retrieves audit logs for a wallet
func (fc *FabricClient) GetAuditLogs(wallet string) ([]map[string]interface{}, error) {
	req := channel.Request{
		ChaincodeID: fc.chaincode,
		Fcn:         "GetAuditLogs",
		Args:        [][]byte{[]byte(wallet)},
	}

	resp, err := fc.channelClient.Query(req)
	if err != nil {
		fc.logger.Error("Failed to get audit logs from Fabric", "error", err)
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}

	var logs []map[string]interface{}
	if err := json.Unmarshal(resp.Payload, &logs); err != nil {
		return nil, fmt.Errorf("failed to parse audit logs: %w", err)
	}

	return logs, nil
}

