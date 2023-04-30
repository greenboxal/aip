package tokenizers

type Token interface {
	float32 | float64 | int | int64 | uint | uint64
}

type Tokenizable interface {
	CalculateTokenCount(tokenizer BasicTokenizer) (int, error)
}

type BasicTokenizer interface {
	Count(text string) (int, error)
}

type Tokenizer[T Token] interface {
	BasicTokenizer

	Encode(text string) ([]T, error)
	Decode(tokens []T) (string, error)
}
