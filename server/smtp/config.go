package smtp

// Config is the SMTP server configuration
type Config struct {
	BindAddr            string `env:"MH2_SMTP_BIND_ADDR"`
	Hostname            string `env:"MH2_SMTP_HOSTNAME"`
	LogData             bool   `env:"MH2_LOG_DATA"`
	LogProto            bool   `env:"MH2_LOG_PROTO"`
	Backend             string `env:"MH2_BACKEND"`
	RecordSessionData   bool   `env:"MH2_RECORD_SESSION_DATA"`
	RecordSessionProto  bool   `env:"MH2_RECORD_SESSION_PROTO"`
	RecordSessionEvents bool   `env:"MH2_RECORD_SESSION_EVENTS"`
	Jim                 JimConfig
}

// JimConfig is the Jim chaos monkey configuration
type JimConfig struct {
	Enabled               bool    `env:"JIM_ENABLED"`
	DisconnectChance      float64 `env:"JIM_DISCONNECT_CHANCE"`
	AcceptChance          float64 `env:"JIM_ACCEPT_CHANCE"`
	LinkSpeedAffect       float64 `env:"JIM_LINKSPEED_AFFECT"`
	LinkSpeedMin          float64 `env:"JIM_LINKSPEED_MIN"`
	LinkSpeedMax          float64 `env:"JIM_LINKSPEED_MAX"`
	RejectSenderChance    float64 `env:"JIM_REJECT_SENDER_CHANCE"`
	RejectRecipientChance float64 `env:"JIM_REJECT_RECIPIENT_CHANCE"`
	RejectAuthChance      float64 `env:"JIM_REJECT_AUTH_CHANCE"`
}
