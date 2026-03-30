// Package protocol defines all messages and responses for the API
package protocol

import (
	"encoding/json"
)

// BaseResponse is the standard API response wrapper
type BaseResponse struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
}

// WithData creates a response with data
func (r *BaseResponse) WithData(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:    r.Code,
		Message: r.Message,
	}
}

// Error sets error message and code
func (r *BaseResponse) Error(msg string) {
	r.Code = RESPONSE_CODE_ERROR
	r.Message = msg
}

// Success sets success code
func (r *BaseResponse) Success() {
	r.Code = RESPONSE_CODE_OK
	r.Message = "success"
}

// NewResponse creates a new response
func NewResponse() *BaseResponse {
	return &BaseResponse{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
	}
}

// NewErrorResponse creates a new error response
func NewErrorResponse(msg string) *BaseResponse {
	return &BaseResponse{
		Code:    RESPONSE_CODE_ERROR,
		Message: msg,
	}
}

// EngineListResp is response for engine list
type EngineListResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    []EngineDesc  `json:"data"`
}

// NewEngineListResp creates engine list response
func NewEngineListResp(engines []EngineDesc) *EngineListResp {
	if engines == nil {
		engines = []EngineDesc{}
	}
	return &EngineListResp{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    engines,
	}
}

// EngineDefaultResp is response for default engine
type EngineDefaultResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    EngineDesc    `json:"data"`
}

// NewEngineDefaultResp creates default engine response
func NewEngineDefaultResp(engine EngineDesc) *EngineDefaultResp {
	return &EngineDefaultResp{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    engine,
	}
}

// EngineParamResp is response for engine parameters
type EngineParamResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    []ParamDesc   `json:"data"`
}

// NewEngineParamResp creates engine parameter response
func NewEngineParamResp(params []ParamDesc) *EngineParamResp {
	if params == nil {
		params = []ParamDesc{}
	}
	return &EngineParamResp{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    params,
	}
}

// VoiceListResp is response for voice list
type VoiceListResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    []VoiceDesc   `json:"data"`
}

// NewVoiceListResp creates voice list response
func NewVoiceListResp(voices []VoiceDesc) *VoiceListResp {
	if voices == nil {
		voices = []VoiceDesc{}
	}
	return &VoiceListResp{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    voices,
	}
}

// ConversationIdResp is response for conversation ID
type ConversationIdResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    string        `json:"data"`
}

// NewConversationIdResp creates conversation ID response
func NewConversationIdResp(id string) *ConversationIdResp {
	return &ConversationIdResp{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    id,
	}
}

// ASREngineOutput is response for ASR inference
type ASREngineOutput struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    string        `json:"data"`
}

// NewASREngineOutput creates ASR output response
func NewASREngineOutput(text string) *ASREngineOutput {
	return &ASREngineOutput{
		Code:    RESPONSE_CODE_OK,
		Message: "success",
		Data:    text,
	}
}

// TTSEngineOutput is response for TTS inference
type TTSEngineOutput struct {
	Code        RESPONSE_CODE `json:"code"`
	Message     string        `json:"message"`
	Data        string        `json:"data"`
	SampleRate  int           `json:"sampleRate"`
	SampleWidth int           `json:"sampleWidth"`
}

// NewTTSEngineOutput creates TTS output response
func NewTTSEngineOutput(audioData string, sampleRate, sampleWidth int) *TTSEngineOutput {
	return &TTSEngineOutput{
		Code:        RESPONSE_CODE_OK,
		Message:     "success",
		Data:        audioData,
		SampleRate:  sampleRate,
		SampleWidth: sampleWidth,
	}
}

// StringResp is generic string response
type StringResp struct {
	Code    RESPONSE_CODE `json:"code"`
	Message string        `json:"message"`
	Data    string        `json:"data"`
}

// MarshalJSON implements custom JSON marshaling for BaseResponse
func (r *BaseResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"code":    r.Code,
		"message": r.Message,
	})
}
