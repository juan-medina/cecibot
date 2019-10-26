package config

import (
	"os"
	"reflect"
	"testing"
)

func TestEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name string
		want Provider
	}{
		{
			"we should get a valid EnvironmentVariableProvider provider",
			EnvironmentVariableProvider{baseKey: environmentVariablesBaseKey},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EnvironmentVariables(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnvironmentVariables() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEnvironmentsVariableProvider_getConfigValue(t *testing.T) {
	type fields struct {
		baseKey string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			"we should get a valid value",
			fields{baseKey: "TEST_ENV_"},
			args{key: "VALUE"},
			"A VALUE",
			false,
		},
		{
			"we should get an error",
			fields{baseKey: "TEST_ENV_"},
			args{key: "VALUE_DO_NOT_EXIST"},
			"",
			true,
		},
	}
	_ = os.Setenv("TEST_ENV_VALUE", "A VALUE")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := EnvironmentVariableProvider{
				baseKey: tt.fields.baseKey,
			}
			got, err := e.getConfigValue(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("getConfigValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getConfigValue() got = %v, want %v", got, tt.want)
				return
			}
		})
	}
	_ = os.Unsetenv("TEST_ENV_VALUE")
}
