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
			Type_:      protocol.ENGINE_TYPE_TTS,
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
	// Use a different approach - Edge TTS via HTTP API instead of WebSocket
	// WebSocket connection is often blocked by Microsoft for server environments
	return e.synthesizeViaHTTP(ctx, text, voice, rate, volume, pitch)
}

func (e *EdgeTTS) synthesizeViaHTTP(ctx context.Context, text, voice string, rate, volume, pitch int) ([]byte, error) {
	// Edge TTS doesn't have a public HTTP API, so we need to use WebSocket
	// Let's try with more complete headers
	wsURL := "wss://speech.platform.bing.com/consumer/speech/synthesize/readaloud/edge/v1?trustedclienttoken=6A5AA1D4EAFF4E9FB37E23D68491D6F4"

	dialer := websocket.DefaultDialer
	dialer.EnableCompression = false
	dialer.HandshakeTimeout = 10 * time.Second

	headers := map[string][]string{
		"Origin":                 {"https://edge.microsoft.com"},
		"User-Agent":             {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"Accept":                 {"*/*"},
		"Accept-Encoding":        {"gzip, deflate, br"},
		"Accept-Language":        {"zh-CN,zh;q=0.9,en;q=0.8"},
		"Sec-WebSocket-Key":      {generateWSKey()},
		"Sec-WebSocket-Version":  {"13"},
		"Sec-WebSocket-Extensions": {"permessage-deflate; client_max_window_bits"},
	}

	conn, resp, err := dialer.DialContext(ctx, wsURL, headers)
	if err != nil {
		if resp != nil {
			return nil, fmt.Errorf("websocket dial failed (status %d): %w", resp.StatusCode, err)
		}
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

func generateWSKey() string {
	// Generate a random base64 string for Sec-WebSocket-Key
	b := make([]byte, 16)
	for i := range b {
		b[i] = byte(time.Now().UnixNano() % 256)
	}
	return base64.StdEncoding.EncodeToString(b)
}
