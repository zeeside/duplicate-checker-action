package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	directory         string
	fileExtension     string
	contentRegexStr   string
	maxFileSizeKB     int
	maxFilesToProcess int
	ignoreExtMap      map[string]struct{}
	ignoreFilesMap    map[string]struct{}
	ignorePathsMap    map[string]struct{}
)

type Config struct {
	CheckName           string
	DirectoryContext    string
	MaxFileSizeKB       int
	MaxFilesToProcess   int
	FileExtension       string
	FileNameRegexString string
	ContentRegex        *regexp.Regexp
	ExcludedExtensions  map[string]struct{}
	LogLevel            string
	IgnoreFiles         map[string]struct{}
	IgnorePaths         map[string]struct{}
}

func NewConfig() (*Config, error) {
	// initially assigning to global variables to allow
	// checking entered values with a call to validate()
	directory = getEnvStr("INPUT_DIRECTORY_SCOPE", ".")
	fileExtension = os.Getenv("INPUT_CHECK_FILE_EXTENSION")
	contentRegexStr = os.Getenv("INPUT_CONTENT_REGEX")
	maxFilesToProcess = getEnvInt("INPUT_MAX_FILES_TO_PROCESS", 500)
	maxFileSizeKB = getEnvInt("INPUT_MAX_FILE_SIZE_BYTES", 200000)

	ignoreExtList := strings.Split(getEnvStr("INPUT_EXCLUDED_EXTENSIONS", ""), ",")
	ignoreExtMap = buildExcludedExtensionsMap(ignoreExtList)

	ignoreFileList := strings.Split(getEnvStr("INPUT_IGNORE_FILES", ""), ",")
	ignoreFilesMap = buildIgnoredFilesMap(ignoreFileList)

	ignorePathList := strings.Split(getEnvStr("INPUT_IGNORE_PATHS_CONTAINING", ""), ",")
	ignorePathsMap = buildIgnoredPathsMap(ignorePathList)

	_, err := validate()
	if err != nil {
		return nil, err
	}

	contentRegex, _ := regexp.Compile(contentRegexStr)
	c := &Config{
		DirectoryContext:   directory,
		MaxFileSizeKB:      maxFileSizeKB,
		MaxFilesToProcess:  maxFilesToProcess,
		FileExtension:      toValidExtensionFormat(fileExtension),
		ContentRegex:       contentRegex,
		LogLevel:           getEnvStr("INPUT_LOG_LEVEL", "info"),
		IgnoreFiles:        ignoreFilesMap,
		IgnorePaths:        ignorePathsMap,
		ExcludedExtensions: ignoreExtMap,
	}
	return c, nil
}

func buildExcludedExtensionsMap(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for i := 0; i < len(list); i++ {
		validExt := toValidExtensionFormat(list[i])
		_, ok := m[validExt]
		if !ok {
			m[validExt] = struct{}{}
		}
	}
	return m
}

func buildIgnoredFilesMap(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for i := 0; i < len(list); i++ {
		d := list[i]
		_, ok := m[d]
		if !ok {
			m[d] = struct{}{}
		}
	}
	return m
}

func buildIgnoredPathsMap(list []string) map[string]struct{} {
	m := make(map[string]struct{}, len(list))
	for i := 0; i < len(list); i++ {
		d := list[i]
		_, ok := m[d]
		if !ok {
			m[d] = struct{}{}
		}
	}
	return m
}

func toValidExtensionFormat(input string) string {
	if len(input) > 0 && input[0:1] != "." {
		return fmt.Sprintf(".%s", fileExtension)
	}

	return input
}

func validate() (bool, error) {
	if len(fileExtension) == 0 || len(fileExtension) == 1 && fileExtension == "." {
		return false, fmt.Errorf("invalid file extension")
	}

	_, err := regexp.Compile(contentRegexStr)
	if err != nil {
		return false, fmt.Errorf("error compiling content regex: %v", err)
	}

	files, err := os.Open(directory)
	if err != nil {
		return false, fmt.Errorf("error opening directory: %v", err)
	}
	files.Close()

	return true, nil
}
