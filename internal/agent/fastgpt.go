package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type FastGPTAgent struct {
	BaseAgent
	baseURL string
	apiKey  string
	uid     string
	client  *http.Client
}

type FastGPTConfig struct {
	BaseURL string
	APIKey  string
	UID     string
}

func NewFastGPTAgent(cfg *FastGPTConfig) *FastGPTAgent {
	return &FastGPTAgent{
		BaseAgent: BaseAgent{
			name:      "FastGPTAgent",
			desc:      "FastGPT platform agent",
			inferType: protocol.INFER_TYPE_STREAM,
		},
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
		uid:     cfg.UID,
		client:  &http.Client{},
	}
}

func (a *FastGPTAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	return uuid.New().String(), nil
}

func (a *FastGPTAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
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

		messageID := uuid.New().String()
		ch <- protocol.SSEEventMessageID(messageID)

		payload := map[string]interface{}{
			"chatId": conversationID,
			"stream": true,
			"detail": false,
			"messages": []map[string]interface{}{
				{
					"role":    "user",
					"content": req.Input,
				},
			},
			"customUid": a.uid,
		}

		body, _ := json.Marshal(payload)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/v1/chat/completions", bytes.NewReader(body))
		if err != nil {
			ch <- protocol.SSEEventError(err.Error())
			return
		}

		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+a.apiKey)

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

func (a *FastGPTAgent) parseSSEStream(reader io.Reader, ch chan<- *protocol.SSEEvent) {
	scanner := bufio.NewScanner(reader)
	pattern := regexp.MustCompile(`data:\s*({.*})`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.Contains(line, "DONE") {
			break
		}

		match := pattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		var data struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(match[1]), &data); err != nil {
			continue
		}

		if len(data.Choices) > 0 && data.Choices[0].Delta.Content != "" {
			ch <- protocol.SSEEventText(data.Choices[0].Delta.Content)
		}
	}
}
