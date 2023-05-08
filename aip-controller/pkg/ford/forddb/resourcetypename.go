package forddb

import (
	"github.com/gertd/go-pluralize"
	"github.com/stoewer/go-strcase"
)

var pluralizeClient = pluralize.NewClient()

type ResourceTypeName struct {
	Name   string
	Plural string
}

func (n ResourceTypeName) ToTitle() string {
	return strcase.UpperCamelCase(n.Name)
}

func (n ResourceTypeName) ToTitlePlural() string {
	return strcase.UpperCamelCase(n.Plural)
}

func ResourceTypeNameFromSingular(name string) ResourceTypeName {
	name = strcase.UpperCamelCase(name)

	return ResourceTypeName{
		Name:   pluralizeClient.Singular(name),
		Plural: pluralizeClient.Plural(name),
	}
}
