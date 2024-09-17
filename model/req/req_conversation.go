package req

type ReqCreateConvers struct {
	ListMember []string `json:"listMember,omitempty" validate:"required"`
	ConversationName  string `json:"conversationName,omitempty" validate:"required"`
}
type ReqLoadMem struct {
	ConversationId string    `json:"conversationId,omitempty" db:"conversation_id, omitempty"`
}