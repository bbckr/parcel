package structs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/bbckr/parcel/helpers"
	"github.com/hashicorp/go-multierror"
)

const (
	templateDirectoryPath = "templates"
	valuesFilePath        = "values.yaml"
	staticDirectoryPath   = "static"
	manifestFilePath      = "manifest.yaml"
)

func NewParcel(cfg *Config, owner, name, version string) (*Parcel, error) {
	parcel := &Parcel{
		Owner:            owner,
		Name:             name,
		Version:          version,
		InstallDirectory: cfg.ParcelInstallDirectory,
	}

	path, err := cfg.Index.PathFrom(parcel.ID())
	if err != nil {
		return nil, fmt.Errorf("Error reading parcel %s/%s version %s from index: %s", owner, name, version, err)
	}

	source, _ := helpers.DecodeBase64(filepath.Base(path))
	parcel.Source = source

	return parcel, nil
}

func NewParcelFromSource(cfg *Config, source string) *Parcel {
	parcel := &Parcel{
		Source:           source,
		InstallDirectory: cfg.ParcelInstallDirectory,
	}

	return parcel
}

type Parcel struct {
	Name             string
	Owner            string
	Version          string
	Description      string
	InstallDirectory string
	Source           string

	Templates []*ParcelTemplate
	Values    map[interface{}]interface{}
}

func (p *Parcel) ID() string {
	return helpers.JoinNonEmptyStrings([]string{p.Owner, p.Name, p.Version}, "/")
}

func (p *Parcel) Ref() string {
	return helpers.EncodeBase64FromArray([]string{p.Owner, p.Name, p.Version})
}

func (p *Parcel) DirectoryName() string {
	return helpers.EncodeBase64(p.Source)
}

func (p *Parcel) InstallPath() string {
	return fmt.Sprintf("%s/%s", p.InstallDirectory, p.DirectoryName())
}

func (p *Parcel) TemplateDirectory() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), templateDirectoryPath)
}

func (p *Parcel) StaticDirectory() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), staticDirectoryPath)
}

func (p *Parcel) ValuesPath() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), valuesFilePath)
}

func (p *Parcel) ManifestPath() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), manifestFilePath)
}

func (p *Parcel) IndexEntry() *Entry {
	return &Entry{
		Name:        p.Name,
		Owner:       p.Owner,
		Version:     p.Version,
		Description: p.Description,
		Path:        p.InstallPath(),
		Source:      p.Source,
	}
}

func (p *Parcel) IsValid() error {
	var result error

	if p.Name == "" {
		result = multierror.Append(result, fmt.Errorf("Name -> must not be empty"))
	}

	if p.Owner == "" {
		result = multierror.Append(result, fmt.Errorf("Owner -> must not be empty"))
	}

	if p.Version == "" {
		result = multierror.Append(result, fmt.Errorf("Version -> must not be empty"))
	}

	if p.InstallDirectory == "" {
		result = multierror.Append(result, fmt.Errorf("InstallDirectory -> must not be empty"))
	}

	if info, err := os.Stat(p.InstallPath()); os.IsNotExist(err) || !info.IsDir() {
		result = multierror.Append(result, fmt.Errorf("InstallPath -> must exist at path %s", p.InstallPath()))
	}

	if info, err := os.Stat(p.ValuesPath()); os.IsNotExist(err) || info.IsDir() {
		result = multierror.Append(result, fmt.Errorf("ValuesPath -> %s must exist", valuesFilePath))
	}

	if info, err := os.Stat(p.TemplateDirectory()); os.IsNotExist(err) || !info.IsDir() {
		result = multierror.Append(result, fmt.Errorf("TemplateDirectory -> %s must exist", templateDirectoryPath))
	}

	return result
}

type ParcelTemplate struct {
	Path     string
	Template *template.Template
}

func (tmpl *ParcelTemplate) Render(values map[interface{}]interface{}) (bytes.Buffer, error) {
	t := tmpl.Template

	var renderedOutput bytes.Buffer
	if err := t.Execute(&renderedOutput, values); err != nil {
		return renderedOutput, fmt.Errorf("Error rendering template: %s", err)
	}
	return renderedOutput, nil
}
