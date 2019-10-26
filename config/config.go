package config

import (
	"errors"
	"fmt"
)

type Config interface {
	GetOwner() string
	GetToken() string
}

const configVariableNotSet = "config error, variable for %s not set"

var errNotTokenConfig = errors.New(fmt.Sprintf(configVariableNotSet, "TOKEN"))
var errNotOwnerConfig = errors.New(fmt.Sprintf(configVariableNotSet, "OWNER"))

type config struct {
	token    string
	owner    string
	provider Provider
}

func (c config) GetOwner() string {
	return c.owner
}

func (c config) GetToken() string {
	return c.token
}

func (c *config) read() error {
	var err error = nil

	c.token, err = c.provider.getConfigValue("TOKEN")
	if err == errKeyNotFound {
		return errNotTokenConfig
	} else {
		c.owner, err = c.provider.getConfigValue("OWNER")
		if err == errKeyNotFound {
			return errNotOwnerConfig
		}
	}

	return err
}

func FromProvider(provider Provider) (Config, error) {
	cfg := config{provider: provider}
	err := cfg.read()
	return cfg, err
}
