package model

import (
	"time"
	"github.com/lib/pq"
)

type Message struct {
	MessageId      string    `json:"messageId,omitempty" db:"message_id, omitempty"`
	ConversationId string    `json:"conversationId,omitempty" db:"conversation_id, omitempty"`
	SenderId       string    `json:"senderId,omitempty" db:"sender_id, omitempty"`
	Content        string    `json:"content,omitempty" db:"content, omitempty"`
	MediaUrl       string    `json:"mediaUrl,omitempty" db:"media_url, omitempty"`
	SendAt         time.Time `json:"sendAt,omitempty" db:"sent_at, omitempty"`
	SeenBy pq.StringArray `json:"seenBy,omitempty" db:"seen_by, omitempty"`
}