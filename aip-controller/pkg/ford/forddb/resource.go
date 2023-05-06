package forddb

type BasicResource interface {
	ResourceMetadata

	OnBeforeSave(self BasicResource)
}

type Resource[ID BasicResourceID] interface {
	BasicResource

	GetResourceID() ID
}
