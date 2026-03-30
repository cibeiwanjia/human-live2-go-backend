package base

import "github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"

type BaseEngine struct {
	Name_      string
	Desc_      string
	InferType_ protocol.INFER_TYPE
}

func (e *BaseEngine) Name() string {
	return e.Name_
}

func (e *BaseEngine) Type() protocol.ENGINE_TYPE {
	return protocol.ENGINE_TYPE_TTS
}

func (e *BaseEngine) InferType() protocol.INFER_TYPE {
	return e.InferType_
}

func (e *BaseEngine) Desc() protocol.EngineDesc {
	return protocol.EngineDesc{
		Name:      e.Name_,
		Type:      protocol.ENGINE_TYPE_TTS,
		InferType: e.InferType_,
		Desc:      e.Desc_,
		Meta:      make(map[string]string),
	}
}

func (e *BaseEngine) Parameters() []protocol.ParamDesc {
	return []protocol.ParamDesc{}
}
