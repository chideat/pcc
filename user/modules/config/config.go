package config

/*
This file is used to manage global configurations

Author: chideat <chinaxiahaifeng@gmail.com>
*/

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	// config file path
	CONFIG_DEFAULT_PATH = "conf/app.conf"
)

type _config struct {
	Group      uint8             `yaml:"group"`
	Name       string            `yaml:"name"`
	Model      string            `yaml:"model"`
	Version    string            `yaml:"version"`
	HTTPAddr   string            `yaml:"http_addr"`
	RPCAddr    string            `yaml:"rpc_addr"`
	LogPath    string            `yaml:"log_dir"`
	StaticPath string            `yaml:"static_dir"`
	Database   string            `yaml:"database"`
	Caches     map[string]string `yaml:"caches"`
}

func (conf *_config) IsDebug() bool {
	if conf.Model == "debug" || os.Getenv("DEBUG") != "" {
		return true
	}
	return false
}

var Conf _config

func init() {
	configFilePath := CONFIG_DEFAULT_PATH
	// get config file path from cmd line.

	if data, err := ioutil.ReadFile(configFilePath); err == nil {
		if err := yaml.Unmarshal(data, &Conf); err != nil {
			panic(err)
		}

		if Conf.Name == "" {
			panic("Please specify app's name in config file.")
		}

		// check log path
		if logDirInfo, err := os.Stat(Conf.LogPath); err == nil {
			if !logDirInfo.IsDir() {
				panic(fmt.Sprintf("%s is NOT a dir", Conf.LogPath))
			}
		} else {
			panic(err)
		}
	} else {
		panic(err)
	}
}
