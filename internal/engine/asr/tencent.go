package asr

import (
	"context"
	"encoding/base64"
	"fmt"

	tencentASR "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/asr/v20190614"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type TencentASREngine struct {
	base.BaseEngine
	secretID  string
	secretKey string
}

func NewTencentASR(config map[string]interface{}) *TencentASREngine {
	secretID := ""
	secretKey := ""

	if config != nil {
		if v, ok := config["secret_id"].(string); ok {
			secretID = v
		}
		if v, ok := config["secret_key"].(string); ok {
			secretKey = v
		}
	}

	return &TencentASREngine{
		BaseEngine: base.BaseEngine{
			Name_:      "TencentASR",
			Desc_:      "Tencent Cloud ASR",
			Type_:      protocol.ENGINE_TYPE_ASR,
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
		secretID:  secretID,
		secretKey: secretKey,
	}
}

func (e *TencentASREngine) Parameters() []protocol.ParamDesc {
	return []protocol.ParamDesc{
		{
			Name:        "secret_id",
			Description: "Tencent Cloud Secret ID",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    true,
			Default:     "",
		},
		{
			Name:        "secret_key",
			Description: "Tencent Cloud Secret Key",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    true,
			Default:     "",
		},
		{
			Name:        "engine_model_type",
			Description: "Engine model type (16k_zh, 16k_en, etc.)",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    false,
			Default:     "16k_zh",
		},
	}
}

func (e *TencentASREngine) Run(ctx context.Context, input *protocol.AudioMessage, config map[string]interface{}) (*protocol.TextMessage, error) {
	secretID := e.secretID
	secretKey := e.secretKey

	if config != nil {
		if v, ok := config["secret_id"].(string); ok && v != "" {
			secretID = v
		}
		if v, ok := config["secret_key"].(string); ok && v != "" {
			secretKey = v
		}
	}

	if secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("Tencent ASR credentials not configured")
	}

	engineModelType := "16k_zh"
	if config != nil {
		if v, ok := config["engine_model_type"].(string); ok && v != "" {
			engineModelType = v
		}
	}

	voiceFormat := "mp3"
	switch input.Type {
	case protocol.AUDIO_TYPE_WAV:
		voiceFormat = "wav"
	case protocol.AUDIO_TYPE_MP3:
		voiceFormat = "mp3"
	}

	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "asr.tencentcloudapi.com"

	client, err := tencentASR.NewClient(credential, "ap-shanghai", cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tencent ASR client: %w", err)
	}

	request := tencentASR.NewSentenceRecognitionRequest()
	request.EngSerViceType = common.StringPtr(engineModelType)
	request.SourceType = common.Uint64Ptr(1)
	request.VoiceFormat = common.StringPtr(voiceFormat)
	request.Data = common.StringPtr(base64.StdEncoding.EncodeToString(input.Data))
	request.DataLen = common.Int64Ptr(int64(len(input.Data)))

	response, err := client.SentenceRecognition(request)
	if err != nil {
		return nil, fmt.Errorf("Tencent ASR API error: %w", err)
	}

	return &protocol.TextMessage{
		Data: *response.Response.Result,
	}, nil
}
