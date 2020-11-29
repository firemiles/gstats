package parser

import "github.com/firemiles/gstats/pkg/model"

type Parser interface {
	ParseSourceDir(dirName string, includeRegex string, excludeRegex string) (model.ParsedSources, error)
}
