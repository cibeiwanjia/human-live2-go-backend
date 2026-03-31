package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type CozeAgent struct {
	BaseAgent
	token  string
	botID  string
	client *http.Client
}

type CozeConfig struct {
	Token string
	BotID string
}

func NewCozeAgent(cfg *CozeConfig) *CozeAgent {
	return &CozeAgent{
		BaseAgent: BaseAgent{
			name:      "CozeAgent",
			desc:      "Coze platform agent",
			inferType: protocol.INFER_TYPE_STREAM,
		},
		token:  cfg.Token,
		botID:  cfg.BotID,
		client: &http.Client{},
	}
}

func (a *CozeAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.coze.cn/v1/conversation/create", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Data.ID, nil
}

func (a *CozeAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
	ch := make(chan *protocol.SSEEvent, 20)

	go func() {
		defer close(ch)

		conversationID := req.ConversationID
		if conversationID == "" {
			var err error
			conversationID, err = a.CreateConversation(ctx, req.Config)
			if err != nil {
				ch <- protocol.SSEEventError(err.Error())
				return
			}
			ch <- protocol.SSEEventConversationID(conversationID)
		}

		payload := map[string]interface{}{
			"bot_id":            a.botID,
			"user_id":           "adh",
			"stream":            true,
			"auto_save_history": true,
			"additional_messages": []map[string]interface{}{
				{
					"role":         "user",
					"content":      req.Input,
					"content_type": "text",
				},
			},
		}

		body, _ := json.Marshal(payload)
		apiURL := "https://api.coze.cn/v3/chat?conversation_id=" + conversationID
		httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
		if err != nil {
			ch <- protocol.SSEEventError(err.Error())
			return
		}

		httpReq.Header.Set("Authorization", "Bearer "+a.token)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := a.client.Do(httpReq)
		if err != nil {
			ch <- protocol.SSEEventError(err.Error())
			return
		}
		defer resp.Body.Close()

		a.parseSSEStream(resp.Body, ch)

		ch <- protocol.SSEEventDone()
	}()

	return ch, nil
}

func (a *CozeAgent) parseSSEStream(reader io.Reader, ch chan<- *protocol.SSEEvent) {
	scanner := bufio.NewScanner(reader)
	var event string

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "event:") {
			event = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}

		if event == "conversation.message.delta" && strings.HasPrefix(line, "data:") {
			dataStr := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if dataStr == "" {
				continue
			}

			var data struct {
				ReasoningContent string `json:"reasoning_content"`
				Content          string `json:"content"`
			}

			if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
				continue
			}

			if data.ReasoningContent != "" {
				ch <- protocol.SSEEventThink(data.ReasoningContent)
			}
			if data.Content != "" {
				ch <- protocol.SSEEventText(data.Content)
			}
		}
	}
}
