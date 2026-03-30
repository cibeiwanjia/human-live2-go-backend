package protocol

import (
	"testing"
)

func TestStructMessage(t *testing.T) {
	tests := []struct {
		name    string
		action  string
		payload []byte
		wantLen int
		wantErr bool
	}{
		{
			name:    "simple message",
			action:  "PING",
			payload: []byte{},
			wantLen: ProtocolHeaderSize,
			wantErr: false,
		},
		{
			name:    "message with payload",
			action:  "ENGINE_START",
			payload: []byte(`{"engine":"test"}`),
			wantLen: ProtocolHeaderSize + 17,
			wantErr: false,
		},
		{
			name:    "action too long",
			action:  "THIS_ACTION_IS_TOO_LONG",
			payload: []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := StructMessage(tt.action, tt.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(msg) != tt.wantLen {
				t.Errorf("StructMessage() len = %v, want %v", len(msg), tt.wantLen)
			}
		})
	}
}

func TestParseMessage(t *testing.T) {
	action := "PING"
	payload := []byte("test data")

	msg, err := StructMessage(action, payload)
	if err != nil {
		t.Fatalf("StructMessage() error = %v", err)
	}

	gotAction, gotPayload, err := ParseMessage(msg)
	if err != nil {
		t.Fatalf("ParseMessage() error = %v", err)
	}

	if gotAction != action {
		t.Errorf("ParseMessage() action = %v, want %v", gotAction, action)
	}

	if string(gotPayload) != string(payload) {
		t.Errorf("ParseMessage() payload = %v, want %v", string(gotPayload), string(payload))
	}
}

func TestSSEEvent(t *testing.T) {
	event := SSEEventText("Hello World")
	sse := event.ToSSEString()

	expected := "event: TEXT\ndata: Hello World\n\n"
	if sse != expected {
		t.Errorf("ToSSEString() = %q, want %q", sse, expected)
	}
}

func TestSSEEventWithNewlines(t *testing.T) {
	event := SSEEventText("line1\nline2")
	sse := event.ToSSEString()

	if sse != "event: TEXT\ndata: line1\\nline2\n\n" {
		t.Errorf("ToSSEString() with newlines = %q", sse)
	}
}
