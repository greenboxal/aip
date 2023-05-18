package sema

import (
	"context"
	"strings"

	"github.com/samber/lo"

	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	"github.com/greenboxal/aip/aip-langchain/pkg/chunkers"
)

type SemanticContent struct {
	Content    string `json:"content"`
	TokenCount int    `json:"token_count"`
}

func (sc SemanticContent) String() string {
	return sc.Content
}

type SemanticUnit struct {
	Contents []SemanticContent
}

func (su SemanticUnit) TokenCount() int {
	return lo.SumBy(su.Contents, func(item SemanticContent) int {
		return item.TokenCount
	})
}

func (su SemanticUnit) String() string {
	strs := lo.Map(su.Contents, func(item SemanticContent, _ int) string {
		return item.String()
	})

	return strings.Join(strs, "\n")
}

func (su SemanticUnit) Add(other SemanticUnit) SemanticUnit {
	contents := make([]SemanticContent, 0, len(su.Contents)+len(other.Contents))
	contents = append(contents, su.Contents...)
	contents = append(contents, other.Contents...)

	return SemanticUnit{
		Contents: contents,
	}
}

func (su SemanticUnit) Reshape(ctx context.Context, chunker chunkers.Chunker, maxTokens int) (result SemanticUnit, err error) {
	var buffer SemanticContent

	flushBuffer := func() {
		if buffer.Content == "" {
			return
		}

		result.Contents = append(result.Contents, buffer)

		buffer = SemanticContent{}
	}

	for _, content := range su.Contents {
		if content.TokenCount > maxTokens {
			flushBuffer()

			chunks, err := chunker.SplitTextIntoChunks(ctx, content.Content, maxTokens, 0)

			if err != nil {
				return result, err
			}

			for _, chunk := range chunks {
				result.Contents = append(result.Contents, SemanticContent{
					Content:    chunk.Content,
					TokenCount: chunk.TokenCount,
				})
			}
		} else {
			if buffer.TokenCount+content.TokenCount > maxTokens {
				flushBuffer()
			}

			buffer.Content += buffer.Content
			buffer.TokenCount += buffer.TokenCount
		}
	}

	flushBuffer()

	return
}

type SemanticNodeID struct {
	forddb.StringResourceID[*SemanticNode] `ipld:",inline"`
}

type SemanticNode struct {
	forddb.ResourceBase[SemanticNodeID, *SemanticNode] `json:"metadata"`

	Spec   SemanticNodeSpec   `json:"spec"`
	Status SemanticNodeStatus `json:"status"`
}

type SemanticNodeSpec struct {
	SemanticRootID     SemanticNodeID `json:"semantic_root_id"`
	SemanticRootRow    int            `json:"semantic_root_row"`
	SemanticRootHeight int            `json:"semantic_root_height"`

	ParentNodeIds []SemanticNodeID `json:"parent_node_ids"`
}

type SemanticNodeStatus struct {
	Value SemanticUnit `json:"value"`
}

func init() {
	forddb.DefineResourceType[SemanticNodeID, *SemanticNode]("semantic_node")
}
