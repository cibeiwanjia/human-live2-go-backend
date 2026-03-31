package engine

import (
	"context"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type Engine interface {
	Name() string
	Type() protocol.ENGINE_TYPE
	Desc() protocol.EngineDesc
	Parameters() []protocol.ParamDesc
	InferType() protocol.INFER_TYPE
}

type TTSEngine interface {
	Engine
	Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error)
	Run(ctx context.Context, input *protocol.TextMessage, config map[string]interface{}) (*protocol.AudioMessage, error)
}

type ASREngine interface {
	Engine
	Run(ctx context.Context, input *protocol.AudioMessage, config map[string]interface{}) (*protocol.TextMessage, error)
}
