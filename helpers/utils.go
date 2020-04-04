package helpers

import (
	"encoding/base64"
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

func ContainsAnyPrefix(s string, arr []string) bool {
	for _, v := range arr {
		if strings.HasPrefix(s, v) {
			return true
		}
	}
	return false
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

func LoadYamlEnsurePath(path string, out interface{}) error {
	f, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(bytes, out); err != nil {
		return err
	}

	return nil
}

func LoadYamlFromPath(path string, out interface{}) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(f, out); err != nil {
		return err
	}

	return nil
}

func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0751)
}

func EncodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func EncodeBase64FromArray(args []string) string {
	s := strings.Join(args, "")
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func DecodeBase64(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
