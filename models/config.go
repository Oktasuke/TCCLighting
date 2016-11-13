package models

func NewConfig() Config {
	return Config{}
}

type Config struct {
	ServerInfo  	serverInfo
	FacebookInfo 	facebookInfo
	WeMoInfo	WeMoInfo
	ShopInfo	ShopInfo
}

type serverInfo struct {
	ListenPort string `toml:"listen_port"`
	GinMode    string `toml:"gin_mode"`
}

type facebookInfo struct {
	VerifyToken string `toml:"verify_token"`
}

type ShopInfo struct {
	OpeningHour string `toml:"opening_time"`
	ClosingTime string `toml:"closing_time"`
}
type WeMoInfo struct {
	Location string `toml:"location"`
	Port	string `toml:"port"`
}