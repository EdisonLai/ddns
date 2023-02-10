package config

var GConf *DdnsConfig

type DomainConfig struct {
	Domain    string `toml:"domain"`
	SubDomain string `toml:"subDomain"`
	Line      string `toml:"line"`
	CheckTime int    `toml:"checkTime"`
}

type ProviderConfig struct {
	Provider  string `toml:"provider"`
	SecretId  string `toml:"secretId"`
	SecretKey string `toml:"secretKey"`
}

type DdnsConfig struct {
	Provider ProviderConfig `toml:"provider"`
	Domain   DomainConfig   `toml:"domain"`
}

func InitConfig(configPath string) (err error) {
	return
}
