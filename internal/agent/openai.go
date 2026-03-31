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

func (a *OpenAIAgent) Parameters() []protocol.ParamDesc {
	return []protocol.ParamDesc{
		{
			Name:        "api_key",
			Description: "OpenAI API Key",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    true,
			Default:     "",
		},
		{
			Name:        "base_url",
			Description: "OpenAI API Base URL (optional, for custom endpoints)",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    false,
			Default:     "",
		},
		{
			Name:        "model",
			Description: "Model name (e.g., gpt-3.5-turbo, gpt-4)",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    false,
			Default:     "gpt-3.5-turbo",
		},
	}
}

func (a *OpenAIAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	return uuid.New().String(), nil
}

func (a *OpenAIAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
	ch := make(chan *protocol.SSEEvent, 20)

	apiKey := ""
	baseURL := ""
	model := a.model

	if req.Config != nil {
		if v, ok := req.Config["api_key"].(string); ok && v != "" {
			apiKey = v
		}
		if v, ok := req.Config["base_url"].(string); ok && v != "" {
			baseURL = v
		}
		if v, ok := req.Config["model"].(string); ok && v != "" {
			model = v
		}
	}

	if apiKey == "" && a.client == nil {
		return nil, errors.New("OpenAI API key not configured. Please provide api_key in frontend settings")
	}

	var client *openai.Client
	if apiKey != "" {
		config := openai.DefaultConfig(apiKey)
		if baseURL != "" {
			config.BaseURL = baseURL
		}
		client = openai.NewClientWithConfig(config)
	} else {
		client = a.client
	}

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

		stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
			Model:    model,
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
