package openai

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Client struct {
	*openai.Client
}

type ClientConfig = openai.ClientConfig
type CompletionRequest = openai.CompletionRequest
type ChatCompletionRequest = openai.ChatCompletionRequest
type ImageRequest = openai.ImageRequest

func NewClient() *Client {
	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	return &Client{
		Client: c,
	}
}

func (c *Client) MaxTokensForModel(model string) int {
	switch model {
	case "gpt-3.5-turbo":
		return 4090
	case "gpt-3.5-turbo-16k":
		return 16384
	case "gpt-4":
		return 8192
	case "gpt-4-32k":
		return 32768
	default:
		panic("not implemented")
	}
}
