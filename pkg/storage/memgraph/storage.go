package memgraph

import (
	"context"
	"net/url"
	"strconv"

	"github.com/prahaladd/gograph/core"
	_ "github.com/prahaladd/gograph/memgraph"
	"github.com/prahaladd/gograph/neo"
	"github.com/prahaladd/gograph/omg"

	"github.com/greenboxal/aip/pkg/collective"
	"github.com/greenboxal/aip/pkg/ford/forddb"
)

type StorageConfig struct {
	ConnectionUri string `json:"connection_uri"`
	Username      string `json:"username"`
	Password      string `json:"password"`
}

type Storage struct {
	forddb.HasListenersBase

	conn  core.Connection
	store omg.Store
}

func (s *Storage) Close() error {
	return s.conn.Close(context.TODO())
}

func NewStorage(config *StorageConfig) (*Storage, error) {
	config.ConnectionUri = "bolt+ssc://54.152.19.198:7687"
	config.Username = "aip"
	config.Password = "aip"

	connectionUrl, err := url.Parse(config.ConnectionUri)

	if err != nil {
		return nil, err
	}
	factory := core.GetConnectorFactory("memgraph")
	port := connectionUrl.Port()

	var p *int32

	if port != "" {
		i, err := strconv.Atoi(port)

		if err != nil {
			return nil, err
		}

		p2 := int32(i)
		p = &p2
	}

	conn, err := factory(
		connectionUrl.Scheme,
		connectionUrl.Hostname(),
		"",
		p,
		map[string]interface{}{
			neo.NEO4J_USER_KEY: config.Username,
			neo.NEO4J_PWD_KEY:  config.Password,
		},
		nil,
	)

	if err != nil {
		return nil, err
	}

	return &Storage{
		conn:  conn,
		store: omg.NewGenericStore(conn, newTypeSystemMapper()),
	}, nil
}

func (s *Storage) AppendSegment(ctx context.Context, segment *collective.MemorySegment) error {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetSegment(ctx context.Context, id collective.MemorySegmentID) (*collective.MemorySegment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetMemory(ctx context.Context, id collective.MemoryID) (*collective.Memory, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) List(ctx context.Context, typ forddb.ResourceTypeID) ([]forddb.BasicResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Get(ctx context.Context, typ forddb.ResourceTypeID, id forddb.BasicResourceID) (forddb.BasicResource, error) {
	var resource vertexWrapper

	resource.ResourceID = id.String()
	resource.ResourceKind = typ.Name()

	if _, err := s.store.ReadVertex(ctx, &resource); err != nil {
		return nil, err
	}

	return forddb.Deserialize([]byte(resource.ResourceManifest), forddb.Json)
}

func (s *Storage) Put(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	serialized, err := forddb.Serialize(forddb.Json, resource)

	if err != nil {
		return nil, err
	}

	err = s.store.PersistVertex(ctx, &vertexWrapper{
		ResourceKind:     resource.GetType().Name(),
		ResourceID:       resource.GetResourceID().String(),
		ResourceManifest: string(serialized),
	})

	if err != nil {
		return nil, err
	}

	// FIXME:
	forddb.FireListeners(&s.HasListenersBase, resource.GetResourceID(), resource, resource)

	return resource, nil
}

func (s *Storage) Delete(ctx context.Context, resource forddb.BasicResource) (forddb.BasicResource, error) {
	_, err := s.conn.ExecuteQuery(ctx, "MATCH (n) WHERE n.metadata.id = $id AND kind = $kind DELETE n", core.Write, map[string]interface{}{
		"id":   resource.GetResourceID().String(),
		"kind": resource.GetType().Name(),
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

type vertexWrapper struct {
	ResourceKind     string `json:"resource_kind"`
	ResourceID       string `json:"resource_id"`
	ResourceManifest string `json:"resource_manifest"`
}

func (v *vertexWrapper) GetLabel() string {
	return v.ResourceKind
}

func (v *vertexWrapper) GetType() omg.GraphObjectType {
	return omg.Vertex
}

type typeSystemMapper struct {
	*omg.ReflectionMapper
}

func newTypeSystemMapper() *typeSystemMapper {
	return &typeSystemMapper{
		ReflectionMapper: omg.NewReflectionMapper(),
	}
}
