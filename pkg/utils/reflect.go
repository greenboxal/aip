package utils

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/ipld/go-ipld-prime/schema"
	"github.com/samber/lo"
)

var NormalizeTypeNameRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

type ParsedTypeName struct {
	Pkg  string
	Name string
	Args []ParsedTypeName
}

func (n ParsedTypeName) FullName() string {
	if n.Pkg != "" {
		return n.Pkg + "." + n.Name
	}

	return n.Name
}

func (n ParsedTypeName) GoString() string {
	return n.String()
}

func (n ParsedTypeName) String() string {
	args := ""

	if len(n.Args) > 0 {
		a := lo.Map(n.Args, func(arg ParsedTypeName, _index int) string {
			return arg.String()
		})

		args = strings.Join(a, ", ")
		args = "[" + args + "]"
	}

	return n.Name + args
}

func (n ParsedTypeName) NormalizedFullNameWithArguments() string {
	args := ""

	if len(n.Args) > 0 {
		a := lo.Map(n.Args, func(arg ParsedTypeName, _index int) string {
			return arg.String()
		})

		args = strings.Join(a, "__")
		args = "___" + args + "___"
	}

	return NormalizeName(n.FullName() + args)
}

func GetParsedTypeName(typ reflect.Type) ParsedTypeName {
	name := typ.Name()

	if name == "" {
		name = typ.Kind().String()
	} else {
		pkg := typ.PkgPath()

		name = pkg + "." + name
	}

	return ParseTypeName(name)
}

func ParseTypeName(fullName string) ParsedTypeName {
	var args []ParsedTypeName
	var name string

	nested := 0
	lastOffset := 0

	for i, c := range fullName {
		if nested == 0 {
			if c == '[' {
				name = fullName[lastOffset:i]
				lastOffset = i + 1
				nested++
			}
		} else {
			if c == '[' {
				nested++
			} else if c == ']' {
				nested--
			}

			if nested == 0 || (c == ',' && nested == 1) {
				arg := fullName[lastOffset:i]
				args = append(args, ParseTypeName(arg))
				lastOffset = i + 1
			}
		}
	}

	if len(name) == 0 {
		name = fullName
	}

	name = strings.TrimSpace(name)

	parts := strings.Split(name, ".")
	pkg := strings.Join(parts[:len(parts)-1], ".")
	name = parts[len(parts)-1]

	name = NormalizeName(name)

	return ParsedTypeName{
		Pkg:  pkg,
		Name: name,
		Args: args,
	}
}

func NormalizeName(str string) string {
	return NormalizeTypeNameRegex.ReplaceAllString(str, "")
}

func NormalizedTypeName(typ reflect.Type) schema.TypeName {
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	name := typ.Name()

	if name == "" {
		name = typ.Kind().String()
	}

	return NormalizeTypeNameRegex.ReplaceAllString(name, "")
}
