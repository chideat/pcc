package config

/*
This file is used to manage global configurations

Author: chideat <chinaxiahaifeng@gmail.com>
*/

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const (
	// config file path
	CONFIG_DEFAULT_PATH = "conf/app.conf"
)

var Config struct {
	Name        string            `yaml:"name"`
	Model       string            `yaml:"model"`
	Version     string            `yaml:"version"`
	HttpAddress string            `yaml:"http_address"`
	LogPath     string            `yaml:"log_dir"`
	StaticPath  string            `yaml:"static_dir"`
	Database    string            `yaml:"database"`
	Caches      map[string]string `yaml:"caches"`
}

func init() {
	configFilePath := CONFIG_DEFAULT_PATH
	// get config file path from cmd line.

	if data, err := ioutil.ReadFile(configFilePath); err == nil {
		if err := yaml.Unmarshal(data, &Config); err != nil {
			panic(err)
		}

		if os.Getenv("DEBUG") == "debug" {
			Config.Model = "debug"
		}

		if Config.Name == "" {
			panic("Please specify app's name in config file.")
		}

		// check log path
		if logDirInfo, err := os.Stat(Config.LogPath); err == nil {
			if !logDirInfo.IsDir() {
				panic(fmt.Sprintf("%s is NOT a dir", Config.LogPath))
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
