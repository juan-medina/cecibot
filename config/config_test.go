package config

import (
	"testing"
)

type FakeProvider struct {
}

func (f FakeProvider) getConfigValue(key string) (string, error) {
	if key == "TOKEN" {
		return "fake token", nil
	}
	return "fake owner", nil
}

type ErrorOnKeyProvider struct {
	key string
}

func (f ErrorOnKeyProvider) getConfigValue(key string) (string, error) {
	if f.key == key {
		return "", errKeyNotFound
	}
	return "fake value", nil
}

func TestFromProvider(t *testing.T) {
	var tests = []struct {
		name      string
		provider Provider
		wantToken string
		wantOwner string
		wantErr   error
	}{
		{
			"we should get the fake values",
			FakeProvider{},
			"fake token",
			"fake owner",
			nil,
		},
		{
			"we should get an token error",
			ErrorOnKeyProvider{key: "TOKEN"},
			"fake token",
			"fake owner",
			errNotTokenConfig,
		},
		{
			"we should get an owner error",
			ErrorOnKeyProvider{key: "OWNER"},
			"fake token",
			"fake owner",
			errNotOwnerConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FromProvider(tt.provider)
			if err != tt.wantErr {
				t.Errorf("FromProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if got == nil {
					t.Errorf("FromProvider() got nil")
					return
				}
				gotToken := got.GetToken()
				if gotToken != tt.wantToken {
					t.Errorf("FromProvider() got token = %q, want %q", gotToken, tt.wantToken)
					return
				}
				gotOwner := got.GetOwner()
				if gotOwner != tt.wantOwner {
					t.Errorf("FromProvider() got owner = %q, want %q", gotOwner, tt.wantOwner)
					return
				}
			}
		})
	}
}
func Test_config_read(t *testing.T) {
	tests := []struct {
		name     string
		provider Provider
		wantErr  error
	}{
		{
			"we should not get an error",
			FakeProvider{},
			nil,
		},
		{
			"we should get an error not token",
			ErrorOnKeyProvider{key: "TOKEN"},
			errNotTokenConfig,
		},
		{
			"we should get an error not owner",
			ErrorOnKeyProvider{key: "OWNER"},
			errNotOwnerConfig,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &config{provider: tt.provider}
			if err := c.read(); err != tt.wantErr {
				t.Errorf("read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
