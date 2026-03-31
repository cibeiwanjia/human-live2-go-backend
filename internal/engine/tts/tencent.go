package tts

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/google/uuid"
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
			Description: "Voice type (ID)",
			Type:        protocol.PARAM_TYPE_STRING,
			Required:    false,
			Default:     "502001",
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
		{Name: "502001", DisplayName: "智小柔", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "502003", DisplayName: "智小敏", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "502004", DisplayName: "智小满", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "502005", DisplayName: "智小解", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "502006", DisplayName: "智小悟", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "502007", DisplayName: "智小虎", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "602004", DisplayName: "暖心阿灿", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "602005", DisplayName: "专业梓欣", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "603000", DisplayName: "懂事少年", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "603001", DisplayName: "潇湘妹妹", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "603002", DisplayName: "软萌心心", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "603003", DisplayName: "随和老李", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "603004", DisplayName: "温柔小柠", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "603005", DisplayName: "知心大林", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "603006", DisplayName: "沉稳青叔", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "603007", DisplayName: "邻家女孩", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "602003", DisplayName: "爱小悠", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "501000", DisplayName: "智斌", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501001", DisplayName: "智兰", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "501002", DisplayName: "智菊", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "501003", DisplayName: "智宇", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501004", DisplayName: "月华", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "501005", DisplayName: "飞镜", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501006", DisplayName: "千嶂", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501007", DisplayName: "浅草", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501008", DisplayName: "WeJames", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "501009", DisplayName: "WeWinny", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "601008", DisplayName: "爱小豪", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "601009", DisplayName: "爱小芊", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "601010", DisplayName: "爱小娇", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "601011", DisplayName: "爱小川", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "601012", DisplayName: "爱小璟", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "601013", DisplayName: "爱小伊", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "601014", DisplayName: "爱小简", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101001", DisplayName: "智瑜", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101004", DisplayName: "智云", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101011", DisplayName: "智燕", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101013", DisplayName: "智辉", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101015", DisplayName: "智萌", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101016", DisplayName: "智甜", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101019", DisplayName: "智彤", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101021", DisplayName: "智瑞", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101026", DisplayName: "智希", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101027", DisplayName: "智梅", Gender: protocol.GENDER_TYPE_FEMALE},
		{Name: "101030", DisplayName: "智柯", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101050", DisplayName: "WeJack", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101054", DisplayName: "智友", Gender: protocol.GENDER_TYPE_MALE},
		{Name: "101055", DisplayName: "智付", Gender: protocol.GENDER_TYPE_FEMALE},
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

	voice := getStringConfig(config, "voice", "502001")
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
	request.SessionId = common.StringPtr(uuid.New().String())
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
		"502001": 502001, "502003": 502003, "502004": 502004, "502005": 502005,
		"502006": 502006, "502007": 502007,
		"602003": 602003, "602004": 602004, "602005": 602005,
		"603000": 603000, "603001": 603001, "603002": 603002, "603003": 603003,
		"603004": 603004, "603005": 603005, "603006": 603006, "603007": 603007,
		"501000": 501000, "501001": 501001, "501002": 501002, "501003": 501003,
		"501004": 501004, "501005": 501005, "501006": 501006, "501007": 501007,
		"501008": 501008, "501009": 501009,
		"601008": 601008, "601009": 601009, "601010": 601010, "601011": 601011,
		"601012": 601012, "601013": 601013, "601014": 601014,
		"101001": 101001, "101004": 101004, "101011": 101011, "101013": 101013,
		"101015": 101015, "101016": 101016, "101019": 101019, "101021": 101021,
		"101026": 101026, "101027": 101027, "101030": 101030, "101050": 101050,
		"101054": 101054, "101055": 101055,
	}
	nameToId := map[string]int{
		"智小柔": 502001, "智小敏": 502003, "智小满": 502004, "智小解": 502005,
		"智小悟": 502006, "智小虎": 502007,
		"暖心阿灿": 602004, "专业梓欣": 602005,
		"懂事少年": 603000, "潇湘妹妹": 603001, "软萌心心": 603002, "随和老李": 603003,
		"温柔小柠": 603004, "知心大林": 603005, "沉稳青叔": 603006, "邻家女孩": 603007,
		"爱小悠": 602003,
		"智斌":  501000, "智兰": 501001, "智菊": 501002, "智宇": 501003,
		"月华": 501004, "飞镜": 501005, "千嶂": 501006, "浅草": 501007,
		"WeJames": 501008, "WeWinny": 501009,
		"爱小豪": 601008, "爱小芊": 601009, "爱小娇": 601010, "爱小川": 601011,
		"爱小璟": 601012, "爱小伊": 601013, "爱小简": 601014,
		"智瑜": 101001, "智云": 101004, "智燕": 101011, "智辉": 101013,
		"智萌": 101015, "智甜": 101016, "智彤": 101019, "智瑞": 101021,
		"智希": 101026, "智梅": 101027, "智柯": 101030, "WeJack": 101050,
		"智友": 101054, "智付": 101055,
	}
	if v, ok := voiceMap[voice]; ok {
		return v
	}
	if v, ok := nameToId[voice]; ok {
		return v
	}
	return 502001
}
