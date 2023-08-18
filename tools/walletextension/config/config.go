package config

// Config contains the configuration required by the WalletExtension.
type Config struct {
	WalletExtensionHost     string
	WalletExtensionPortHTTP int
	WalletExtensionPortWS   int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
	LogPath                 string
	DBPathOverride          string // Overrides the database file location. Used in tests.
	Hosted                  string
	VerboseFlag             bool
}
