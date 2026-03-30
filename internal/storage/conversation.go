package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type ConversationStore struct {
	pg    *pgxpool.Pool
	redis *redis.Client
}

func NewConversationStore(pg *pgxpool.Pool, redis *redis.Client) *ConversationStore {
	return &ConversationStore{pg: pg, redis: redis}
}

func (s *ConversationStore) Create(ctx context.Context, userID, agentName string) (*Conversation, error) {
	conv := &Conversation{
		ID:        uuid.New().String(),
		UserID:    userID,
		AgentName: agentName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := s.pg.Exec(ctx,
		`INSERT INTO conversations (id, user_id, agent_name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		conv.ID, conv.UserID, conv.AgentName, conv.CreatedAt, conv.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create conversation: %w", err)
	}

	s.cacheConversation(ctx, conv)

	return conv, nil
}

func (s *ConversationStore) Get(ctx context.Context, id string) (*Conversation, error) {
	if s.redis != nil {
		cached, err := s.redis.Get(ctx, s.convKey(id)).Result()
		if err == nil {
			var conv Conversation
			if json.Unmarshal([]byte(cached), &conv) == nil {
				return &conv, nil
			}
		}
	}

	var conv Conversation
	err := s.pg.QueryRow(ctx,
		`SELECT id, user_id, agent_name, created_at, updated_at
		 FROM conversations WHERE id = $1`,
		id,
	).Scan(&conv.ID, &conv.UserID, &conv.AgentName, &conv.CreatedAt, &conv.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("conversation not found: %w", err)
	}

	s.cacheConversation(ctx, &conv)

	return &conv, nil
}

func (s *ConversationStore) AddMessage(ctx context.Context, convID, role, content string) (*Message, error) {
	msg := &Message{
		ID:             uuid.New().String(),
		ConversationID: convID,
		Role:           role,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	_, err := s.pg.Exec(ctx,
		`INSERT INTO messages (id, conversation_id, role, content, created_at)
		 VALUES ($1, $2, $3, $4, $5)`,
		msg.ID, msg.ConversationID, msg.Role, msg.Content, msg.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add message: %w", err)
	}

	s.invalidateConversationCache(ctx, convID)

	return msg, nil
}

func (s *ConversationStore) GetMessages(ctx context.Context, convID string) ([]Message, error) {
	rows, err := s.pg.Query(ctx,
		`SELECT id, conversation_id, role, content, created_at
		 FROM messages WHERE conversation_id = $1
		 ORDER BY created_at ASC`,
		convID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.Role, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

func (s *ConversationStore) convKey(id string) string {
	return fmt.Sprintf("conv:%s", id)
}

func (s *ConversationStore) cacheConversation(ctx context.Context, conv *Conversation) {
	if s.redis == nil {
		return
	}
	data, err := json.Marshal(conv)
	if err != nil {
		return
	}
	s.redis.Set(ctx, s.convKey(conv.ID), data, 30*time.Minute)
}

func (s *ConversationStore) invalidateConversationCache(ctx context.Context, convID string) {
	if s.redis == nil {
		return
	}
	s.redis.Del(ctx, s.convKey(convID))
}
