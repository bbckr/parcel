package structs

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/bbckr/parcel/helpers"

	"gopkg.in/yaml.v2"
)

const (
	indexFilePath = "index.yaml"
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
	Index                   *Index
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

func NewConfig() (*Config, error) {
	cfg := DefaultConfig()
	cfg.initializeFromFields(defaultFields)

	err := loadIndex(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func loadIndex(cfg *Config) error {
	index := &Index{
		Path:    fmt.Sprintf("%s/%s", cfg.ParcelInstallDirectory, indexFilePath),
		Entries: make(map[string]*Entry),
	}

	if err := helpers.LoadYamlEnsurePath(index.Path, &index); err != nil {
		return fmt.Errorf("Unable to load index: %s", err)
	}

	cfg.Index = index
	return nil
}

type Entry struct {
	Name        string `yaml:"name"`
	Owner       string `yaml:"owner"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Path        string `yaml:"path"`
	Source      string `yaml:"source"`
}

type Index struct {
	Path    string            `yaml:"-"`
	Entries map[string]*Entry `yaml:"entries"`
}

func (i *Index) PathFrom(k string) (string, error) {
	entry, ok := i.Entries[k]
	if !ok {
		return "", fmt.Errorf("Entry does not exist")
	}
	return entry.Path, nil
}

func (i *Index) AddEntry(key string, entry *Entry) {
	i.Entries[key] = entry
}

func (i *Index) Persist() error {
	data, err := yaml.Marshal(i)
	if err != nil {
		return fmt.Errorf("Could not save index: %s", err)
	}

	if err = ioutil.WriteFile(i.Path, data, 0644); err != nil {
		return fmt.Errorf("Could not write to index: %s", err)
	}

	return nil
}
