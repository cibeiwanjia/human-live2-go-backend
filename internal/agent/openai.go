package agent

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type OpenAIAgent struct {
	BaseAgent
	client *openai.Client
	model  string
	config map[string]interface{}
}

type OpenAIConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

func NewOpenAIAgent(cfg *OpenAIConfig) *OpenAIAgent {
	config := openai.DefaultConfig(cfg.APIKey)
	if cfg.BaseURL != "" {
		config.BaseURL = cfg.BaseURL
	}

	return &OpenAIAgent{
		BaseAgent: BaseAgent{
			name:      "OpenAIAgent",
			desc:      "OpenAI compatible agent",
			inferType: protocol.INFER_TYPE_STREAM,
		},
		client: openai.NewClientWithConfig(config),
		model:  cfg.Model,
	}
}

func (a *OpenAIAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	return uuid.New().String(), nil
}

func (a *OpenAIAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
	ch := make(chan *protocol.SSEEvent, 20)

	go func() {
		defer close(ch)

		if req.ConversationID == "" {
			req.ConversationID = uuid.New().String()
			ch <- protocol.SSEEventConversationID(req.ConversationID)
		}

		messageID := uuid.New().String()
		ch <- protocol.SSEEventMessageID(messageID)

		messages := []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Input,
			},
		}

		stream, err := a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model:    a.model,
			Messages: messages,
			Stream:   true,
		})
		if err != nil {
			ch <- protocol.SSEEventError(err.Error())
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if errors.Is(err, context.Canceled) {
				break
			}
			if err != nil {
				ch <- protocol.SSEEventError(err.Error())
				return
			}

			if len(response.Choices) == 0 {
				continue
			}

			delta := response.Choices[0].Delta

			if delta.ReasoningContent != "" {
				ch <- protocol.SSEEventThink(delta.ReasoningContent)
			}

			if delta.Content != "" {
				ch <- protocol.SSEEventText(delta.Content)
			}
		}

		ch <- protocol.SSEEventDone()
	}()

	return ch, nil
}
