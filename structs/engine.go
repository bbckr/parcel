package structs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		if err := helpers.LoadYamlFromPath(opts.ValuesFilePath, loadedValues); err != nil {
			return err
		}
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
		"Meta": map[string]string{
			"name":    parcel.Name,
			"owner":   parcel.Owner,
			"version": parcel.Version,
		},
		"Static": map[string]string{
			"path": parcel.StaticDirectory(),
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

func (e *Engine) Load(owner, name, version string) (*Parcel, error) {
	parcel, err := NewParcel(e.cfg, owner, name, version)
	if err != nil {
		return nil, fmt.Errorf("Invalid parcel: %s,", err)
	}

	if err := parcel.IsValid(); err != nil {
		return nil, fmt.Errorf("Invalid parcel: %s", err)
	}

	var values map[interface{}]interface{}
	if err := helpers.LoadYamlFromPath(parcel.ValuesPath(), values); err != nil {
		return nil, fmt.Errorf("Unable to load parcel: %s", err)
	}
	parcel.Values = values

	templateFiles, err := ioutil.ReadDir(parcel.TemplateDirectory())
	if err != nil {
		return nil, fmt.Errorf("Error reading template dir: %s", err)
	}

	templates := make([]*ParcelTemplate, 0)
	funcs := funcMap()
	for _, f := range templateFiles {
		templatePath := fmt.Sprintf("%s/%s", parcel.TemplateDirectory(), f.Name())

		t, err := template.New(f.Name()).Funcs(funcs).ParseFiles(templatePath)
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
	var (
		protocolPrefixes = []string{"git::ssh", "git::file"}
	)

	ok := helpers.ContainsAnyPrefix(source, protocolPrefixes)
	if !ok {
		return fmt.Errorf("Invalid source protocol: must be of prefix %s or ref to [OWNER]/[NAME]", protocolPrefixes)
	}

	parcel := NewParcelFromSource(e.cfg, source)
	installPath := parcel.InstallPath()

	if _, err := os.Stat(installPath); !os.IsNotExist(err) && !opts.Force {
		logger.WithField("force", opts.Force).Info("Parcel found at install path, not pulling")
		return nil
	}

	logger.WithField("path", installPath).Info("Installing parcel to path")
	if _, err := helpers.GetWithGoGetter(installPath, source); err != nil {
		return fmt.Errorf("Unable to pull parcel from source: %s", err)
	}

	var manifest map[string]string
	if err := helpers.LoadYamlFromPath(parcel.ManifestPath(), manifest); err != nil {
		return fmt.Errorf("Unable to load parcel manifest: %s", err)
	}

	entry := &Entry{
		Name:        manifest["name"],
		Owner:       manifest["owner"],
		Version:     manifest["version"],
		Description: manifest["description"],
		Path:        installPath,
		Source:      source,
	}

	e.cfg.Index.AddEntry(parcel.Ref(), entry)
	if err := e.cfg.Index.Persist(); err != nil {
		return err
	}

	return nil
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"toJSON": func(i interface{}) string {
			data, err := json.Marshal(i)
			if err != nil {
				return ""
			}
			return string(data)
		},
		"replace": func(old, new string, n int, s string) string {
			return strings.Replace(s, old, new, n)
		},
		"replaceAll": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"quote": func(s string) string {
			return fmt.Sprintf(`"%s"`, s)
		},
		"b64encode": func(s string) string {
			return base64.StdEncoding.EncodeToString([]byte(s))
		},
		"b64decode": func(s string) string {
			data, err := base64.StdEncoding.DecodeString(s)
			if err != nil {
				return ""
			}
			return string(data)
		},
	}
}
