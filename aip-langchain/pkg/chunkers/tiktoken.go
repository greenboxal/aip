package chunkers

import (
	"context"

	"github.com/pkoukk/tiktoken-go"
	"github.com/samber/lo"
)

type TikToken struct{}

func (tt TikToken) SplitTextIntoStrings(ctx context.Context, text string, chunkSize int, overlapSize int) ([]string, error) {
	chunks, err := tt.SplitTextIntoChunks(ctx, text, chunkSize, overlapSize)

	if err != nil {
		return nil, err
	}

	return lo.Map(chunks, func(item Chunk, index int) string {
		return item.Content
	}), nil
}

func (tt TikToken) SplitTextIntoChunks(ctx context.Context, text string, chunkSize int, overlapSize int) ([]Chunk, error) {
	// Tokenize the text using tiktoken
	tokenizer, err := tiktoken.EncodingForModel("text-embedding-ada-002")

	if err != nil {
		return nil, err
	}

	tokens := tokenizer.Encode(text, nil, nil)

	// Calculate the number of chunks based on the chunk size and overlap size
	numChunks := (len(tokens) + chunkSize - 1) / chunkSize

	chunks := make([]Chunk, numChunks)

	// Generate the chunks by iterating over the tokens
	for i := 0; i < numChunks; i++ {
		startIndex := i * (chunkSize - overlapSize)
		endIndex := Min(startIndex+chunkSize, len(tokens))
		chunkTokens := tokens[startIndex:endIndex]

		chunks[i] = Chunk{
			Index:      i,
			Content:    tokenizer.Decode(chunkTokens),
			TokenCount: endIndex - startIndex,
		}
	}

	return chunks, nil
}
