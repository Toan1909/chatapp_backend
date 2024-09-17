package req
type ReqSignUp struct{
	FullName string `json:"fullName,omitempty" validate:"required"`
	Email string	`json:"email,omitempty" validate:"required,email"`
	Phone string	`json:"phone,omitempty" validate:"required"`
	Password string	`json:"password,omitempty" validate:"required"`
}
type ReqSignIn struct{
	Email string	`json:"email,omitempty" validate:"required,email"`
	Password string	`json:"password,omitempty" validate:"required"`
}