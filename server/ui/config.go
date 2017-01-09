package ui

// Config is the UI server configuration
type Config struct {
	BindAddr string `env:"MH2_UI_BIND_ADDR"`
}
