package git

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/slon/shad-go/gitfame/configs"
)

type Language struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Extension []string `json:"extensions"`
}

func SelectByExtensions(files, fileExtensions, languages []string) ([]string, error) {
	langExtensions, err := selectLanguages(languages)
	if err != nil {
		return nil, err
	}

	extensions := mergeExtensions(fileExtensions, langExtensions)

	return selectByExtensions(files, extensions), nil
}

func selectLanguages(languages []string) ([]string, error) {
	var jsonLanguages []Language

	if err := json.Unmarshal(configs.LanguageExtensions, &jsonLanguages); err != nil {
		return nil, fmt.Errorf("unmarshal language_extensions.json: %v", err)
	}

	m := make(map[string]Language)

	for _, lang := range jsonLanguages {
		m[strings.ToLower(lang.Name)] = lang
	}

	var extension []string

	for _, lang := range languages {
		if l, ok := m[strings.ToLower(lang)]; ok {
			extension = append(extension, l.Extension...)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "unknown language: %q\n", lang)
		}
	}

	return extension, nil
}

func mergeExtensions(fileExtensions, langExtensions []string) []string {
	if len(fileExtensions) == 0 {
		return langExtensions
	}

	if len(langExtensions) == 0 {
		return fileExtensions
	}

	m := make(map[string]struct{})

	for _, ext := range langExtensions {
		m[ext] = struct{}{}
	}

	var extensions []string

	for _, ext := range fileExtensions {
		if _, ok := m[ext]; ok {
			extensions = append(extensions, ext)
		}
	}

	return extensions
}

func selectByExtensions(files, extensions []string) []string {
	if len(extensions) == 0 {
		return files
	}

	var includedFiles []string

	for _, file := range files {
		for _, extension := range extensions {
			if filepath.Ext(file) == extension {
				includedFiles = append(includedFiles, file)
				break
			}
		}
	}

	return includedFiles
}
