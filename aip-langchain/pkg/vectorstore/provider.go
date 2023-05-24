package vectorstore

type Provider interface {
	Collection(name string, dim int) Collection
}
