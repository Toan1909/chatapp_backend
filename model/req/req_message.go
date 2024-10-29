package req
type ReqGetMessage struct {
	ConversId string `json:"conversId,omitempty" validate:"required"`
}

type SendMessage struct {
	ConversationId string    `json:"conversationId,omitempty" db:"conversation_id, omitempty"`
	SenderId       string    `json:"senderId,omitempty" db:"sender_id, omitempty"`
	Content        string    `json:"content,omitempty" db:"content, omitempty"`
	MediaUrl       string    `json:"mediaUrl,omitempty" db:"media_url, omitempty"`
}
type ReqReadReceipt struct {
	MessageId      string    `json:"messageId,omitempty" db:"message_id, omitempty"`
	ConversationId string	`json:"conversationId,omitempty" db:"conversation_id, omitempty"`
}