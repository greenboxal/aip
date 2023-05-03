package openai

import (
	"os"

	"github.com/sashabaranov/go-openai"
)

type Client = openai.Client
type ClientConfig = openai.ClientConfig

func NewClient() *Client {
	return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
}
