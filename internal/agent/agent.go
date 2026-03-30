// Package agent provides agent implementations and management
package agent

import (
	"context"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

// Agent defines the interface for all agent implementations
type Agent interface {
	Name() string
	Type() protocol.ENGINE_TYPE
	Desc() protocol.EngineDesc
	Parameters() []protocol.ParamDesc
	InferType() protocol.INFER_TYPE

	CreateConversation(ctx context.Context, config map[string]interface{}) (string, error)
	Run(ctx context.Context, req *AgentRequest) (<-chan *protocol.SSEEvent, error)
}

// AgentRequest represents a request to run an agent
type AgentRequest struct {
	Input          string
	ConversationID string
	Streaming      bool
	Config         map[string]interface{}
	User           *protocol.UserDesc
}

// BaseAgent provides common functionality for agents
type BaseAgent struct {
	name      string
	desc      string
	inferType protocol.INFER_TYPE
}

func (a *BaseAgent) Name() string {
	return a.name
}

func (a *BaseAgent) Type() protocol.ENGINE_TYPE {
	return protocol.ENGINE_TYPE_AGENT
}

func (a *BaseAgent) InferType() protocol.INFER_TYPE {
	return a.inferType
}

func (a *BaseAgent) Desc() protocol.EngineDesc {
	return protocol.EngineDesc{
		Name:      a.name,
		Type:      protocol.ENGINE_TYPE_AGENT,
		InferType: a.inferType,
		Desc:      a.desc,
		Meta:      make(map[string]string),
	}
}

func (a *BaseAgent) Parameters() []protocol.ParamDesc {
	return []protocol.ParamDesc{}
}
