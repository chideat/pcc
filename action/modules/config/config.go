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

var Conf struct {
	Name        string `yaml:"name"`
	Model       string `yaml:"model"`
	Version     string `yaml:"version"`
	HttpAddress string `yaml:"http_address"`
	LogPath     string `yaml:"log_dir"`
	Database    string `yaml:"database"`
	Queue       struct {
		NsqdAddress    string `yaml:"nsqd_address"`
		LookupdAddress string `yaml:"lookupd_address"`
	} `yaml:"queue"`
}

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

		Conf.Model = os.Getenv("DEBUG")

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
