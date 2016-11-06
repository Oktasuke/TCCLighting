package models

func NewConfig() Config {
	return Config{}
}

type Config struct {
	Server   Server
	Facebook Facebook
}

type Server struct {
	Listen_port string `toml:"listen_port"`
	Gin_mode    string `toml:"gin_mode"`
}

type Facebook struct {
	App_id       string `toml:"app_id"`
	App_secret   string `toml:"app_secret"`
	Verify_token string `toml:"verify_token"`
}
