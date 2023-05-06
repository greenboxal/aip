package documents

import (
	"context"
)

type Store interface {
	Put(ctx context.Context, document Document) error
}
