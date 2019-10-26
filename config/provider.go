package config

import (
	"errors"
	"os"
)

var errKeyNotFound = errors.New("config error, key not found")

type Provider interface {
	getConfigValue(key string) (string, error)
}

const environmentVariablesBaseKey = "CECIBOT_"

type EnvironmentVariableProvider struct {
	baseKey string
}

func (e EnvironmentVariableProvider) getConfigValue(key string) (string, error) {
	var value = os.Getenv(e.baseKey + key)

	if value == "" {
		return "", errKeyNotFound
	}

	return value, nil
}

func EnvironmentVariables() Provider {
	return EnvironmentVariableProvider{environmentVariablesBaseKey}
}
