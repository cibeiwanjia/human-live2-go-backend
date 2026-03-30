// Package protocol defines WebSocket binary protocol
package protocol

import (
	"encoding/binary"
	"errors"
)

const (
	ActionHeaderSize   = 18
	ProtocolHeaderSize = 22 // 18 + 4
)

var (
	ErrMessageTooShort     = errors.New("message too short")
	ErrMessageSizeMismatch = errors.New("message size mismatch")
	ErrActionNameTooLong   = errors.New("action name exceeds 18 bytes")
)

// StructMessage creates a binary WebSocket message
// Format: [Action(18B)] + [Payload Size(4B)] + [Payload(Variable)]
func StructMessage(action string, payload []byte) ([]byte, error) {
	if len(action) > ActionHeaderSize {
		return nil, ErrActionNameTooLong
	}

	actionBytes := formatAction(action)
	payloadSize := uint32(len(payload))

	header := make([]byte, ProtocolHeaderSize)
	copy(header[0:ActionHeaderSize], actionBytes)
	binary.BigEndian.PutUint32(header[ActionHeaderSize:ProtocolHeaderSize], payloadSize)

	return append(header, payload...), nil
}

// ParseMessage parses a binary WebSocket message
func ParseMessage(data []byte) (action string, payload []byte, err error) {
	if len(data) < ProtocolHeaderSize {
		return "", nil, ErrMessageTooShort
	}

	actionBytes := data[0:ActionHeaderSize]
	action = string(actionBytes)
	for i, c := range action {
		if c == ' ' {
			action = action[:i]
			break
		}
	}

	payloadSize := binary.BigEndian.Uint32(data[ActionHeaderSize:ProtocolHeaderSize])

	expectedSize := ProtocolHeaderSize + int(payloadSize)
	if len(data) != expectedSize {
		return "", nil, ErrMessageSizeMismatch
	}

	if payloadSize > 0 {
		payload = data[ProtocolHeaderSize:expectedSize]
	} else {
		payload = []byte{}
	}

	return action, payload, nil
}

// formatAction pads action name to 18 bytes with spaces
func formatAction(action string) []byte {
	padded := make([]byte, ActionHeaderSize)
	copy(padded, []byte(action))
	for i := len(action); i < ActionHeaderSize; i++ {
		padded[i] = ' '
	}
	return padded
}

// MustStructMessage creates a binary message, panics on error
func MustStructMessage(action string, payload []byte) []byte {
	msg, err := StructMessage(action, payload)
	if err != nil {
		panic(err)
	}
	return msg
}
