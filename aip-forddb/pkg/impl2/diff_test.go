package forddbimpl

import (
	"context"
	"testing"
	"time"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"

	"github.com/greenboxal/aip/aip-forddb/pkg/typesystem"
)

type TestEmbeddedStruct struct {
	Meta1 string `json:"meta_1"`
}

type TestInlined string

func (t *TestInlined) UnmarshalText(text []byte) error {
	*t = TestInlined(text)
	return nil
}

func (t TestInlined) MarshalText() (text []byte, err error) {
	return []byte(t), nil
}

type TestInline struct {
	TestInlined `ipld:",inline"`
}

type TestNamedEmbeddedStruct struct {
	ID TestInline `json:"id"`

	TestEmbeddedStruct
}

type TestScalars struct {
	A bool
	B int
	C string

	Time     time.Time
	Duration time.Duration
}

type TestComplexScalars struct {
}

type TestPointers struct {
	A *bool
	B *int
	C *string

	Time     *time.Time
	Duration *time.Duration
}

type TestNils struct {
	List []string
}

type TestLists struct {
	TestListScalar []string
	TestListStruct []TestScalars
}

type TestMaps struct {
	TestMapStringStruct map[string]TestScalars
	TestMapIntStruct    map[int]TestScalars
	TestMapUint64Struct map[uint64]TestScalars

	TestMapStringString map[string]string
	TestMapIntInt       map[int]int
	TestMapUint64Uint64 map[uint64]uint64
}

type TestStruct1 struct {
	TestNamedEmbeddedStruct `json:"metadata"`

	Nils    TestNils
	Scalars TestScalars
	Lists   TestLists
	Maps    TestMaps
	Complex TestComplexScalars

	TestPointersNull TestPointers
	TestPointers     TestPointers

	TestIfaceNode interface{}
}

func TestDiff(t *testing.T) {
	f := faker.New()

	a := TestStruct1{}
	b := TestStruct1{}

	f.Struct().Fill(&a)
	f.Struct().Fill(&b)

	an := typesystem.Wrap(a)
	bn := typesystem.Wrap(b)

	diff, err := DiffObject(context.Background(), an, bn)

	require.NoError(t, err)
	require.NotEmpty(t, diff.Changes)

	serialized, err := ipld.Encode(typesystem.Wrap(diff), dagjson.Encode)

	require.NoError(t, err)
	require.NotEmpty(t, serialized)
}
