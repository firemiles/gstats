package externalapi

import (
	"github.com/firemiles/gstats/pkg/annotation"
)

const (
	// TypeExternalAPI : ...
	TypeExternalAPI = "ExternalAPI"
	ParamServer     = "server"
	ParamURL        = "url"
	ParamMethod     = "method"
	ParamStatusCode = "code"
	ParamReturnBody = "body"
)

// Get : ...
func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeExternalAPI,
			ParamNames: []string{ParamURL, ParamMethod, ParamStatusCode, ParamReturnBody},
			Validator:  validateExternalAPIAnnotation,
		},
	}
}

func validateExternalAPIAnnotation(anno annotation.Annotation) bool {
	return true
}
