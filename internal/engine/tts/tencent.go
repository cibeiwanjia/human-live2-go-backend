package tts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type TencentTTSEngine struct {
	base.BaseEngine
	appID     string
	secretID  string
	secretKey string
}

type tencentTTSRequest struct {
	Text      string `json:"Text"`
	SessionID string `json:"SessionId"`
	ModelType int    `json:"ModelType"`
	Speed     int    `json:"Speed"`
	VoiceType string `json:"VoiceType"`
	Codec     string `json:"Codec"`
}

type tencentTTSResponse struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
	Audio   string `json:"Audio"`
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

func (e *TencentTTSEngine) Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error) {
	// Return common Chinese voices
	voices := []protocol.VoiceDesc{
		{Name: "ZhiMei", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "ZhiYu", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "RuiXin", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "YunJie", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "YunXi", Gender: protocol.GENDER_TYPE_MALE},
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
		return nil, fmt.Errorf("Tencent TTS credentials not configured. Please provide app_id, secret_id, and secret_key in frontend settings or config file")
	}

	voice := getStringConfig(config, "voice", "ZhiMei")
	speed := getIntConfig(config, "speed", 0)

	// Build request
	reqBody := tencentTTSRequest{
		Text:      input.Data,
		SessionID: fmt.Sprintf("%d", time.Now().UnixNano()),
		ModelType: 1, // 1 for standard, 2 for premium
		Speed:     speed,
		VoiceType: voice,
		Codec:     "mp3",
	}

	if _, err := json.Marshal(reqBody); err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Note: This is a simplified implementation
	// In production, you would need to:
	// 1. Sign the request with Tencent Cloud signature
	// 2. Call the actual Tencent TTS API endpoint (e.g., https://tts.cloud.tencent.com)
	// 3. Handle the response and decode the audio

	// For now, return a placeholder error indicating the implementation is incomplete
	return nil, fmt.Errorf("Tencent TTS requires proper API integration with Tencent Cloud credentials. Please configure app_id, secret_id, and secret_key in your config")
}
