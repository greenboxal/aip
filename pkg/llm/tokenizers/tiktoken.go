package tokenizers

import "github.com/pkoukk/tiktoken-go"

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
