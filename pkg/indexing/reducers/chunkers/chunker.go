package chunkers

import "github.com/pkoukk/tiktoken-go"

func SplitTextIntoChunks(text string, chunkSize int, overlapSize int) ([]string, error) {
	// Tokenize the text using tiktoken
	tokenizer, err := tiktoken.EncodingForModel("text-embedding-ada-002")

	if err != nil {
		return nil, err
	}

	tokens := tokenizer.Encode(text, nil, nil)

	// Calculate the number of chunks based on the chunk size and overlap size
	numChunks := (len(tokens) + chunkSize - 1) / chunkSize

	chunks := make([]string, numChunks)

	// Generate the chunks by iterating over the tokens
	for i := 0; i < numChunks; i++ {
		startIndex := i * (chunkSize - overlapSize)
		endIndex := Min(startIndex+chunkSize, len(tokens))
		chunkTokens := tokens[startIndex:endIndex]

		chunks[i] = tokenizer.Decode(chunkTokens)
	}

	return chunks, nil
}

func Min(i int, i2 int) int {
	if i < i2 {
		return i
	}

	return i2
}
