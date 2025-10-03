package user_dto

type RequestPasswordReset struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPassword struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
