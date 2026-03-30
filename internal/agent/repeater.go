package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

// RepeaterAgent simply repeats user input
type RepeaterAgent struct {
	BaseAgent
}

// NewRepeaterAgent creates a new RepeaterAgent
func NewRepeaterAgent() *RepeaterAgent {
	return &RepeaterAgent{
		BaseAgent: BaseAgent{
			name:      "RepeaterAgent",
			desc:      "Repeat user input",
			inferType: protocol.INFER_TYPE_STREAM,
		},
	}
}

// CreateConversation creates a new conversation
func (a *RepeaterAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	return uuid.New().String(), nil
}

// Run executes the agent and returns events channel
func (a *RepeaterAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
	ch := make(chan *protocol.SSEEvent, 10)

	go func() {
		defer close(ch)

		if req.ConversationID == "" {
			req.ConversationID = uuid.New().String()
			ch <- protocol.SSEEventConversationID(req.ConversationID)
		}

		messageID := uuid.New().String()
		ch <- protocol.SSEEventMessageID(messageID)

		ch <- protocol.SSEEventText(req.Input)

		ch <- protocol.SSEEventDone()
	}()

	return ch, nil
}
