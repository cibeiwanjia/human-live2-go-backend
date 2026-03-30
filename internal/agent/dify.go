package agent

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type DifyAgent struct {
	BaseAgent
	apiServer string
	apiKey    string
	username  string
	client    *http.Client
}

type DifyConfig struct {
	APIServer string
	APIKey    string
	Username  string
}

func NewDifyAgent(cfg *DifyConfig) *DifyAgent {
	return &DifyAgent{
		BaseAgent: BaseAgent{
			name:      "DifyAgent",
			desc:      "Dify platform agent",
			inferType: protocol.INFER_TYPE_STREAM,
		},
		apiServer: cfg.APIServer,
		apiKey:    cfg.APIKey,
		username:  cfg.Username,
		client:    &http.Client{},
	}
}

func (a *DifyAgent) CreateConversation(ctx context.Context, config map[string]interface{}) (string, error) {
	payload := map[string]interface{}{
		"inputs":          map[string]interface{}{},
		"query":           "hello",
		"response_mode":   "blocking",
		"user":            a.username,
		"conversation_id": "",
		"files":           []interface{}{},
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", a.apiServer+"/chat-messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ConversationID string `json:"conversation_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.ConversationID, nil
}

func (a *DifyAgent) Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error) {
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
			"inputs":          map[string]interface{}{},
			"query":           req.Input,
			"response_mode":   "streaming",
			"user":            a.username,
			"conversation_id": conversationID,
			"files":           []interface{}{},
		}

		body, _ := json.Marshal(payload)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", a.apiServer+"/chat-messages", bytes.NewReader(body))
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

func (a *DifyAgent) parseSSEStream(reader io.Reader, ch chan<- *protocol.SSEEvent) {
	scanner := bufio.NewScanner(reader)
	pattern := regexp.MustCompile(`data:\s*({.*})`)
	messageID := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		match := pattern.FindStringSubmatch(line)
		if match == nil {
			continue
		}

		var data struct {
			Event          string `json:"event"`
			Answer         string `json:"answer"`
			ConversationID string `json:"conversation_id"`
			MessageID      string `json:"message_id"`
		}

		if err := json.Unmarshal([]byte(match[1]), &data); err != nil {
			continue
		}

		if messageID == "" && data.MessageID != "" {
			messageID = data.MessageID
			ch <- protocol.SSEEventMessageID(messageID)
		}

		if data.Event == "message" && data.Answer != "" {
			ch <- protocol.SSEEventText(data.Answer)
		}
	}
}
