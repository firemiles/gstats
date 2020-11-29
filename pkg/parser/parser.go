package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/firemiles/gstats/pkg/model"
)

var debugAstOfSources = false

type myParser struct {
}

// New : ...
func New() Parser {
	return &myParser{}
}

func (p *myParser) ParseSourceDir(dirName string, includeRegex string, excludeRegex string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFilesInDir(dirName)
	}
	packages, err := parseDir(dirName, includeRegex, excludeRegex)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return model.ParsedSources{}, err
	}

	v := &astVisitor{
		Imports: map[string]string{},
	}
	for _, aPackage := range packages {
		parsePackage(aPackage, v)
	}

	return model.ParsedSources{
		Operations: v.Operations,
	}, nil
}

func parsePackage(aPackage *ast.Package, v *astVisitor) {
	for _, fileEntry := range sortedFileEntries(aPackage.Files) {
		v.CurrentFilename = fileEntry.key

		appEngineOnly := true
		for _, commentGroup := range fileEntry.file.Comments {
			if commentGroup != nil {
				for _, comment := range commentGroup.List {
					if comment != nil && comment.Text == "// +build !appengine" {
						appEngineOnly = false
					}
				}
			}
		}
		if appEngineOnly {
			ast.Walk(v, &fileEntry.file)
		}
	}
}

func parseSourceFile(srcFilename string) (model.ParsedSources, error) {
	if debugAstOfSources {
		dumpFile(srcFilename)
	}

	v, err := doParseFile(srcFilename)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return model.ParsedSources{}, err
	}

	return model.ParsedSources{
		Operations: v.Operations,
	}, nil
}

func doParseFile(srcFilename string) (*astVisitor, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src-file %s: %s", srcFilename, err.Error())
		return nil, err
	}
	v := &astVisitor{
		Imports: map[string]string{},
	}
	v.CurrentFilename = srcFilename
	ast.Walk(v, file)

	return v, nil
}

type fileEntry struct {
	key  string
	file ast.File
}

type fileEntries []fileEntry

func (list fileEntries) Len() int {
	return len(list)
}

func (list fileEntries) Less(i, j int) bool {
	return list[i].key < list[j].key
}

func (list fileEntries) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func sortedFileEntries(fileMap map[string]*ast.File) fileEntries {
	var fileEntries fileEntries = make([]fileEntry, 0, len(fileMap))
	for key, file := range fileMap {
		if file != nil {
			fileEntries = append(fileEntries, fileEntry{
				key:  key,
				file: *file,
			})
		}
	}
	sort.Sort(fileEntries)
	return fileEntries
}

func parseDir(dirName string, includeRegex string, excludeRegex string) (map[string]*ast.Package, error) {
	var includePattern = regexp.MustCompile(includeRegex)
	var excludePattern = regexp.MustCompile(excludeRegex)

	fileSet := token.NewFileSet()
	packageMap, err := parser.ParseDir(fileSet, dirName, func(fi os.FileInfo) bool {
		if excludePattern.MatchString(fi.Name()) {
			return false
		}
		return includePattern.MatchString(fi.Name())
	}, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
		return packageMap, err
	}

	return packageMap, nil
}

func dumpFile(srcFilename string) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, srcFilename, nil, parser.ParseComments)
	if err != nil {
		log.Printf("error parsing src %s: %s", srcFilename, err.Error())
		return
	}
	ast.Print(fileSet, file)
}

func dumpFilesInDir(dirName string) {
	fileSet := token.NewFileSet()
	packageMap, err := parser.ParseDir(
		fileSet,
		dirName,
		nil,
		parser.ParseComments)
	if err != nil {
		log.Printf("error parsing dir %s: %s", dirName, err.Error())
	}
	for _, aPackage := range packageMap {
		for _, file := range aPackage.Files {
			ast.Print(fileSet, file)
		}
	}
}

// =====================================================================================================================

type astVisitor struct {
	CurrentFilename string
	PackageName     string
	Filename        string
	Imports         map[string]string
	Operations      []model.Operation
}

func (v *astVisitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {

		// package-name is in isolated node
		if packageName, ok := extractPackageName(node); ok {
			v.PackageName = packageName
		}

		// extract all imports into a map
		v.extractGenDeclImports(node)

		v.parseAsOperation(node)

	}
	return v
}

func (v *astVisitor) extractGenDeclImports(node ast.Node) {
	if genDecl, ok := node.(*ast.GenDecl); ok {
		for _, spec := range genDecl.Specs {
			if importSpec, ok := spec.(*ast.ImportSpec); ok {
				quotedImport := importSpec.Path.Value
				unquotedImport := strings.Trim(quotedImport, "\"")
				init, last := filepath.Split(unquotedImport)
				if init == "" {
					last = init
				}
				v.Imports[last] = unquotedImport
			}
		}
	}
}

func (v *astVisitor) parseAsOperation(node ast.Node) {
	// if mOperation, get its signature
	if mOperation := extractOperation(node, v.Imports); mOperation != nil {
		mOperation.PackageName = v.PackageName
		mOperation.Filename = v.CurrentFilename
		v.Operations = append(v.Operations, *mOperation)
	}
}

// =====================================================================================================================

func extractPackageName(node ast.Node) (string, bool) {
	if file, ok := node.(*ast.File); ok {
		if file.Name != nil {
			return file.Name.Name, true
		}
		return "", true
	}
	return "", false
}

// ----------------------------------------------------- OPERATION -----------------------------------------------------

func extractOperation(node ast.Node, imports map[string]string) *model.Operation {
	if funcDecl, ok := node.(*ast.FuncDecl); ok {
		mOperation := model.Operation{
			DocLines: extractComments(funcDecl.Doc),
		}

		if funcDecl.Recv != nil {
			fields := extractFieldList(funcDecl.Recv, imports)
			if len(fields) >= 1 {
				mOperation.RelatedStruct = &(fields[0])
			}
		}

		if funcDecl.Name != nil {
			mOperation.Name = funcDecl.Name.Name
		}

		if funcDecl.Type.Params != nil {
			mOperation.InputArgs = extractFieldList(funcDecl.Type.Params, imports)
		}

		if funcDecl.Type.Results != nil {
			mOperation.OutputArgs = extractFieldList(funcDecl.Type.Results, imports)
		}
		return &mOperation
	}
	return nil
}
