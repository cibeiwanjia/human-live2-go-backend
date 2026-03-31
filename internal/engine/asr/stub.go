package asr

import (
	"context"
	"errors"

	"github.com/gorilla/websocket"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type StreamASREngine interface {
	Name() string
	Desc() protocol.EngineDesc
	Parameters() []protocol.ParamDesc
	RunStream(ctx context.Context, conn *websocket.Conn, config map[string]interface{}) error
}

type StubASREngine struct {
	base.BaseEngine
}

func NewFunASR(config map[string]interface{}) *StubASREngine {
	return &StubASREngine{
		BaseEngine: base.BaseEngine{
			Name_:      "FunASR",
			Desc_:      "FunASR Streaming ASR",
			Type_:      protocol.ENGINE_TYPE_ASR,
			InferType_: protocol.INFER_TYPE_STREAM,
		},
	}
}

func NewDifyASR(config map[string]interface{}) *StubASREngine {
	return &StubASREngine{
		BaseEngine: base.BaseEngine{
			Name_:      "DifyASR",
			Desc_:      "Dify ASR",
			Type_:      protocol.ENGINE_TYPE_ASR,
			InferType_: protocol.INFER_TYPE_STREAM,
		},
	}
}

func NewCozeASR(config map[string]interface{}) *StubASREngine {
	return &StubASREngine{
		BaseEngine: base.BaseEngine{
			Name_:      "CozeASR",
			Desc_:      "Coze ASR",
			Type_:      protocol.ENGINE_TYPE_ASR,
			InferType_: protocol.INFER_TYPE_STREAM,
		},
	}
}

func (e *StubASREngine) Run(ctx context.Context, input *protocol.AudioMessage, config map[string]interface{}) (*protocol.TextMessage, error) {
	return nil, errors.New("not implemented")
}

func (e *StubASREngine) RunStream(ctx context.Context, conn *websocket.Conn, config map[string]interface{}) error {
	return errors.New("not implemented")
}
