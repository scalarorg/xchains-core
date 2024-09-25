package types

import (
	"time"

	"github.com/btcsuite/btcd/rpcclient"
)

// EVMConfig contains all EVM module configuration values
type BTCConfig struct {
	// Host is the IP address and port of the RPC server you want to connect
	// to.
	Host string

	// Endpoint is the websocket endpoint on the RPC server.  This is
	// typically "ws".
	Endpoint string

	// User is the username to use to authenticate to the RPC server.
	User string

	// Pass is the passphrase to use to authenticate to the RPC server.
	Pass string

	// CookiePath is the path to a cookie file containing the username and
	// passphrase to use to authenticate to the RPC server.  It is used
	// instead of User and Pass if non-empty.
	CookiePath string

	cookieLastCheckTime time.Time
	cookieLastModTime   time.Time
	cookieLastUser      string
	cookieLastPass      string
	cookieLastErr       error

	// Params is the string representing the network that the server
	// is running. If there is no parameter set in the config, then
	// mainnet will be used by default.
	Params string

	// DisableTLS specifies whether transport layer security should be
	// disabled.  It is recommended to always use TLS if the RPC server
	// supports it as otherwise your username and password is sent across
	// the wire in cleartext.
	DisableTLS bool

	// Certificates are the bytes for a PEM-encoded certificate chain used
	// for the TLS connection.  It has no effect if the DisableTLS parameter
	// is true.
	Certificates []byte

	// Proxy specifies to connect through a SOCKS 5 proxy server.  It may
	// be an empty string if a proxy is not required.
	Proxy string

	// ProxyUser is an optional username to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyUser string

	// ProxyPass is an optional password to use for the proxy server if it
	// requires authentication.  It has no effect if the Proxy parameter
	// is not set.
	ProxyPass string

	// DisableAutoReconnect specifies the client should not automatically
	// try to reconnect to the server when it has been disconnected.
	DisableAutoReconnect bool

	// DisableConnectOnNew specifies that a websocket client connection
	// should not be tried when creating the client with New.  Instead, the
	// client is created and returned unconnected, and Connect must be
	// called manually.
	DisableConnectOnNew bool

	// HTTPPostMode instructs the client to run using multiple independent
	// connections issuing HTTP POST requests instead of using the default
	// of websockets.  Websockets are generally preferred as some of the
	// features of the client such notifications only work with websockets,
	// however, not all servers support the websocket extensions, so this
	// flag can be set to true to use basic HTTP POST requests instead.
	HTTPPostMode bool

	// ExtraHeaders specifies the extra headers when perform request. It's
	// useful when RPC provider need customized headers.
	ExtraHeaders map[string]string

	// EnableBCInfoHacks is an option provided to enable compatibility hacks
	// when connecting to blockchain.info RPC server
	EnableBCInfoHacks bool
	Name              string `mapstructure:"name"`
	ChainID           uint64 `mapstructure:"chain_id"`
	ID                string `mapstructure:"id"`
}

// DefaultConfig returns a configuration populated with default values
func DefaultConfig() []BTCConfig {
	return []BTCConfig{{
		Host:    "http://127.0.0.1:18332",
		User:    "user",
		Pass:    "password",
		Name:    "bitcoin",
		ChainID: 2,
		ID:      "bitcoin-regtest",
	}}
}
func (c *BTCConfig) GetRPCConfig() *rpcclient.ConnConfig {

	return &rpcclient.ConnConfig{
		Host:                 c.Host,
		Endpoint:             c.Endpoint,
		User:                 c.User,
		Pass:                 c.Pass,
		CookiePath:           c.CookiePath,
		Params:               c.Params,
		DisableTLS:           c.DisableTLS,
		Certificates:         c.Certificates,
		Proxy:                c.Proxy,
		ProxyUser:            c.ProxyUser,
		ProxyPass:            c.ProxyPass,
		DisableAutoReconnect: c.DisableAutoReconnect,
		DisableConnectOnNew:  c.DisableConnectOnNew,
		HTTPPostMode:         c.HTTPPostMode,
		ExtraHeaders:         c.ExtraHeaders,
		EnableBCInfoHacks:    c.EnableBCInfoHacks,
	}
}
