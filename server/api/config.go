package api

// Config is the API server configuration
type Config struct {
	BindAddr string `env:"MH2_API_BIND_ADDR"`
	Backend  string `env:"MH2_BACKEND"`
}
