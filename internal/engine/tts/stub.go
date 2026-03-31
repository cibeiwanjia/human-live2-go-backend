package tts

import (
	"context"
	"errors"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type StubTTSEngine struct {
	base.BaseEngine
}


func (e *StubTTSEngine) Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error) {
	return []protocol.VoiceDesc{}, nil
}

func (e *StubTTSEngine) Run(ctx context.Context, input *protocol.TextMessage, config map[string]interface{}) (*protocol.AudioMessage, error) {
	return nil, errors.New("not implemented")
}

func NewDifyTTS(config map[string]interface{}) *StubTTSEngine {
	return &StubTTSEngine{
		BaseEngine: base.BaseEngine{
			Name_:      "DifyTTS",
			Desc_:      "Dify TTS",
			Type_:      protocol.ENGINE_TYPE_TTS,
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
	}
}

func NewCozeTTS(config map[string]interface{}) *StubTTSEngine {
	return &StubTTSEngine{
		BaseEngine: base.BaseEngine{
			Name_:      "CozeTTS",
			Desc_:      "Coze TTS",
			Type_:      protocol.ENGINE_TYPE_TTS,
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
	}
}
