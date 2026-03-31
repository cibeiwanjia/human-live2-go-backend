package handlers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/agent"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

// AgentHandler handles agent API requests
type AgentHandler struct {
	pool *agent.AgentPool
}

// NewAgentHandler creates a new AgentHandler
func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		pool: agent.GetPool(),
	}
}

// GetEngineList returns list of available agents
func (h *AgentHandler) GetEngineList(c *gin.Context) {
	names := h.pool.List()
	engines := make([]protocol.EngineDesc, 0, len(names))

	for _, name := range names {
		ag, err := h.pool.Get(name)
		if err != nil {
			continue
		}
		engines = append(engines, ag.Desc())
	}

	resp := protocol.NewEngineListResp(engines)
	c.JSON(200, resp)
}

// GetDefaultEngine returns the default agent
func (h *AgentHandler) GetDefaultEngine(c *gin.Context) {
	defaultName := h.pool.Default()
	if defaultName == "" {
		c.JSON(404, protocol.NewErrorResponse("no default agent configured"))
		return
	}
	ag, err := h.pool.Get(defaultName)
	if err != nil {
		c.JSON(404, protocol.NewErrorResponse("default agent not found"))
		return
	}

	resp := protocol.NewEngineDefaultResp(ag.Desc())
	c.JSON(200, resp)
}

// GetEngineParams returns parameters for a specific agent
func (h *AgentHandler) GetEngineParams(c *gin.Context) {
	engineName := c.Param("engine")

	ag, err := h.pool.Get(engineName)
	if err != nil {
		c.JSON(404, protocol.NewErrorResponse("agent not found"))
		return
	}

	params := ag.Parameters()
	resp := protocol.NewEngineParamResp(params)
	c.JSON(200, resp)
}

// CreateConversation creates a new conversation
func (h *AgentHandler) CreateConversation(c *gin.Context) {
	engineName := c.Param("engine")

	ag, err := h.pool.Get(engineName)
	if err != nil {
		c.JSON(404, protocol.NewErrorResponse("agent not found"))
		return
	}

	var input struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, protocol.NewErrorResponse("invalid request body"))
		return
	}

	conversationID, err := ag.CreateConversation(context.Background(), input.Data)
	if err != nil {
		c.JSON(500, protocol.NewErrorResponse(err.Error()))
		return
	}

	resp := protocol.NewConversationIdResp(conversationID)
	c.JSON(200, resp)
}

// AgentEngineInput represents agent inference request
type AgentEngineInput struct {
	Engine         string                 `json:"engine"`
	Config         map[string]interface{} `json:"config"`
	Data           string                 `json:"data"`
	ConversationID string                 `json:"conversation_id"`
}

// StreamInfer handles streaming agent inference
func (h *AgentHandler) StreamInfer(c *gin.Context) {
	var input AgentEngineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, protocol.NewErrorResponse("invalid request body"))
		return
	}

	ag, err := h.pool.Get(input.Engine)
	if err != nil {
		c.JSON(404, protocol.NewErrorResponse("agent not found"))
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	flusher, ok := c.Writer.(interface{ Flush() })
	if !ok {
		c.JSON(500, protocol.NewErrorResponse("streaming not supported"))
		return
	}

	req := &agent.AgentRequest{
		Input:          input.Data,
		ConversationID: input.ConversationID,
		Streaming:      true,
		Config:         input.Config,
	}

	events, err := ag.Run(c.Request.Context(), req)
	if err != nil {
		c.Writer.Write([]byte(protocol.SSEEventError(err.Error()).ToSSEString()))
		flusher.Flush()
		return
	}

	for event := range events {
		_, err := c.Writer.Write([]byte(event.ToSSEString()))
		if err != nil {
			return
		}
		flusher.Flush()
	}
}

func generateUUID() string {
	return uuid.New().String()
}
