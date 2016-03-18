package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func parseConfig() tomlConfig {
	var config tomlConfig
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
	}
  return config
}

// tomlConfig is an example
type tomlConfig struct {
	Username string
	Password string
	Email    string
}
