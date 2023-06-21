package openai

import (
	"context"

	"github.com/sashabaranov/go-openai"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm"
	"github.com/greenboxal/aip/aip-langchain/pkg/llm/chat"
	"github.com/greenboxal/aip/aip-langchain/pkg/tokenizers"
)

type LanguageModel struct {
	Client *Client

	Model       string
	Temperature float32
}

func (lm *LanguageModel) MaxTokens() int {
	return lm.Client.MaxTokensForModel(lm.Model)
}

func (lm *LanguageModel) Predict(ctx context.Context, prompt string, options ...llm.PredictOption) (string, error) {
	opts := llm.NewPredictOptions(options...)

	if opts.AutoMaxTokens {
		tokenizer := tokenizers.TikTokenForModel(lm.Model)
		count, err := tokenizer.Count(prompt)

		if err != nil {
			return "", err
		}

		opts.MaxTokens = lm.MaxTokens() - count
	}

	result, err := lm.Client.CreateCompletion(ctx, openai.CompletionRequest{
		Model:       lm.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Prompt:      prompt,
	})

	if err != nil {
		return "", err
	}

	return result.Choices[0].Text, nil
}

func (lm *LanguageModel) PredictChat(ctx context.Context, msg chat.Message, options ...llm.PredictOption) (chat.Message, error) {
	result, err := lm.Predict(ctx, msg.String(), options...)

	if err != nil {
		return chat.Message{}, err
	}

	return chat.Message{
		Entries: []chat.MessageEntry{
			{
				Role: msn.RoleAI,
				Text: result,
			},
		},
	}, nil
}

type ChatLanguageModel struct {
	Client *Client

	Model       string
	Temperature float32
}

func (lm *ChatLanguageModel) MaxTokens() int {
	return lm.Client.MaxTokensForModel(lm.Model)
}

func (lm *ChatLanguageModel) Predict(ctx context.Context, prompt string, options ...llm.PredictOption) (string, error) {
	msg := chat.Compose(
		chat.Entry(msn.RoleUser, prompt),
		chat.Entry(msn.RoleAI, ""),
	)

	result, err := lm.PredictChat(ctx, msg, options...)

	if err != nil {
		return "", nil
	}

	return result.Entries[0].Text, nil
}

func (lm *ChatLanguageModel) PredictChat(ctx context.Context, msg chat.Message, options ...llm.PredictOption) (chat.Message, error) {
	opts := llm.NewPredictOptions(options...)
	request, err := lm.buildChatCompletionRequest(msg, opts)

	if err != nil {
		return chat.Message{}, err
	}

	result, err := lm.Client.CreateChatCompletion(ctx, *request)

	if err != nil {
		return chat.Message{}, nil
	}

	return buildMsnMessages(result.Choices), nil
}

type messageStream struct {
	stream *openai.ChatCompletionStream
}

func (m *messageStream) Recv() (chat.MessageFragment, error) {
	reply, err := m.stream.Recv()

	if err != nil {
		return chat.MessageFragment{}, err
	}

	return chat.MessageFragment{
		MessageIndex: reply.Choices[0].Index,
		Delta:        reply.Choices[0].Delta.Content,
	}, nil
}

func (m *messageStream) Close() error {
	return m.Close()
}

func (lm *ChatLanguageModel) PredictChatStream(ctx context.Context, msg chat.Message, options ...llm.PredictOption) (chat.MessageStream, error) {
	opts := llm.NewPredictOptions(options...)
	request, err := lm.buildChatCompletionRequest(msg, opts)

	if err != nil {
		return nil, err
	}

	stream, err := lm.Client.CreateChatCompletionStream(ctx, *request)

	if err != nil {
		return nil, err
	}

	return &messageStream{stream: stream}, nil
}

func (lm *ChatLanguageModel) buildChatCompletionRequest(msg chat.Message, opts llm.PredictOptions) (*openai.ChatCompletionRequest, error) {
	messages := buildChatMessages(msg)

	if opts.AutoMaxTokens {
		tokenizer := tokenizers.TikTokenForModel(lm.Model)
		count, err := tokenizer.Count(msg.String())

		if err != nil {
			return nil, err
		}

		opts.MaxTokens = lm.MaxTokens() - count - len(messages)*4
	}

	if opts.MaxTokens <= 0 {
		opts.MaxTokens = lm.MaxTokens() / 2
	}

	req := &openai.ChatCompletionRequest{
		Model:       lm.Model,
		Temperature: opts.Temperature,
		MaxTokens:   opts.MaxTokens,
		Stop:        opts.StopSequences,
		Messages:    messages,
	}

	if opts.Functions != nil {
		req.Functions = make([]*openai.FunctionDefine, 0, len(opts.Functions))

		for _, v := range opts.Functions {

			req.Functions = append(req.Functions, &openai.FunctionDefine{
				Name:        v.Name,
				Description: v.Description,
				Parameters: &openai.FunctionParams{
					Type:       v.Parameters.Type,
					Required:   v.Parameters.Required,
					Properties: v.Parameters.Properties,
				},
			})
		}
	}

	if len(req.Functions) > 0 {
		if opts.AllowFunctionCall {
			if opts.ForceFunctionCall != nil {
				panic("not supported by the OpenAI client yet")
			} else {
				req.FunctionCall = "auto"
			}
		} else {
			req.FunctionCall = "none"
		}
	}

	return req, nil
}

func buildMsnMessages(choices []openai.ChatCompletionChoice) chat.Message {
	entries := make([]chat.MessageEntry, len(choices))

	for i, c := range choices {
		var role msn.Role

		switch c.Message.Role {
		case "system":
			role = msn.RoleSystem
		case "user":
			role = msn.RoleUser
		case "assistant":
			role = msn.RoleAI
		case "function":
			role = msn.RoleFunction
		default:
			panic("unknown role")
		}

		msg := chat.MessageEntry{
			Role: role,
			Name: c.Message.Name,
			Text: c.Message.Content,
		}

		if c.Message.FunctionCall != nil {
			msg.Fn = c.Message.FunctionCall.Name
			msg.FnArgs = c.Message.FunctionCall.Arguments
		}

		entries[i] = msg
	}

	return chat.Compose(entries...)
}

func buildChatMessages(msg chat.Message) []openai.ChatCompletionMessage {
	messages := make([]openai.ChatCompletionMessage, len(msg.Entries))

	for i, m := range msg.Entries {
		var role string

		switch m.Role {
		case msn.RoleSystem:
			role = "system"
		case msn.RoleUser:
			role = "user"
		case msn.RoleAI:
			role = "assistant"
		case msn.RoleFunction:
			role = "function"
		default:
			panic("unknown role")
		}

		msg := openai.ChatCompletionMessage{
			Name:    m.Name,
			Role:    role,
			Content: m.Text,
		}

		if m.Fn != "" {
			msg.FunctionCall = &openai.FunctionCall{
				Name:      m.Fn,
				Arguments: m.FnArgs,
			}
		}

		messages[i] = msg
	}

	return messages
}
