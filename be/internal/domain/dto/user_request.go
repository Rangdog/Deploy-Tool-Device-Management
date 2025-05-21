package dto

type UserRegisterRequest struct {
	Password    string `json:"password" binding:"required"`
	Email       string `json:"email" binding:"required"`
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResponse struct {
	Id        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	IsActive  bool   `json:"isActivate"`
}

type UserRequestResetPassword struct {
	NewPassword string `json:"newPassword"  binding:"required"`
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
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UpdateInformationUserRequest struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
}

type UpdateRoleUserRequest struct {
	RoleTitle string `json:"roleTitle" binding:"required"`
}
