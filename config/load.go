package config

import (
	"io/ioutil"

	"github.com/mcuadros/go-defaults"
	"github.com/rolldever/go-json5"
)

func Load(jsonPath string) (*Config, error) {
	bDoc, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err = json5.Unmarshal(bDoc, &config); err != nil {
		return nil, err
	}
	defaults.SetDefaults(&config)
	return &config, nil
}
