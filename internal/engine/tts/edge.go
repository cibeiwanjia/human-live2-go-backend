package tts

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type EdgeTTS struct {
	base.BaseEngine
	voices []protocol.VoiceDesc
}

func NewEdgeTTS() *EdgeTTS {
	return &EdgeTTS{
		BaseEngine: base.BaseEngine{
			Name_:      "EdgeTTS",
			Desc_:      "Microsoft Edge TTS",
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
		voices: edgeVoiceList,
	}
}

func (e *EdgeTTS) Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error) {
	return e.voices, nil
}

func (e *EdgeTTS) Run(ctx context.Context, input *protocol.TextMessage, config map[string]interface{}) (*protocol.AudioMessage, error) {
	voice := getStringConfig(config, "voice", "zh-CN-XiaoxiaoNeural")
	rate := getIntConfig(config, "rate", 0)
	volume := getIntConfig(config, "volume", 0)
	pitch := getIntConfig(config, "pitch", 0)

	audioData, err := e.synthesize(ctx, input.Data, voice, rate, volume, pitch)
	if err != nil {
		return nil, fmt.Errorf("edge tts synthesis failed: %w", err)
	}

	encodedData := base64.StdEncoding.EncodeToString(audioData)

	return &protocol.AudioMessage{
		Data:        []byte(encodedData),
		Type:        protocol.AUDIO_TYPE_MP3,
		SampleRate:  24000,
		SampleWidth: 2,
	}, nil
}

func (e *EdgeTTS) synthesize(ctx context.Context, text, voice string, rate, volume, pitch int) ([]byte, error) {
	wsURL := "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4"

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("websocket dial failed: %w", err)
	}
	defer conn.Close()

	if err := e.sendSSML(conn, text, voice, rate, volume, pitch); err != nil {
		return nil, err
	}

	audioData, err := e.receiveAudio(conn)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}

func (e *EdgeTTS) sendSSML(conn *websocket.Conn, text, voice string, rate, volume, pitch int) error {
	rateStr := fmt.Sprintf("%+d%%", rate)
	volumeStr := fmt.Sprintf("%+d%%", volume)
	pitchStr := fmt.Sprintf("%+dHz", pitch)

	ssml := fmt.Sprintf(
		`<speak version="1.0" xmlns="http://www.w3.org/2001/10/synthesis" xml:lang="en-US"><voice name="%s"><prosody pitch="%s" rate="%s" volume="%s">%s</prosody></voice></speak>`,
		voice, pitchStr, rateStr, volumeStr, escapeXML(text),
	)

	msg := fmt.Sprintf(
		"X-Timestamp:%s\r\nContent-Type:application/ssml+xml\r\nX-RequestId:%s\r\nPATH:speech.config\r\n\r\n%s",
		time.Now().UTC().Format("Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"),
		generateUUID(),
		ssml,
	)

	return conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func (e *EdgeTTS) receiveAudio(conn *websocket.Conn) ([]byte, error) {
	var audioData []byte

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			continue
		}

		if messageType == websocket.BinaryMessage {
			if len(data) > 2 && data[0] == 0x00 && data[1] == 0x67 {
				audioData = append(audioData, data[2:]...)
			}
		}
	}

	return audioData, nil
}

func getStringConfig(config map[string]interface{}, key string, defaultVal string) string {
	if v, ok := config[key].(string); ok {
		return v
	}
	return defaultVal
}

func getIntConfig(config map[string]interface{}, key string, defaultVal int) int {
	if v, ok := config[key].(int); ok {
		return v
	}
	if v, ok := config[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
