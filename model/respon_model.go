package model

type Response struct {
	StatusCode int         `json:"statusCode,omitempty"`
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}
type ResponseWs struct {
	Clients []ConversationMember      `json:"-"`
	Type    string      `json:"type,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
