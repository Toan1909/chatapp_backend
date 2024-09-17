package model

import "time"

type Conversation struct {
	ConversationId   string `json:"conversationId,omitempty" db:"conversation_id, omitempty"`
	ConversationName string	`json:"conversationName,omitempty" db:"conversation_name, omitempty"`
	IsGroup          bool	`json:"isGroup" db:"is_group, omitempty"`
	CreatedAt        time.Time`json:"createdAt,omitempty" db:"created_at, omitempty"`
	UpdatedAt 		 time.Time`json:"updatedAt,omitempty" db:"updated_at, omitempty"`
}