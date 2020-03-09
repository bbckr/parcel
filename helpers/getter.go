package helpers

import (
	"path/filepath"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-getter"
)

var goGetterGetters = map[string]getter.Getter{
	"git": new(getter.GitGetter),
}

var goGetterDetectors = []getter.Detector{
	new(getter.GitHubDetector),
	new(getter.GitDetector),
}

var goGetterNoDetectors = []getter.Detector{}

var goGetterNoDecompressors = map[string]getter.Decompressor{}

var getterHTTPClient = cleanhttp.DefaultClient()

var getterHTTPGetter = &getter.HttpGetter{
	Client: getterHTTPClient,
	Netrc:  true,
}

func GetWithGoGetter(outputDir, source string) (string, error) {
	client := getter.Client{
		Src: source,
		Dst: outputDir,
		Pwd: outputDir,

		Mode: getter.ClientModeDir,

		Detectors:     goGetterNoDetectors,
		Decompressors: goGetterNoDecompressors,
		Getters:       goGetterGetters,
	}

	err := client.Get()
	if err != nil {
		return "", err
	}

	return filepath.Clean(outputDir), nil
}
