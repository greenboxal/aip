package typesystem

import (
	"testing"

	"github.com/ipld/go-ipld-prime"
	"github.com/ipld/go-ipld-prime/codec/dagjson"
	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/require"
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

	Nils  TestNils
	Lists TestLists
	Maps  TestMaps

	TestIfaceNode interface{}
}

func TestTS(t *testing.T) {
	f := faker.New()

	ifaceValue := TestScalars{}
	f.Struct().Fill(&ifaceValue)

	initialValue := TestStruct1{}

	initialValue.Maps.TestMapIntInt = map[int]int{0: 1, 4342: 435345}
	initialValue.Maps.TestMapStringString = map[string]string{"saddas": "", "": "osdkfods", "asofdasa": "adoad"}
	initialValue.Maps.TestMapUint64Uint64 = map[uint64]uint64{0: 1, 4342: 435345}
	initialValue.Maps.TestMapIntStruct = map[int]TestScalars{0: ifaceValue, 4342: ifaceValue}
	initialValue.Maps.TestMapStringStruct = map[string]TestScalars{"saddas": ifaceValue, "": ifaceValue, "asofdasa": ifaceValue}
	initialValue.Maps.TestMapUint64Struct = map[uint64]TestScalars{0: ifaceValue, 4342: ifaceValue}

	f.Struct().Fill(&initialValue)

	initialValue.Nils = TestNils{}
	initialValue.TestIfaceNode = ifaceValue

	typ := TypeOf(TestStruct1{})
	wrapped := Wrap(initialValue)

	require.NotNil(t, typ)
	require.NotNil(t, wrapped)

	iface, err := wrapped.LookupByString("TestIfaceNode")

	require.NoError(t, err)
	require.NotNil(t, iface)

	data, err := ipld.Encode(wrapped, dagjson.Encode)

	require.NoError(t, err)
	require.NotNil(t, data)

	node, err := ipld.Decode(data, dagjson.Decode)

	require.NoError(t, err)
	require.NotNil(t, node)

	nodeWithProto, err := ipld.DecodeUsingPrototype(data, dagjson.Decode, typ.IpldPrototype())

	require.NoError(t, err)
	require.NotNil(t, nodeWithProto)

	unwrapped := Unwrap(nodeWithProto)

	require.NotNil(t, unwrapped)
	require.EqualValues(t, initialValue, unwrapped)
}
