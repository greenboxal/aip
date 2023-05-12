package forddb

import (
	"context"
)

func GetOrCreateResourceNode[T ResourceNode](
	resource BasicResource,
	constructor func(resource BasicResource) T,
) (def T) {
	if resource == nil {
		return
	}

	if node := resource.GetResourceNode(); node != nil {
		return node.(T)
	}

	node := constructor(resource)

	resource.setResourceNode(node)

	return node
}

type ResourceNode interface {
	HasListeners

	ResourceID() BasicResourceID
	ResourceType() TypeID

	Get(ctx context.Context) (BasicResource, error)
}
