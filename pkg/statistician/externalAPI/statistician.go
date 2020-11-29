package externalapi

import (
	"encoding/json"
	"fmt"

	"github.com/firemiles/gstats/pkg/annotation"
	"github.com/firemiles/gstats/pkg/model"
	"github.com/firemiles/gstats/pkg/statistician"
)

// Statistician : ...
type Statistician struct {
}

type ExternalAPIInfo struct {
	Server string
	Method string
	URL    string
	Code   string
	Body   string
}

// NewStatistician : ...
func NewStatistician() statistician.Statistician {
	return &Statistician{}
}

// GetAnnotations : ...
func (s *Statistician) GetAnnotations() []annotation.AnnotationDescriptor {
	return Get()
}

// Statistics : ...
func (s *Statistician) Statistics(inputDir string, parsedSource model.ParsedSources) ([]statistician.Info, error) {
	ops := parsedSource.Operations
	var infos []statistician.Info
	for _, op := range ops {
		ann, ok := GetAnnotation(op)
		if !ok {
			continue
		}
		infos = append(infos, statistician.Info{
			Src: fmt.Sprintf("%s.%s", op.PackageName, op.Name),
			Data: ExternalAPIInfo{
				Server: GetServer(&ann),
				Method: GetMethod(&ann),
				URL:    GetURL(&ann),
				Code:   GetStatusCode(&ann),
				Body:   GetBody(&ann),
			},
		})
	}
	return infos, nil
}

// PrettyFormat : ...
func (s *Statistician) PrettyFormat(infos []statistician.Info) string {
	b, _ := json.Marshal(infos)
	return string(b)
}

// GetAnnotation : get external api annotation
func GetAnnotation(s model.Operation) (annotation.Annotation, bool) {
	return annotation.NewRegistry(Get()).ResolveAnnotationByName(s.DocLines, TypeExternalAPI)
}

// GetServer : ...
func GetServer(ann *annotation.Annotation) string {
	if ann != nil {
		return ann.Attributes[ParamServer]
	}
	return ""
}

// GetURL : ...
func GetURL(ann *annotation.Annotation) string {
	if ann != nil {
		return ann.Attributes[ParamURL]
	}
	return ""
}

// GetMethod : ...
func GetMethod(ann *annotation.Annotation) string {
	if ann != nil {
		return ann.Attributes[ParamMethod]
	}
	return ""
}

// GetStatusCode : ...
func GetStatusCode(ann *annotation.Annotation) string {
	if ann != nil {
		return ann.Attributes[ParamStatusCode]
	}
	return ""
}

// GetBody : ...
func GetBody(ann *annotation.Annotation) string {
	if ann != nil {
		return ann.Attributes[ParamReturnBody]
	}
	return ""
}
