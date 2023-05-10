package iam

import (
	"context"
	"crypto"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type IdentityID struct {
	forddb.StringResourceID[*Identity]
}

type Identity struct {
	forddb.ResourceBase[IdentityID, *Identity] `json:"metadata"`

	PublicKey crypto.PublicKey `json:"public_key"`
}

type KeyStore interface {
	GenerateKey(ctx context.Context, id IdentityID) (crypto.PrivateKey, error)
	GetKey(ctx context.Context, id IdentityID) (crypto.PrivateKey, error)
}

type IdentityStore interface {
	CreateIdentity(ctx context.Context, id IdentityID) (*Identity, error)
	GetIdentity(ctx context.Context, id IdentityID) (*Identity, error)
}
