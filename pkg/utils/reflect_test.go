package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestSimple struct {
}

type TestGenericA[T any] struct {
}

type TestGenericB[T any] struct {
}

func TestParseTypeName(t *testing.T) {
	simpleName := GetParsedTypeName(reflect.TypeOf(TestSimple{}))

	genericNameAny1 := GetParsedTypeName(reflect.TypeOf(TestGenericA[int]{}))
	genericNameAny2 := GetParsedTypeName(reflect.TypeOf(TestGenericB[bool]{}))
	genericNameNested1 := GetParsedTypeName(reflect.TypeOf(TestGenericB[TestGenericA[string]]{}))

	require.Equal(t, simpleName.Name, "TestSimple")
	require.Len(t, simpleName.Args, 0)

	require.Equal(t, genericNameAny1.Name, "TestGenericA")
	require.Len(t, genericNameAny1.Args, 1)
	require.Equal(t, genericNameAny1.Args[0].Name, "int")

	require.Equal(t, genericNameAny2.Name, "TestGenericB")
	require.Len(t, genericNameAny2.Args, 1)
	require.Equal(t, genericNameAny2.Args[0].Name, "bool")

	require.Equal(t, genericNameNested1.Name, "TestGenericB")
	require.Len(t, genericNameNested1.Args, 1)
	require.Equal(t, genericNameNested1.Args[0].Name, "TestGenericA")
	require.Len(t, genericNameNested1.Args[0].Args, 1)
	require.Equal(t, genericNameNested1.Args[0].Args[0].Name, "string")
}
