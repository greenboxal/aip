package golang

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type SourceFile struct {
	module *ModuleBuilder

	name     string
	contents string

	parsed *ast.File
}

func NewSourceFile(module *ModuleBuilder, name string, contents string) *SourceFile {
	return &SourceFile{
		module:   module,
		name:     name,
		contents: contents,
	}
}

func (sf *SourceFile) Name() string     { return sf.name }
func (sf *SourceFile) Contents() string { return sf.contents }
func (sf *SourceFile) Root() ast.Node   { return sf.parsed }

func (sf *SourceFile) Build() error {
	node, err := parser.ParseFile(sf.module.fset, sf.name, sf.contents, parser.ParseComments)

	if err != nil {
		return err
	}

	sf.parsed = node

	return nil
}

type ModuleBuilder struct {
	name  string
	fset  *token.FileSet
	files map[string]*SourceFile
}

func NewModuleBuilder(name string) *ModuleBuilder {
	return &ModuleBuilder{
		name:  name,
		fset:  token.NewFileSet(),
		files: map[string]*SourceFile{},
	}
}

func (mb *ModuleBuilder) AddFile(name string, contents string) {
	sf := NewSourceFile(mb, name, contents)

	mb.files[name] = sf
}

func (mb *ModuleBuilder) Build() (Module, error) {
	for _, file := range mb.files {
		if err := file.Build(); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

type Module interface {
	Name() string
}
