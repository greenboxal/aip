package forddb

import (
	"github.com/stoewer/go-strcase"
)

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
		Name:   name,
		Plural: name + "s",
	}
}
