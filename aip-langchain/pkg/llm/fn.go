package llm

import "github.com/sashabaranov/go-openai"

type FunctionDeclaration struct {
	Name        string
	Description string
	Parameters  *FunctionParams
}

type FunctionParams struct {
	Type       openai.JSONSchemaType               `json:"type"`
	Properties map[string]*openai.JSONSchemaDefine `json:"properties,omitempty"`
	Required   []string                            `json:"required,omitempty"`
}
