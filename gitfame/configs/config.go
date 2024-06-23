package configs

import (
	_ "embed"
)

var (
	//go:embed language_extensions.json
	LanguageExtensions []byte
)
