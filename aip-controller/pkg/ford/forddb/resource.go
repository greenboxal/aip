package forddb

type BasicResource interface {
	IResourceMetadata

	GetMetadata() *BasicResourceMetadata

	onBeforeSerialize()
}

type Resource[ID BasicResourceID] interface {
	BasicResource

	GetID() ID
}
