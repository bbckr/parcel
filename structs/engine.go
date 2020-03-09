package structs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bbckr/parcel/helpers"
	log "github.com/sirupsen/logrus"
)

func NewEngineFromConfig(cfg *Config) *Engine {
	return &Engine{
		cfg: cfg,
	}
}

type Engine struct {
	cfg *Config
}

type RenderOptions struct {
	ValuesFilePath  string
	OutputDirectory string
}

func (e *Engine) Render(parcel *Parcel, opts *RenderOptions, logger *log.Entry) error {
	var loadedValues map[interface{}]interface{}
	if opts.ValuesFilePath != "" {
		logger.WithField("path", opts.ValuesFilePath).Info("Loading provided values")
		values, err := helpers.LoadYamlFromPath(opts.ValuesFilePath)
		if err != nil {
			return err
		}
		loadedValues = values
	}

	var mergedValues = make(map[interface{}]interface{})
	for k, v := range parcel.Values {
		// copy parcel values as base
		mergedValues[k] = v
	}
	for k, v := range loadedValues {
		// override parcel values
		mergedValues[k] = v
	}
	data := map[interface{}]interface{}{
		"Values": mergedValues,
		"Parcel": map[string]string{
			"name":       parcel.Name,
			"version":    parcel.Version(),
			"repository": parcel.Repository,
			"tag":        parcel.Tag,
		},
		"Assets": map[string]string{
			"path": parcel.AssetsDirectory(),
		},
	}

	for _, t := range parcel.Templates {
		templateBasePath := filepath.Base(t.Path)
		logger.WithField("template", templateBasePath).Info("Rendering template to output directory")
		outputPath := fmt.Sprintf("%s/%s", opts.OutputDirectory, templateBasePath)
		buffer, err := t.Render(data)
		err = ioutil.WriteFile(outputPath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) Load(source string) (*Parcel, error) {
	parcel, err := NewParcel(e.cfg, source)
	if err != nil {
		return nil, err
	}

	if err := parcel.IsValid(); err != nil {
		return nil, fmt.Errorf("Invalid parcel: %s", err)
	}

	values, err := helpers.LoadYamlFromPath(parcel.ValuesPath())
	if err != nil {
		return nil, fmt.Errorf("Unable to load parcel: %s", err)
	}
	parcel.Values = values

	templateFiles, err := ioutil.ReadDir(parcel.TemplateDirectory())
	if err != nil {
		return nil, fmt.Errorf("Error reading template dir: %s", err)
	}

	templates := make([]*ParcelTemplate, 0)
	for _, f := range templateFiles {
		templatePath := fmt.Sprintf("%s/%s", parcel.TemplateDirectory(), f.Name())

		t, err := template.ParseFiles(templatePath)
		if err != nil {
			return nil, err
		}
		templates = append(templates, &ParcelTemplate{
			Path:     templatePath,
			Template: t,
		})
	}
	parcel.Templates = templates
	return parcel, nil
}

type PullOptions struct {
	Force bool
}

func (e *Engine) Pull(source string, opts *PullOptions, logger *log.Entry) error {
	parcel, err := NewParcel(e.cfg, source)
	if err != nil {
		return err
	}

	installPath := parcel.InstallPath()

	if _, err := os.Stat(installPath); !os.IsNotExist(err) && !opts.Force {
		logger.WithField("force", opts.Force).Info("Parcel found at install path, not pulling")
		return nil
	}

	logger.WithField("path", installPath).Info("Installing parcel to path")
	if _, err := helpers.GetWithGoGetter(installPath, source); err != nil {
		return fmt.Errorf("Unable to pull parcel from source: %s", err)
	}
	return nil
}
