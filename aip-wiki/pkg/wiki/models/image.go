package models

import (
	"encoding/json"

	"github.com/multiformats/go-multibase"
	"github.com/multiformats/go-multihash"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
)

type ImageID struct {
	forddb.StringResourceID[*Image]
}

type Image struct {
	forddb.ContentAddressedResourceBase[ImageID, *Image] `json:"metadata"`

	Spec   ImageSpec   `json:"spec"`
	Status ImageStatus `json:"status"`
}

func (i *Image) GetContentAddressableRoot() any {
	return i.Spec
}

type ImageStatus struct {
	URL    string `json:"url"`
	Prompt string `json:"prompt"`
}

type ImageSpec struct {
	Path   string `json:"path"`
	Prompt string `json:"prompt"`
}

func BuildImageID(spec ImageSpec) ImageID {
	data, err := json.Marshal(spec)

	if err != nil {
		panic(err)
	}

	h, err := multihash.Sum(data, multihash.SHA1, -1)

	if err != nil {
		panic(err)
	}

	b, err := multibase.Encode(multibase.Base36, h)

	if err != nil {
		panic(err)
	}

	return forddb.NewStringID[ImageID](b)
}
