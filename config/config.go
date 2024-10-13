package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the configuration values
type Config struct {
	Env      string `mapstructure:"env"`
	HttpPort string `mapstructure:"http_port"`
	GrpcPort string `mapstructure:"grpc_port"`
	DB       struct {
		User     string `mapstructure:"user"`
		Name     string `mapstructure:"name"`
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"password"`
		Port     string `mapstructure:"port"`
	} `mapstructure:"db"`
}

// LoadConfig loads the configuration using Viper
func LoadConfig() (Config, error) {
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.SetConfigType("yml")                   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("../")                   // path to look for the config file in
	viper.AddConfigPath("$HOME/.grpc-todo-list") // call multiple times to add many search paths
	viper.AddConfigPath(".")                     // optionally look for config in the working directory

	viper.AutomaticEnv()                                   // automatically try to load missing keys from from env vars
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`)) // replace . from nested variables to _

	var config Config

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		return config, fmt.Errorf("error reading config file: %v", err)
	}

	// Unmarshal the config into the Config struct
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("unable to decode into struct: %v", err)
	}

	return config, nil
}
