package forddb

type BasicResource interface {
	ResourceMetadata

	GetResourceNode() ResourceNode
	setResourceNode(node ResourceNode)

	OnBeforeSave(self BasicResource)
}

type Resource[ID BasicResourceID] interface {
	BasicResource

	GetResourceID() ID
}
