package milvus

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/fx"

	collective2 "github.com/greenboxal/aip/aip-controller/pkg/collective"
	forddb "github.com/greenboxal/aip/aip-controller/pkg/ford/forddb"
	"github.com/greenboxal/aip/aip-controller/pkg/llm/chunkers"
	"github.com/greenboxal/aip/aip-controller/pkg/utils"
)

type Storage struct {
	db  client.Client
	oai *openai.Client

	collection string
}

func NewStorage(lc fx.Lifecycle, oai *openai.Client) (*Storage, error) {
	ctx := context.Background()

	s := &Storage{
		oai:        oai,
		collection: "memories",
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return s.initialize(ctx)
		},
	})

	endpoint := os.Getenv("MILVUS_ENDPOINT")
	username := os.Getenv("MILVUS_USERNAME")
	password := os.Getenv("MILVUS_PASSWORD")

	c, err := client.NewDefaultGrpcClientWithURI(
		ctx,
		endpoint,
		username,
		password,
	)

	if err != nil {
		return nil, err
	}

	s.db = c

	return s, nil
}

func (s *Storage) initialize(ctx context.Context) error {
	ok, _ := s.db.HasCollection(ctx, s.collection)

	if !ok {

		numShards := int32(2)

		schema := &entity.Schema{
			CollectionName: s.collection,
			Fields: []*entity.Field{
				{Name: "_id", PrimaryKey: true, DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "timestamp", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "64"}},
				{Name: "segment_id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "memory_id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "parent_memory_id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "branch_memory_id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "root_memory_id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "clock", DataType: entity.FieldTypeInt64},
				{Name: "height", DataType: entity.FieldTypeInt64},
				{Name: "text", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "embedding", DataType: entity.FieldTypeFloatVector, TypeParams: map[string]string{"dim": "1536"}},
			},
		}

		err := s.db.CreateCollection(ctx, schema, numShards)

		if err != nil {
			return err
		}
	}

	ok, _ = s.db.HasCollection(ctx, "forddb")

	if !ok {
		numShards := int32(2)

		schema := &entity.Schema{
			CollectionName: "forddb",
			Fields: []*entity.Field{
				{Name: "_pk", PrimaryKey: true, DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "256"}},
				{Name: "_vec", DataType: entity.FieldTypeFloatVector, TypeParams: map[string]string{"dim": "1"}},
				{Name: "id", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "kind", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "namespace", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "128"}},
				{Name: "version", DataType: entity.FieldTypeInt64},
				{Name: "data", DataType: entity.FieldTypeVarChar, TypeParams: map[string]string{"max_length": "65535"}},
			},
		}

		err := s.db.CreateCollection(ctx, schema, numShards)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) AppendSegment(ctx context.Context, segment *collective2.MemorySegment) error {
	chunkIndex := 0

	allChunks := make([]string, 0, len(segment.Memories))
	segmentIds := make([]string, 0, len(segment.Memories))
	rootIds := make([]string, 0, len(segment.Memories))
	branchIds := make([]string, 0, len(segment.Memories))
	parentIds := make([]string, 0, len(segment.Memories))
	memoryIds := make([]string, 0, len(segment.Memories))
	timestamps := make([]string, 0, len(segment.Memories))
	heights := make([]int64, 0, len(segment.Memories))
	clocks := make([]int64, 0, len(segment.Memories))

	for _, memory := range segment.Memories {
		chunks, err := chunkers.SplitTextIntoChunks(string(memory.Data.Text), 512, 64)

		if err != nil {
			return err
		}

		firstIndex := chunkIndex
		lastIndex := firstIndex + len(chunks)
		chunkIndex = lastIndex

		allChunks = utils.Grow(allChunks, len(chunks))
		segmentIds = utils.Grow(segmentIds, len(chunks))
		memoryIds = utils.Grow(memoryIds, len(chunks))
		parentIds = utils.Grow(parentIds, len(chunks))
		branchIds = utils.Grow(branchIds, len(chunks))
		rootIds = utils.Grow(rootIds, len(chunks))
		timestamps = utils.Grow(timestamps, len(chunks))
		clocks = utils.Grow(clocks, len(chunks))
		heights = utils.Grow(heights, len(chunks))

		for i := firstIndex; i < lastIndex; i++ {
			allChunks[i] = chunks[i-firstIndex]
			segmentIds[i] = segment.ID.String()
			memoryIds[i] = memory.ID.String()
			parentIds[i] = memory.ParentMemoryID.String()
			branchIds[i] = memory.BranchMemoryID.String()
			rootIds[i] = memory.RootMemoryID.String()
			timestamps[i] = memory.CreatedAt.String()
			clocks[i] = int64(memory.Clock)
			heights[i] = int64(memory.Height)
		}
	}

	result, err := s.oai.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Input: allChunks,
		Model: openai.AdaEmbeddingV2,
	})

	if err != nil {
		return err
	}

	embeddings := make([][]float32, len(allChunks))

	for i, data := range result.Data {
		embeddings[i] = data.Embedding
	}

	columns := []entity.Column{
		entity.NewColumnVarChar("timestamp", timestamps),
		entity.NewColumnVarChar("segment_id", segmentIds),
		entity.NewColumnVarChar("memory_id", memoryIds),
		entity.NewColumnVarChar("parent_memory_id", parentIds),
		entity.NewColumnVarChar("branch_memory_id", branchIds),
		entity.NewColumnVarChar("root_memory_id", rootIds),
		entity.NewColumnInt64("clock", clocks),
		entity.NewColumnInt64("height", heights),
		entity.NewColumnVarChar("text", allChunks),
		entity.NewColumnFloatVector("embedding", 1536, embeddings),
	}

	_, err = s.db.Insert(ctx, s.collection, "_default", columns...)

	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetSegment(ctx context.Context, id collective2.MemorySegmentID) (*collective2.MemorySegment, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) GetMemory(ctx context.Context, id collective2.MemoryID) (*collective2.Memory, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Storage) List(
	ctx context.Context,
	typ forddb.TypeID,
	opts forddb.ListOptions,
) ([]forddb.RawResource, error) {
	res, err := s.db.Query(
		ctx,
		"forddb",
		[]string{"_default"},
		"kind in [\""+typ.Name()+"\"]",
		[]string{"data"},
	)

	if err != nil {
		return nil, err
	}

	result := make([]forddb.RawResource, 0, 32)

	for _, col := range res {
		if col.Name() != "data" {
			continue
		}

		dataColumn, ok := col.(*entity.ColumnVarChar)

		if !ok {
			if len(res) == 0 {
				return nil, forddb.ErrNotFound
			}
		}

		values := dataColumn.Data()

		for _, v := range values {
			var raw forddb.RawResource

			if err := json.Unmarshal([]byte(v), &raw); err != nil {
				return nil, err
			}

			result = append(result, raw)
		}
	}

	return result, nil
}

func (s *Storage) Get(
	ctx context.Context,
	typ forddb.TypeID,
	id forddb.BasicResourceID,
	opts forddb.GetOptions,
) (forddb.RawResource, error) {
	pk := fmt.Sprintf("%s:%s", id.String(), typ.Name())

	primaryKeys := []string{pk}
	pkColumns := entity.NewColumnVarChar("_pk", primaryKeys)

	result, err := s.db.QueryByPks(ctx, "forddb", []string{"_default"}, pkColumns, []string{"data"})

	if err != nil {
		return nil, err
	}

	for _, col := range result {
		if col.Name() != "data" {
			continue
		}

		dataColumn, ok := col.(*entity.ColumnVarChar)

		if !ok {
			if len(result) == 0 {
				return nil, forddb.ErrNotFound
			}
		}

		values := dataColumn.Data()

		if len(values) == 0 {
			return nil, forddb.ErrNotFound
		}

		var raw forddb.RawResource

		if err := json.Unmarshal([]byte(values[0]), &raw); err != nil {
			return nil, err
		}

		return raw, nil
	}

	return nil, forddb.ErrNotFound
}

func (s *Storage) Put(
	ctx context.Context,
	resource forddb.RawResource,
	opts forddb.PutOptions,
) (forddb.RawResource, error) {
	serialized, err := json.Marshal(resource)

	if err != nil {
		return nil, err
	}

	pk := fmt.Sprintf("%s:%s", resource.GetResourceBasicID().String(), resource.GetResourceTypeID().Name())

	primaryKeys := []string{pk}
	primaryVecs := [][]float32{{math.Pi}}
	resourceIds := []string{resource.GetResourceBasicID().String()}
	resourceKinds := []string{resource.GetResourceTypeID().Name()}
	resourceVersions := []int64{int64(resource.GetResourceVersion())}
	resourceNamespaces := []string{resource.GetResourceMetadata().Namespace}
	resourcesData := []string{string(serialized)}

	columns := []entity.Column{
		entity.NewColumnVarChar("_pk", primaryKeys),
		entity.NewColumnFloatVector("_vec", 1, primaryVecs),
		entity.NewColumnVarChar("id", resourceIds),
		entity.NewColumnVarChar("kind", resourceKinds),
		entity.NewColumnInt64("version", resourceVersions),
		entity.NewColumnVarChar("namespace", resourceNamespaces),
		entity.NewColumnVarChar("data", resourcesData),
	}

	_, err = s.db.Insert(ctx, "forddb", "_default", columns...)

	if err != nil {
		return nil, err
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
	return s.db.Close()
}
