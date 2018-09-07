// config.go
package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Configuration struct {
	RemoteHosts []RemoteHostAlias
}

var configFile = filepath.Join(os.Getenv("HOME"), ".config/copy.conf")

func ReadConfig(filename string) (*Configuration, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	config := &Configuration{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func WriteConfig(filename string, config *Configuration) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}
	return nil
}
