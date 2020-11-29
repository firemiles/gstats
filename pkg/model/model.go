package model

type ParsedSources struct {
	Operations []Operation `json:"operations,omitempty"`
}

type Operation struct {
	PackageName   string   `json:"packageName,omitempty"`
	Filename      string   `json:"filename,omitempty"`
	DocLines      []string `json:"docLines,omitempty"`
	RelatedStruct *Field   `json:"relatedStruct,omitempty"` // optional
	Name          string   `json:"name"`
	InputArgs     []Field  `json:"inputArgs,omitempty"`
	OutputArgs    []Field  `json:"outputArgs,omitempty"`
	CommentLines  []string `json:"commentLines,omitempty"`
}

type Field struct {
	PackageName  string   `json:"packageName,omitempty"`
	DocLines     []string `json:"docLines,omitempty"`
	Name         string   `json:"name,omitempty"`
	TypeName     string   `json:"typeName,omitempty"`
	Tag          string   `json:"tag,omitempty"`
	CommentLines []string `json:"commentLines,omitempty"`
}
