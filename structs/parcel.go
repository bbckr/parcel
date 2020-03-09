package structs

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"

	"github.com/bbckr/parcel/helpers"
	"github.com/hashicorp/go-multierror"
)

const (
	templateDirectoryPath = "templates"
	valuesFilePath        = "values.yaml"
	assetsDirectoryPath   = "assets"
)

var (
	protocolPatterns = map[string]string{
		"git::ssh":  "git::ssh://(?P<Host>.*)/(?P<Owner>.*)/(?P<Repository>.*).git(//(?P<SubDir>[-\\w/]*))?([?]ref=(?P<Ref>.*))?",
		"git::file": "git::file://(?P<BaseDir>.*)/(?P<Repository>.*)/.git(//(?P<SubDir>[-\\w/]*))?([?]ref=(?P<Ref>.*))?",
	}
)

func NewParcel(cfg *Config, source string) (*Parcel, error) {
	pattern, ok := helpers.FindValueFromKeyAsPrefix(source, protocolPatterns)
	if !ok {
		return nil, fmt.Errorf("Invalid source protocol: must be %s", reflect.ValueOf(protocolPatterns).MapKeys())
	}

	parsedSource := helpers.MustCompileRegexSubmatch(source, pattern)
	if len(parsedSource) == 0 {
		return nil, fmt.Errorf("Could not parse source: %s", source)
	}

	name := parsedSource["SubDir"]
	if name != "" {
		name = filepath.Base(parsedSource["SubDir"])
	}
	ref := strings.Replace(strings.Replace(parsedSource["Ref"], "?ref=", "", 1), fmt.Sprintf("%s-", name), "", 1)
	parcel := &Parcel{
		Name:             name,
		Repository:       parsedSource["Repository"],
		Tag:              ref,
		InstallDirectory: cfg.ParcelInstallDirectory,
		Owner:            parsedSource["Owner"],
	}

	return parcel, nil
}

type Parcel struct {
	Name             string
	Repository       string
	Owner            string
	Tag              string
	InstallDirectory string

	Templates []*ParcelTemplate
	Values    map[interface{}]interface{}
}

func (p *Parcel) ID() string {
	return helpers.JoinNonEmptyStrings([]string{p.Owner, p.Repository, p.Name, p.Tag}, "/")
}

func (p *Parcel) DirectoryName() string {
	return helpers.JoinNonEmptyStrings([]string{p.Owner, p.Repository, p.Name, p.Tag}, "-")
}

func (p *Parcel) InstallPath() string {
	return fmt.Sprintf("%s/%s", p.InstallDirectory, p.DirectoryName())
}

func (p *Parcel) TemplateDirectory() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), templateDirectoryPath)
}

func (p *Parcel) AssetsDirectory() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), assetsDirectoryPath)
}

func (p *Parcel) ValuesPath() string {
	return fmt.Sprintf("%s/%s", p.InstallPath(), valuesFilePath)
}

func (p *Parcel) Version() string {
	if len(p.Tag) == 0 {
		return "latest"
	}

	return strings.Replace(p.Tag, p.Name, "", 1)
}

func (p *Parcel) IsValid() error {
	var result error

	if p.Name == "" {
		result = multierror.Append(result, fmt.Errorf("Name -> must not be empty"))
	}

	if p.Repository == "" {
		result = multierror.Append(result, fmt.Errorf("Repository -> must not be empty"))
	}

	if p.Tag == "" {
		result = multierror.Append(result, fmt.Errorf("Tag -> must not be empty"))
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
	var renderedOutput bytes.Buffer
	if err := tmpl.Template.Execute(&renderedOutput, values); err != nil {
		return renderedOutput, fmt.Errorf("Error rendering template: %s", err)
	}
	return renderedOutput, nil
}
