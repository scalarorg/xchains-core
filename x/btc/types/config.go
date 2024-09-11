package types

import "github.com/btcsuite/btcd/rpcclient"

// EVMConfig contains all EVM module configuration values
type BTCConfig struct {
	rpcclient.ConnConfig
	Name    string `mapstructure:"name"`
	ChainID string `mapstructure:"chain_id"`
}

// DefaultConfig returns a configuration populated with default values
func DefaultConfig() []BTCConfig {
	return []BTCConfig{{
		rpcclient.ConnConfig{
			Host: "http://127.0.0.1:18332",
			User: "user",
			Pass: "password",
		},
		"bitcoin",
		"bitcoin-regtest",
	}}
}
func (c *BTCConfig) GetRPCConfig() *rpcclient.ConnConfig {
	return &c.ConnConfig
}
