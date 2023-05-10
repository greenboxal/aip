package forddb

import (
	"time"

	"github.com/mashingan/smapping"
)

type RawResource smapping.Mapped

func (r RawResource) GetResourceMetadata() *Metadata {
	var result Metadata
	var fakeMetadata struct {
		Kind      string    `json:"kind"`
		Namespace string    `json:"namespace"`
		Version   uint64    `json:"version"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	}

	metadata := r["metadata"]

	if metadata == nil {
		return nil
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return nil
	}

	if err := smapping.FillStructByTags(&fakeMetadata, metadataMapped, "json"); err != nil {
		panic(err)
	}

	result.Kind = NewStringID[TypeID](fakeMetadata.Kind)
	result.UpdatedAt = fakeMetadata.UpdatedAt
	result.CreatedAt = fakeMetadata.CreatedAt
	result.Version = fakeMetadata.Version
	result.Namespace = fakeMetadata.Namespace

	return &result
}

func (r RawResource) GetResourceBasicID() BasicResourceID {
	metadata := r["metadata"]

	if metadata == nil {
		return nil
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return nil
	}

	idVal, ok := metadataMapped["id"]

	if !ok {
		return nil
	}

	return r.GetResourceTypeID().Type().CreateID(idVal.(string))
}

func (r RawResource) GetResourceTypeID() TypeID {
	metadata := r["metadata"]

	if metadata == nil {
		return ""
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return ""
	}

	kindVal, ok := metadataMapped["kind"]

	if !ok {
		return ""
	}

	return NewStringID[TypeID](kindVal.(string))
}

func (r RawResource) GetResourceVersion() uint64 {
	metadata := r["metadata"]

	if metadata == nil {
		return 0
	}

	metadataMapped, ok := metadata.(map[string]interface{})

	if !ok {
		return 0
	}

	versionVal, ok := metadataMapped["version"]

	if !ok {
		return 0
	}

	if v, ok := versionVal.(uint64); ok {
		return v
	}

	if v, ok := versionVal.(float64); ok {
		return uint64(v)
	}

	return 0
}