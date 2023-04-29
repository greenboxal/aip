package milvus

import (
	"context"
	"os"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/pkg/indexing"
	"github.com/greenboxal/aip/pkg/indexing/reducers"
)

type Storage struct {
	db  client.Client
	oai *openai.Client

	collection string
}

func NewStorage(oai *openai.Client) (*Storage, error) {
	ctx := context.Background()

	c, err := client.NewDefaultGrpcClientWithTLSAuth(
		ctx,
		os.Getenv("MILVUS_ENDPOINT"),
		os.Getenv("MILVUS_USERNAME"),
		os.Getenv("MILVUS_PASSWORD"),
	)

	if err != nil {
		return nil, err
	}

	s := &Storage{
		db:         c,
		oai:        oai,
		collection: "memories",
	}

	if err = s.initialize(ctx); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Storage) AppendSegment(ctx context.Context, segment *indexing.MemorySegment) error {
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
		chunks, err := reducers.SplitTextIntoChunks(string(memory.Data.Data), 512, 64)

		if err != nil {
			return err
		}

		firstIndex := chunkIndex
		lastIndex := firstIndex + len(chunks)
		chunkIndex = lastIndex

		allChunks = reducers.Grow(allChunks, len(chunks))
		segmentIds = reducers.Grow(segmentIds, len(chunks))
		memoryIds = reducers.Grow(memoryIds, len(chunks))
		parentIds = reducers.Grow(parentIds, len(chunks))
		branchIds = reducers.Grow(branchIds, len(chunks))
		rootIds = reducers.Grow(rootIds, len(chunks))
		timestamps = reducers.Grow(timestamps, len(chunks))
		clocks = reducers.Grow(clocks, len(chunks))
		heights = reducers.Grow(heights, len(chunks))

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

func (s *Storage) initialize(ctx context.Context) error {
	ok, err := s.db.HasCollection(ctx, s.collection)

	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	numShards := int32(2)

	schema := &entity.Schema{
		CollectionName: s.collection,
		AutoID:         true,
		Fields: []*entity.Field{
			{Name: "_id", PrimaryKey: true, AutoID: true, DataType: entity.FieldTypeInt64},
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

	err = s.db.CreateCollection(ctx, schema, numShards)

	if err != nil {
		return err
	}

	return nil
}
