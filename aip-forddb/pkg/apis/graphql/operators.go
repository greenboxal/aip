package graphql

var operatorMap = map[string]string{
	"==": "eq",
	"!=": "neq",
	"<":  "lt",
	"<=": "lte",
	">":  "gt",
	">=": "gte",
}

var reverseOperatorMap = map[string]string{}

func init() {
	for k, v := range operatorMap {
		reverseOperatorMap[v] = k
	}
}
