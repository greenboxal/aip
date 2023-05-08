package firestore

import (
	"context"
	"encoding/json"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
)

type Config struct {
	ProjectID    string
	CollectionID string
}

func (c *Config) SetDefaults() {
	if c.ProjectID == "" {
		c.ProjectID = "uncyclo-385820"
	}

	if c.CollectionID == "" {

	}
}

type Storage struct {
	config *Config
	client *firestore.Client
}

func NewStorage(config *Config) (*Storage, error) {
	ctx := context.Background()

	conf := &firebase.Config{ProjectID: config.ProjectID}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)

	if err != nil {
		return nil, err
	}

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) List(
	ctx context.Context,
	typ forddb.TypeID,
	opts forddb.ListOptions,
) ([]forddb.RawResource, error) {
	collection := typ.Name()

	col := s.client.Collection(collection)

	query := col.Query

	if opts.Offset > 0 {
		query = col.Offset(opts.Offset)
	}

	if opts.Limit > 0 {
		query = col.Limit(opts.Limit)
	}

	for _, item := range opts.SortFields {
		itemPath := item.Path

		if itemPath == "id" {
			itemPath = "metadata.id"
		}

		if item.Order == forddb.Asc {
			query = query.OrderBy(itemPath, firestore.Asc)
		} else {
			query = query.OrderBy(itemPath, firestore.Desc)
		}
	}

	if len(opts.ResourceIDs) > 0 {
		ids := lo.Map(opts.ResourceIDs, func(id forddb.BasicResourceID, _index int) string {
			return id.String()
		})

		query = query.Where("metadata.id", "in", ids)
	}

	iterator := query.Documents(ctx)

	all, err := iterator.GetAll()

	if err != nil {
		return nil, err
	}

	raws := make([]forddb.RawResource, len(all))

	for i, v := range all {
		var raw forddb.RawResource

		data := v.Data()

		serialized, err := json.Marshal(data)

		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(serialized, &raw); err != nil {
			return nil, err
		}

		raws[i] = raw
	}

	return raws, nil
}

func (s *Storage) Get(
	ctx context.Context,
	typ forddb.TypeID,
	id forddb.BasicResourceID,
	opts forddb.GetOptions,
) (forddb.RawResource, error) {
	collection := typ.Name()

	col := s.client.Collection(collection)
	doc := col.Doc(id.String())

	result, err := doc.Get(ctx)

	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, forddb.ErrNotFound
		}

		return nil, err
	}

	serialized, err := json.Marshal(result.Data())

	if err != nil {
		return nil, err
	}

	var raw forddb.RawResource

	if err := json.Unmarshal(serialized, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}

func (s *Storage) Put(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.PutOptions,
) (forddb.RawResource, error) {
	var fields map[string]interface{}
	var preconditions []firestore.Precondition

	serialized, err := json.Marshal(resource)

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(serialized, &fields); err != nil {
		return nil, err
	}

	collection := resource.GetResourceTypeID().Name()
	col := s.client.Collection(collection)
	doc := col.Doc(resource.GetResourceBasicID().String())

	switch opts.OnConflict {
	case forddb.OnConflictReplace:
		if _, err := doc.Set(ctx, fields); err != nil {
			return nil, err
		}

	default:
		metadata := resource.GetResourceMetadata()
		updates := make([]firestore.Update, 0, len(fields))

		for k, v := range fields {
			updates = append(updates, firestore.Update{
				Path:  k,
				Value: v,
			})
		}

		updatedAt := metadata.UpdatedAt
		epoch := time.Unix(0, 0)

		if updatedAt.Before(epoch) {
			metadata.UpdatedAt = epoch
		}

		preconditions = append(preconditions, firestore.LastUpdateTime(metadata.UpdatedAt))

		if _, err := doc.Update(ctx, updates, preconditions...); err != nil {
			return nil, err
		}
	}

	return resource, nil
}

func (s *Storage) Delete(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.DeleteOptions,
) (forddb.RawResource, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) Close() error {
	return s.client.Close()
}
