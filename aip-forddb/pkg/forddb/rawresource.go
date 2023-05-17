package forddb

import (
	"github.com/ipld/go-ipld-prime/schema"
)

type RawResource struct{ schema.TypedNode }

func (r RawResource) IsEmpty() bool {
	return r.TypedNode == nil
}

type RawResourceHeader struct {
	ID string `json:"id"`

	Metadata
}

func (r RawResource) GetResourceMetadata() *Metadata {
	metadata, err := r.LookupByString("metadata")

	if err != nil {
		panic(err)
	}

	if metadata == nil {
		return nil
	}

	m, err := ConvertNode[RawResourceHeader](metadata)

	if err != nil {
		panic(err)
	}

	return &m.Metadata
}

func (r RawResource) GetResourceBasicID() BasicResourceID {
	metadata, err := r.LookupByString("metadata")

	if err != nil {
		panic(err)
	}

	if metadata == nil {
		return nil
	}

	id, err := metadata.LookupByString("id")

	if err != nil {
		panic(err)
	}

	idStr, err := id.AsString()

	if err != nil {
		panic(err)
	}

	return r.GetResourceTypeID().Type().CreateID(idStr)
}

func (r RawResource) GetResourceTypeID() TypeID {
	metadata, err := r.LookupByString("metadata")

	if err != nil {
		panic(err)
	}

	kind, err := metadata.LookupByString("kind")

	if err != nil {
		panic(err)
	}

	kindStr, err := kind.AsString()

	if err != nil {
		panic(err)
	}

	return NewStringID[TypeID](kindStr)
}

func (r RawResource) GetResourceVersion() uint64 {
	if r.IsEmpty() {
		return 0
	}

	metadata, err := r.LookupByString("metadata")

	if err != nil {
		panic(err)
	}

	version, err := metadata.LookupByString("version")

	if err != nil {
		panic(err)
	}

	versionVal, err := version.AsInt()

	if err != nil {
		panic(err)
	}

	return uint64(versionVal)
}
