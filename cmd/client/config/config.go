package config

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

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

type EIPServiceConfig struct {
	Type   string `toml:"type"`
	Server string `toml:"server"`
}

type DdnsConfig struct {
	Provider  ProviderConfig   `toml:"provider"`
	Domain    DomainConfig     `toml:"domain"`
	EIPMethod EIPServiceConfig `toml:"eipGetServer"`
}

func InitConfig() (err error) {
	var configPath string
	flag.StringVar(&configPath, "config", "", "Path to a ddnsconfig.")

	flag.Parse()

	if configPath == "" {
		panic("invalid config file name")
	}
	file, err := os.OpenFile(configPath, os.O_RDONLY, 0666)
	if err != nil {
		file, err = os.OpenFile(path.Join("..", configPath), os.O_RDWR, 0666)
		if err != nil {
			panic(err.Error())
		}
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err.Error())
	}

	_, err = toml.Decode(string(content), &GConf)
	if err != nil {
		panic(err.Error())
	}

	if err = checkConfValidation(); err != nil {
		panic(err.Error())
	}

	return
}

func checkConfValidation() error {
	if GConf.Provider.Provider == "" {
		return fmt.Errorf("dns service provider is empty")
	}
	if GConf.Domain.Domain == "" {
		return fmt.Errorf("domain is empty")
	}
	if GConf.Domain.CheckTime == 0 {
		GConf.Domain.CheckTime = 30
	}
	if GConf.EIPMethod.Type == "" {
		GConf.EIPMethod.Type = "http"
	}
	if GConf.EIPMethod.Type == "http" && GConf.EIPMethod.Server == "" {
		return fmt.Errorf("eip service address is empty")
	}
	return nil
}
