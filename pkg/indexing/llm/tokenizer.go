package llm

import "github.com/pkoukk/tiktoken-go"

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

type TikTokenTokenizer struct {
	tokenizer *tiktoken.Tiktoken
}

func (t *TikTokenTokenizer) Count(text string) (int, error) {
	data := t.tokenizer.Encode(text, nil, nil)

	return len(data), nil
}

func TikTokenForModel(model string) (*TikTokenTokenizer, error) {
	tokenizer, err := tiktoken.EncodingForModel(model)

	if err != nil {
		return nil, err
	}

	return &TikTokenTokenizer{
		tokenizer: tokenizer,
	}, nil
}

func (t *TikTokenTokenizer) Encode(text string) ([]int, error) {
	return t.tokenizer.Encode(text, nil, nil), nil
}

func (t *TikTokenTokenizer) Decode(tokens []int) (string, error) {
	return t.tokenizer.Decode(tokens), nil
}
