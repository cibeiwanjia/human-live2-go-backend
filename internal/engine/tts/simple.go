package tts

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/engine/base"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

type SimpleTTSEngine struct {
	base.BaseEngine
}

func NewSimpleTTS(config map[string]interface{}) *SimpleTTSEngine {
	return &SimpleTTSEngine{
		BaseEngine: base.BaseEngine{
			Name_:      "SimpleTTS",
			Desc_:      "Simple Local TTS",
			Type_:      protocol.ENGINE_TYPE_TTS,
			InferType_: protocol.INFER_TYPE_NORMAL,
		},
	}
}

func (e *SimpleTTSEngine) Voices(ctx context.Context, config map[string]interface{}) ([]protocol.VoiceDesc, error) {
	voices := []protocol.VoiceDesc{
		{Name: "default", Gender: protocol.GENDER_TYPE_FEMALE},
	}
	return voices, nil
}

func (e *SimpleTTSEngine) Run(ctx context.Context, input *protocol.TextMessage, config map[string]interface{}) (*protocol.AudioMessage, error) {
	osType := getOS()

	var cmdName string
	switch osType {
	case "darwin":
		cmdName = "say"
	case "linux":
		cmdName = "espeak"
	default:
		return nil, fmt.Errorf("SimpleTTS is not supported on this operating system (%s). Please use EdgeTTS or TencentTTS instead", osType)
	}

	if _, err := exec.LookPath(cmdName); err != nil {
		return nil, fmt.Errorf("SimpleTTS requires '%s' command. Please install it first:\n  - macOS: pre-installed\n  - Linux: sudo apt-get install espeak\nOr use EdgeTTS/TencentTTS instead", cmdName)
	}

	tmpDir := os.TempDir()
	timestamp := time.Now().UnixNano()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("tts_%d.wav", timestamp))
	defer os.Remove(tmpFile)

	var cmd *exec.Cmd
	switch osType {
	case "darwin":
		cmd = exec.Command("say", "-v", "Ting-Ting", "-o", tmpFile, input.Data)
	case "linux":
		cmd = exec.Command("espeak", "-v", "zh", "-w", tmpFile, input.Data)
	}

	if err := cmd.Run(); err != nil {
		if osType == "darwin" {
			cmd = exec.Command("say", "-o", tmpFile, input.Data)
			if err := cmd.Run(); err != nil {
				return nil, fmt.Errorf("failed to generate TTS audio: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to generate TTS audio: %w", err)
		}
	}

	audioData, err := os.ReadFile(tmpFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read generated audio: %w", err)
	}

	return &protocol.AudioMessage{
		Data:        audioData,
		Type:        protocol.AUDIO_TYPE_WAV,
		SampleRate:  22050,
		SampleWidth: 2,
	}, nil
}

func getOS() string {
	switch os := os.Getenv("GOOS"); os {
	case "darwin":
		return "darwin"
	case "linux":
		return "linux"
	case "windows":
		return "windows"
	default:
		return detectOS()
	}
}

func detectOS() string {
	cmd := exec.Command("uname", "-s")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	osStr := string(output)
	if osStr == "Darwin\n" {
		return "darwin"
	} else if osStr == "Linux\n" {
		return "linux"
	}
	return "unknown"
}
