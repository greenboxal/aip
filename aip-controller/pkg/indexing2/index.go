package indexing2

import (
	"context"
)

type IndexConfiguration struct {
	Reducer Reducer
}

type Index interface {
	Configuration() IndexConfiguration

	OpenSession(ctx context.Context, options SessionOptions) (Session, error)
}
