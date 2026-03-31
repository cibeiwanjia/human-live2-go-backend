package tts

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tts/v20190823"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type TencentTTSEngine struct {
	base.BaseEngine
	appID     string
	secretID  string
	secretKey string
}

func NewTencentTTS(config map[string]interface{}) *TencentTTSEngine {
	appID := ""
	secretID := ""
	secretKey := ""

	if config != nil {
		if v, ok := config["app_id"].(string); ok {
			appID = v
		}
		if v, ok := config["secret_id"].(string); ok {
			secretID = v
		}
		if v, ok := config["secret_key"].(string); ok {
			secretKey = v
		}
	}

	return &TencentTTSEngine{
		BaseEngine: base.BaseEngine{
			Name_:      "TencentTTS",
			Desc_:      "Tencent Cloud TTS",
			Type_:      protocol.ENGINE_TYPE_TTS,
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
		appID:     appID,
		secretID:  secretID,
		secretKey: secretKey,
	}
}

func (e *TencentTTSEngine) Parameters() []protocol.ParamDesc {
	return []protocol.ParamDesc{
		{
			Name:        "app_id",
			Description: "Tencent Cloud App ID",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    true,
			Default:     "",
		},
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
			Name:        "voice",
			Description: "Voice type",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    false,
			Default:     "101001",
			Choices:     []string{"101001", "101002", "101003", "101004", "101005"},
		},
		{
			Name:        "speed",
			Description: "Speech speed (-2 to 2)",
			Type:        protocol.PARAM_TYPE_INT,
			Required:    false,
			Default:     0,
			Range:       []string{"-2", "2"},
		},
	}
}

func (e *TencentTTSEngine) Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error) {
	voices := []protocol.VoiceDesc{
		{Name: "101001", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101002", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101003", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101004", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101005", Gender: protocol.GENDER_TYPE_MALE},
	}
	return voices, nil
}

func (e *TencentTTSEngine) Run(ctx context.Context, input *protocol.TextMessage, config map[string]interface{}) (*protocol.AudioMessage, error) {
	appID := e.appID
	secretID := e.secretID
	secretKey := e.secretKey

	if config != nil {
		if v, ok := config["app_id"].(string); ok && v != "" {
			appID = v
		}
		if v, ok := config["secret_id"].(string); ok && v != "" {
			secretID = v
		}
		if v, ok := config["secret_key"].(string); ok && v != "" {
			secretKey = v
		}
	}

	if appID == "" || secretID == "" || secretKey == "" {
		return nil, fmt.Errorf("Tencent TTS credentials not configured. Please provide app_id, secret_id, and secret_key in frontend settings")
	}

	voice := getStringConfig(config, "voice", "101001")
	speed := getIntConfig(config, "speed", 0)

	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tts.tencentcloudapi.com"

	client, err := tts.NewClient(credential, "ap-shanghai", cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tencent TTS client: %w", err)
	}

	request := tts.NewTextToVoiceRequest()
	request.Text = common.StringPtr(input.Data)
	request.VoiceType = common.Int64Ptr(int64(getVoiceType(voice)))
	request.Speed = common.Float64Ptr(float64(speed))
	request.Codec = common.StringPtr("mp3")

	response, err := client.TextToVoice(request)
	if err != nil {
		return nil, fmt.Errorf("Tencent TTS API error: %w", err)
	}

	audioBase64 := *response.Response.Audio
	audioData, err := base64.StdEncoding.DecodeString(audioBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio: %w", err)
	}

	return &protocol.AudioMessage{
		Data:        audioData,
		Type:        protocol.AUDIO_TYPE_MP3,
		SampleRate:  16000,
		SampleWidth: 2,
	}, nil
}

func getVoiceType(voice string) int {
	voiceMap := map[string]int{
		"101001": 101001,
		"101002": 101002,
		"101003": 101003,
		"101004": 101004,
		"101005": 101005,
	}
	if v, ok := voiceMap[voice]; ok {
		return v
	}
	return 101001
}
