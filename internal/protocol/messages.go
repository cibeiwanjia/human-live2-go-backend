// Package protocol defines all messages and responses for the API
package protocol

import (
	"encoding/json"
)

// BaseMessage is the base for all message types
type BaseMessage struct{}

// AudioMessage represents audio data with metadata
type AudioMessage struct {
	Data        []byte     `json:"data,omitempty"`
	Type        AUDIO_TYPE `json:"type"`
	SampleRate  int        `json:"sampleRate"`
	SampleWidth int        `json:"sampleWidth"`
}

// TextMessage represents text data
type TextMessage struct {
	Data string `json:"data,omitempty"`
}

// RoleMessage represents a message with a role (for conversation history)
type RoleMessage struct {
	Role    ROLE_TYPE `json:"role"`
	Content string    `json:"content"`
}

// VoiceDesc describes a TTS voice
type VoiceDesc struct {
	Name   string      `json:"name"`
	Gender GENDER_TYPE `json:"gender"`
}

// ParamDesc describes an engine parameter
type ParamDesc struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        PARAM_TYPE  `json:"type"`
	Required    bool        `json:"required"`
	Range       []string    `json:"range,omitempty"`
	Choices     []string    `json:"choices,omitempty"`
	Default     interface{} `json:"default"`
}

// EngineDesc describes an engine
type EngineDesc struct {
	Name      string            `json:"name"`
	Type      ENGINE_TYPE       `json:"type"`
	InferType INFER_TYPE        `json:"infer_type"`
	Desc      string            `json:"desc"`
	Meta      map[string]string `json:"meta"`
}

// EngineConfig represents engine configuration
type EngineConfig struct {
	Name   string                 `json:"name"`
	Type   ENGINE_TYPE            `json:"type"`
	Config map[string]interface{} `json:"config"`
}

// UserDesc describes user information
type UserDesc struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
	Cookie    string `json:"cookie"`
}

// SSEEvent represents a Server-Sent Event
type SSEEvent struct {
	Type EVENT_TYPE `json:"-"`
	Data string     `json:"-"`
}

// ToSSEString converts the event to SSE format string
// Format: "event: {TYPE}\ndata: {DATA}\n\n"
func (e *SSEEvent) ToSSEString() string {
	return "event: " + string(e.Type) + "\ndata: " + escapeNewlines(e.Data) + "\n\n"
}

// escapeNewlines escapes newlines in data for SSE format
func escapeNewlines(s string) string {
	// Replace literal newlines with escaped newlines
	result := ""
	for _, c := range s {
		if c == '\n' {
			result += "\\n"
		} else {
			result += string(c)
		}
	}
	return result
}

// NewSSEEvent creates a new SSE event
func NewSSEEvent(eventType EVENT_TYPE, data string) *SSEEvent {
	return &SSEEvent{
		Type: eventType,
		Data: data,
	}
}

// SSE Event constructors
func SSEEventText(data string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_TEXT, data)
}

func SSEEventThink(data string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_THINK, data)
}

func SSEEventConversationID(id string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_CONVERSATION_ID, id)
}

func SSEEventMessageID(id string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_MESSAGE_ID, id)
}

func SSEEventTask(taskID string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_TASK, taskID)
}

func SSEEventDone() *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_DONE, "Done")
}

func SSEEventError(err string) *SSEEvent {
	return NewSSEEvent(EVENT_TYPE_ERROR, err)
}

// IsSSEEvent checks if a string is an SSE event format
func IsSSEEvent(s string) bool {
	return len(s) > 6 && s[:6] == "event:"
}

// MarshalJSON implements json.Marshaler for custom JSON output if needed
func (e *SSEEvent) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"event": string(e.Type),
		"data":  e.Data,
	})
}
