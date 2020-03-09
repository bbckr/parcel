package structs

import (
	"reflect"

	"github.com/bbckr/parcel/helpers"
)

var defaultFields = []Field{
	Field{
		Name:                "ParcelInstallDirectory",
		EnvironmentVariable: "PARCEL_INSTALL_DIR",
		DefaultValue:        ".parcel",
	},
	Field{
		Name:                "LogLevel",
		EnvironmentVariable: "LOG_LEVEL",
	},
	Field{
		Name:                "LogFormat",
		EnvironmentVariable: "LOG_FORMAT",
	},
	Field{
		Name:                "Debug",
		EnvironmentVariable: "DEBUG",
	},
}

type Field struct {
	Name                string
	EnvironmentVariable string
	DefaultValue        string
}

type Config struct {
	ParcelInstallDirectory  string
	ParcelTemplateDirectory string
	LogLevel                string
	LogFormat               string
	Debug                   bool
}

func (cfg *Config) initializeFromFields(optionalFields []Field) {
	for _, df := range defaultFields {
		field := reflect.ValueOf(cfg).Elem().FieldByName(df.Name)

		if field.Type().Kind() == reflect.Bool {
			value := helpers.GetEnvBool(df.EnvironmentVariable, df.DefaultValue)
			field.SetBool(value)
			return
		}

		value := helpers.GetEnvWithDefault(df.EnvironmentVariable, df.DefaultValue)
		field.SetString(value)
		return
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func NewConfig() *Config {
	cfg := DefaultConfig()
	cfg.initializeFromFields(defaultFields)
	return cfg
}
