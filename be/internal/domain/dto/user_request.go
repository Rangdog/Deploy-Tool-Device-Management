package dto

type UserRegisterRequest struct {
	Password    string `json:"password" binding:"required,min=6"`
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"firstName" binding:"required,min=2,max=30"`
	LastName    string `json:"lastName" binding:"required,min=2,max=30"`
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
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
	Email       string `json:"email" binding:"required,email"`
	RedirectUrl string `json:"redirectUrl" binding:"required"`
}

type DeleteUserRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type UpdateInformationUserRequest struct {
	FirstName string `json:"firstName" binding:"required,min=2,max=30"`
	LastName  string `json:"lastName" binding:"required,min=2,max=30"`
}

type UpdateRoleUserRequest struct {
	Id   int64  `json:"userId" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}

type UserResponse struct {
	Id         int64                   `json:"id"`
	FirstName  string                  `json:"firstName"`
	LastName   string                  `json:"lastName"`
	Email      string                  `json:"email"`
	IsActive   bool                    `json:"isActivate"`
	Role       UserRoleResponse        `json:"role"`
	Department *UserDepartmentResponse `json:"department,omitempty"`
}

type UserRoleResponse struct {
	Id   int64  `json:"id"`
	Slug string `json:"slug"`
}

type UserDepartmentResponse struct {
	Id             int64  `json:"id"`
	DepartmentName string `json:"departmentName"`
}

type UserUpdateDepartmentRequest struct {
	UserId       int64 `json:"userId"`
	DepartmentId int64 `json:"departmentId"`
}
