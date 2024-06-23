package git

import (
	"fmt"
	"path/filepath"
)

func SelectByGlob(files, excludeGlobs, strictGlobs []string) ([]string, error) {
	files, err := excludeGlob(files, excludeGlobs)
	if err != nil {
		return nil, err
	}

	return strictToGlob(files, strictGlobs)
}

func excludeGlob(files, globs []string) ([]string, error) {
	if len(globs) == 0 {
		return files, nil
	}

	var includedFiles []string

	for _, file := range files {
		var match bool
		for _, glob := range globs {
			tryMatch, err := filepath.Match(glob, file)
			if err != nil {
				return nil, fmt.Errorf("filepath match %s: %v", glob, err)
			}

			if tryMatch {
				match = true
			}
		}

		if !match {
			includedFiles = append(includedFiles, file)
		}
	}

	return includedFiles, nil
}

func strictToGlob(files, globs []string) ([]string, error) {
	if len(globs) == 0 {
		return files, nil
	}

	var includedFiles []string

	for _, file := range files {
		for _, glob := range globs {
			match, err := filepath.Match(glob, file)
			if err != nil {
				return nil, fmt.Errorf("filepath match %s: %v", glob, err)
			}

			if match {
				includedFiles = append(includedFiles, file)
			}
		}
	}

	return includedFiles, nil
}
