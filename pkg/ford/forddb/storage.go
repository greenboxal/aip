package forddb

type Storage interface {
	List(typ ResourceTypeID) ([]BasicResource, error)
	Get(typ ResourceTypeID, id BasicResourceID) (BasicResource, error)
	Put(resource BasicResource) (BasicResource, error)
	Delete(resource BasicResource) (BasicResource, error)
}
