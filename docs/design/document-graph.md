# Document Graph

```go
package documentgraph

type Document struct {
    ID string
	Text string
	
	Links []DocumentLink
}

type DocumentLink struct {
    DocumentID string
}

type DocumentManager interface {
	Resolve(link DocumentLink) (*Document, error)
}

```
