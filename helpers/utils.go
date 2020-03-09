package helpers

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

func GetEnvWithDefault(key, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func GetEnvBool(key, defaultValue string) bool {
	val, err := strconv.ParseBool(GetEnvWithDefault(key, defaultValue))
	if err != nil {
		log.Fatal(err)
	}

	return val
}

func FindValueFromKeyAsPrefix(s string, m map[string]string) (string, bool) {
	for k, v := range m {
		if strings.HasPrefix(s, k) {
			return v, true
		}
	}
	return "", false
}

func MustCompileRegexSubmatch(s, regexStr string) map[string]string {
	re := regexp.MustCompile(regexStr)
	match := re.FindStringSubmatch(s)
	submatchMap := make(map[string]string)
	for index, m := range match {
		submatchMap[re.SubexpNames()[index]] = m
	}
	return submatchMap
}

func JoinNonEmptyStrings(s []string, sep string) string {
	var segments []string
	for _, s := range s {
		if strings.TrimSpace(s) != "" {
			segments = append(segments, s)
		}
	}
	return strings.Join(segments, sep)
}

func LoadYamlFromPath(path string) (map[interface{}]interface{}, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	result := make(map[interface{}]interface{})
	if err = yaml.Unmarshal(f, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0751)
}
