package indexing

import (
	"context"
)

type IndexConfiguration struct {
}

type Index interface {
	Configuration() IndexConfiguration

	OpenSession(ctx context.Context, options SessionOptions) (Session, error)
}
