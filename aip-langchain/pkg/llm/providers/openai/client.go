package openai

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Client = openai.Client
type ClientConfig = openai.ClientConfig
type CompletionRequest = openai.CompletionRequest
type ChatCompletionRequest = openai.ChatCompletionRequest
type ImageRequest = openai.ImageRequest

func NewClient() *Client {
	return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
