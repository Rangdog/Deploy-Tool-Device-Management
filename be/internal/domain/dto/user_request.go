package dto

type UserRegisterRequest struct {
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	IsActive  bool   `json:"is_activate"`
}

type UserRequestResetPassword struct {
	NewPassword string `json:"new_password"  binding:"required"`
	Token       string `json:"token" binding:"required"`
}

type CheckPasswordReset struct {
	Email       string `json:"email" binding:"required"`
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type DeleteUserRequest struct {
	Email string `json:"email" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
