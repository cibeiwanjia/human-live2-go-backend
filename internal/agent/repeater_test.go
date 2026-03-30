package agent

import (
	"context"
	"testing"

	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/protocol"
)

func TestRepeaterAgent_Name(t *testing.T) {
	agent := NewRepeaterAgent()
	if agent.Name() != "RepeaterAgent" {
		t.Errorf("Name() = %v, want RepeaterAgent", agent.Name())
	}
}

func TestRepeaterAgent_Desc(t *testing.T) {
	agent := NewRepeaterAgent()
	desc := agent.Desc()

	if desc.Name != "RepeaterAgent" {
		t.Errorf("Desc().Name = %v, want RepeaterAgent", desc.Name)
	}
	if desc.Type != protocol.ENGINE_TYPE_AGENT {
		t.Errorf("Desc().Type = %v, want AGENT", desc.Type)
	}
}

func TestRepeaterAgent_CreateConversation(t *testing.T) {
	agent := NewRepeaterAgent()
	ctx := context.Background()

	id, err := agent.CreateConversation(ctx, nil)
	if err != nil {
		t.Errorf("CreateConversation() error = %v", err)
	}
	if id == "" {
		t.Error("CreateConversation() returned empty ID")
	}
}

func TestRepeaterAgent_Run(t *testing.T) {
	agent := NewRepeaterAgent()
	ctx := context.Background()

	req := &AgentRequest{
		Input:     "Hello World",
		Streaming: true,
	}

	events, err := agent.Run(ctx, req)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	eventList := make([]*protocol.SSEEvent, 0)
	for event := range events {
		eventList = append(eventList, event)
	}

	if len(eventList) < 3 {
		t.Errorf("Run() returned %d events, want at least 3", len(eventList))
	}

	foundText := false
	foundDone := false
	for _, event := range eventList {
		if event.Type == protocol.EVENT_TYPE_TEXT && event.Data == "Hello World" {
			foundText = true
		}
		if event.Type == protocol.EVENT_TYPE_DONE {
			foundDone = true
		}
	}

	if !foundText {
		t.Error("Run() did not return TEXT event with input")
	}
	if !foundDone {
		t.Error("Run() did not return DONE event")
	}
}
