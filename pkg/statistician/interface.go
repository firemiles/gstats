package statistician

import (
	"github.com/firemiles/gstats/pkg/annotation"
	"github.com/firemiles/gstats/pkg/model"
)

// Statistician : ...
type Statistician interface {
	GetAnnotations() []annotation.AnnotationDescriptor
	Statistics(inputDir string, parsedSources model.ParsedSources) ([]Info, error)
	PrettyFormat([]Info) string
}
