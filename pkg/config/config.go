package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerCfg `mapstructure:"server"`
}

type ServerCfg struct {
	Host          string        `mapstructure:"host"`
	Port          int           `mapstructure:"port"`
	WriteTimeout  time.Duration `mapstructure:"writetimeout"`
	ReadTimeout   time.Duration `mapstructure:"readtimeout"`
	MaxHeaderByte int           `mapstructure:"maxheaderbyte"`
}

func Init(path, file string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(file)

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("error  read in config")
		return nil, err
	}

	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Println("error  config unmarshal")
		return nil, err
	}
	return &cfg, nil
}
