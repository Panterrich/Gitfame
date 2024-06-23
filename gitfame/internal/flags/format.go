package flags

import (
	"errors"
	"flag"
)

type Format string

var _ flag.Value = (*Format)(nil)

const (
	FormatTabular   Format = "tabular"
	FormatCSV       Format = "csv"
	FormatJSON      Format = "json"
	FormatJSONLines Format = "json-lines"
)

func (f *Format) String() string {
	return string(*f)
}

func (f *Format) Set(v string) error {
	switch v {
	case string(FormatTabular), string(FormatCSV), string(FormatJSON), string(FormatJSONLines):
		*f = Format(v)
		return nil
	default:
		return errors.New(`must be one of "tabular", "csv", "json", or "json-lines"`)
	}
}

func (f *Format) Type() string {
	return "format"
}
