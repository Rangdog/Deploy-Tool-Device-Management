package handler

import (
	"BE_Manage_device/config"
	"BE_Manage_device/constant"
	"BE_Manage_device/internal/domain/dto"
	service "BE_Manage_device/internal/service/user"

	"BE_Manage_device/pkg"
	"BE_Manage_device/pkg/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// User godoc
// @Summary      Register user
// @Description  Register user
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserRegisterRequest   true  "Data"
// @Router       /api/auth/register [post]
// @Success      201   {object}  dto.ApiResponseSuccessNoData
// @Failure      500   {object}  dto.ApiResponseFail
func (h *UserHandler) Register(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}
	_, err := h.service.Register(user.FirstName, user.LastName, user.Password, user.Email, user.RedirectUrl)
	if err != nil {
		log.Error("Happened error when saving data to database. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, err.Error())
	}
	c.JSON(http.StatusCreated, pkg.BuildReponseSuccessNoData(http.StatusCreated, constant.Success))
}

// User godoc
// @Summary      Login
// @Description  Login
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user   body    dto.UserLoginRequest   true  "Data"
// @Router       /api/auth/login [post]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      500   {object}  dto.ApiResponseFail
func (h *UserHandler) Login(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var user dto.UserLoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE.")
	}

	userLogin, accessToken, refreshToken, err := h.service.Login(user.Email, user.Password)
	if err != nil {
		log.Error("Happened error when login. Error", err)
		pkg.PanicExeption(constant.Invalidemailorpassword)
	}
	dataResponese := map[string]interface{}{}
	if userLogin.IsActive {
		dataResponese = map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"is_active":     userLogin.IsActive,
			"roleSlug":      userLogin.Role.Slug,
		}
	} else {
		dataResponese = map[string]interface{}{
			"is_active": userLogin.IsActive,
		}
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, dataResponese))
}

func (h *UserHandler) Activate(c *gin.Context) {
	defer pkg.PanicHandler(c)
	token, exist := c.GetQuery("token")
	if !exist {
		log.Error("Happened error when mapping request from FE. Error: Don't see token in url")
		pkg.PanicExeption(constant.InvalidRequest, "Don't see token in url")
	}
	redirectUrl, exist := c.GetQuery("redirectUrl")
	if !exist {
		log.Error("Happened error when mapping request from FE. Error: Don't see token in url")
		pkg.PanicExeption(constant.InvalidRequest, "Error: Don't see token in url")
	}
	err := h.service.Activate(token)
	if err != nil {
		log.Error("Happened error when activate. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when activate")
	}
	c.Redirect(http.StatusFound, redirectUrl)
}

// User godoc
// @Summary      Refresh Token
// @Description  Refresh Token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        refresh_token   body    dto.RefreshRequest   true  "Data"
// @Router       /api/auth/refresh [POST]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      401   {object}  dto.ApiResponseFail
// @Failure      403   {object}  dto.ApiResponseFail
// @Failure      500   {object}  dto.ApiResponseFail
func (h *UserHandler) Refresh(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var rq dto.RefreshRequest
	if err := c.ShouldBindJSON(&rq); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE")
	}

	ok := h.service.CheckRefreshToken(rq.RefreshToken)
	if !ok {
		log.Error("Happened error refresh token was invoked")
		pkg.PanicExeption(constant.Unauthorized, "Refresh token was revoked")
	}
	refreshToken, err := jwt.Parse(rq.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.RefreshSecret), nil
	})
	if err != nil || !refreshToken.Valid {
		pkg.PanicExeption(constant.Unauthorized, "Refresh token was expired")
		return
	}

	if claims, ok := refreshToken.Claims.(jwt.MapClaims); ok && refreshToken.Valid {
		exp, ok := claims["exp"].(float64)
		if !ok {
			pkg.PanicExeption(constant.Unauthorized, "Invalid refresh token")
			return
		}
		if int64(exp) < time.Now().Unix() {
			pkg.PanicExeption(constant.StatusForbidden, "Refresh token was expired")
			return
		}
		email := claims["email"].(string)
		user, err := h.service.FindUserByEmail(email)
		newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": user.Id,
			"email":  email,
			"exp":    time.Now().Add(time.Minute * 1).Unix(),
		})

		if err != nil {
			pkg.PanicExeption(constant.UnknownError, "Happened error when create access token")
		}

		newAccessTokenString, err := newAccessToken.SignedString([]byte(config.AccessSecret))
		if err != nil {
			pkg.PanicExeption(constant.UnknownError, "Happened error when create access token")
		}
		data := map[string]interface{}{
			"access_token": newAccessTokenString,
		}
		c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, data))
	} else {
		pkg.PanicExeption(constant.Unauthorized, "invalid refresh token")
	}
}

// User godoc
// @Summary      Password-reset
// @Description  reset password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Password-reset   body    dto.UserRequestResetPassword   true  "Data"
// @Router       /api/user/password-reset [PATCH]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @Failure      401   {object}  dto.ApiResponseFail
// @Failure      500   {object}  dto.ApiResponseFail
func (h *UserHandler) ResetPassword(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.UserRequestResetPassword
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when mapping request from FE. Error")
	}
	token, err := jwt.Parse(request.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, http.ErrAbortHandler
		}
		return []byte(config.PasswordSecret), nil
	})
	if err != nil || !token.Valid {
		pkg.PanicExeption(constant.Unauthorized, "Token was expired")
		c.Abort()
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				pkg.PanicExeption(constant.Unauthorized, "Token expired")
			}
		}
		if email, ok := claims["email"].(string); ok {
			user, err := h.service.FindUserByEmail(email)
			if err != nil {
				log.Error("Happened error when email don't exist. Error", err)
				pkg.PanicExeption(constant.UnknownError, "Email don't exist")
			}
			err = h.service.ResetPassword(user, request.NewPassword)
			if err != nil {
				log.Error("Happened error when resert password. Error", err)
				pkg.PanicExeption(constant.UnknownError, err.Error())
			}
			c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
		} else {
			pkg.PanicExeption(constant.Unauthorized, "Invalid token")
		}
	} else {
		pkg.PanicExeption(constant.Unauthorized, "Invalid token")
	}
}

// User godoc
// @Summary      Get session
// @Description  Get session
// @Tags         Users
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/user/session [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      401   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) Session(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	user, err := h.service.FindByUserId(userId)
	if err != nil {
		log.Error("Happened error when reset password. Error", err)
		pkg.PanicExeption(constant.Unauthorized, err.Error())
	}
	userResponse := utils.ConvertUserToUserResponse(user)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, userResponse))
}

// User godoc
// @Summary      Email reset password
// @Description   Email reset password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Email_Reset_Password   body    dto.CheckPasswordReset   true  "Data"
// @Router       /api/user/forget-password [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @Failure      500   {object}  dto.ApiResponseFail
func (h *UserHandler) CheckPasswordReset(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.CheckPasswordReset
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err := h.service.CheckPasswordReset(request.Email, request.RedirectUrl)
	if err != nil {
		log.Error("Happened error when reset password. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Delete user
// @Description   Delete user via email
// @Tags         Users
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Param		email	path		string				true	"email"
// @Router       /api/user/{email} [DELETE]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @Failure      500   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) DeleteUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	email := c.Param("email")
	err := h.service.DeleteUser(email)
	if err != nil {
		log.Error("Happened error when delete user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Logout
// @Description   Logout
// @Tags         Auth
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/auth/logout [POST]
// @Success      200   {object}  dto.ApiResponseSuccessNoData
// @Failure      500   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) Logout(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	userSession, err := h.service.FindSessionById(userId)
	if err != nil {
		log.Error("Happened error when logout user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	err = h.service.UpdateInvoked(*userSession)
	if err != nil {
		log.Error("Happened error when logout user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Get all user
// @Description   Get all user
// @Tags         Users
// @Accept       json
// @Produce      json
// @param Authorization header string true "Authorization"
// @Router       /api/users [GET]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) GetAllUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	// userId := utils.GetUserIdFromContext(c)
	users := h.service.GetAllUser()
	usersResponses := utils.ConvertUsersToUserResponses(users)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, usersResponses))
}

// User godoc
// @Summary      Update Information
// @Description   Update Information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        Information   body    dto.UpdateInformationUserRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/user/information [PATCH]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      500   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateInformationUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.UpdateInformationUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	userUpdated, err := h.service.UpdateInformation(userId, request.FirstName, request.LastName)
	if err != nil {
		log.Error("Happened error when update information user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	usersResponse := utils.ConvertUserToUserResponse(userUpdated)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, usersResponse))
}

// Role godoc
// @Summary      Update role by id
// @Description   Update role by id
// @Tags         Roles
// @Accept       json
// @Produce      json
// @Param        Role   body    dto.UpdateRoleUserRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/users/role [PATCH]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      500   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateRoleUser(c *gin.Context) {
	defer pkg.PanicHandler(c)
	userId := utils.GetUserIdFromContext(c)
	var request dto.UpdateRoleUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	userUpdated, err := h.service.UpdateRole(userId, request.Id, request.Slug)
	if err != nil {
		log.Error("Happened error when update role user. Error", err)
		pkg.PanicExeption(constant.UnknownError, err.Error())
	}
	usersResponse := utils.ConvertUserToUserResponse(userUpdated)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, usersResponse))
}

// Role godoc
// @Summary      Get all user by department_id
// @Description  Get all user by department_id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		department_id	path		string				true	"department_id"
// @param Authorization header string true "Authorization"
// @Router       /api/user/department/{department_id} [GET]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) GetAllUserOfDepartment(c *gin.Context) {
	defer pkg.PanicHandler(c)
	idStr := c.Param("department_id")
	departmentId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert assetId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert assetId to int64")
	}
	users, err := h.service.GetAllUserOfDepartment(departmentId)
	if err != nil {
		log.Error("Happened error when get all user of department. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all user of department")
	}
	userResponses := utils.ConvertUsersToUserResponses(users)
	c.JSON(http.StatusOK, pkg.BuildReponseSuccess(http.StatusOK, constant.Success, userResponses))
}

// Role godoc
// @Summary      Update department by id
// @Description   Update department by id
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        department   body    dto.UserUpdateDepartmentRequest   true  "Data"
// @param Authorization header string true "Authorization"
// @Router       /api/user/department [PATCH]
// @Success      200   {object}  dto.ApiResponseSuccessStruct
// @Failure      500   {object}  dto.ApiResponseFail
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateDepartment(c *gin.Context) {
	defer pkg.PanicHandler(c)
	var request dto.UserUpdateDepartmentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Error("Happened error when mapping request from FE. Error", err)
		pkg.PanicExeption(constant.InvalidRequest)
	}
	err := h.service.UpdateDepartment(request.UserId, request.DepartmentId)
	if err != nil {
		log.Error("Happened error when get all user of department. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when get all user of department")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Update head department by userId
// @Description  Update head department by userId
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		user_id	path		string				true	"user_id"
// @param Authorization header string true "Authorization"
// @Router       /api/user/head-department/{user_id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateHeadDep(c *gin.Context) {
	defer pkg.PanicHandler(c)
	idStr := c.Param("user_id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert userId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert userId to int64")
	}
	err = h.service.UpdateHeadDep(userId)
	if err != nil {
		log.Error("Happened error when update user is head department. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when update user is head department")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Update head department by userId
// @Description  Update head department by userId
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		user_id	path		string				true	"user_id"
// @param Authorization header string true "Authorization"
// @Router       /api/user/manager-department/{user_id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateManagerDep(c *gin.Context) {
	defer pkg.PanicHandler(c)
	idStr := c.Param("user_id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert userId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert userId to int64")
	}
	err = h.service.UpdateManagerDep(userId)
	if err != nil {
		log.Error("Happened error when update user is manager department. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when update user is manager department")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}

// User godoc
// @Summary      Update can-export by userId
// @Description  Update can-export by userId
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		user_id	path		string				true	"user_id"
// @param Authorization header string true "Authorization"
// @Router       /api/user/can-export/{user_id} [PATCH]
// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
func (h *UserHandler) UpdateCanExport(c *gin.Context) {
	defer pkg.PanicHandler(c)
	idStr := c.Param("user_id")
	userId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Error("Happened error when convert userId to int64. Error", err)
		pkg.PanicExeption(constant.InvalidRequest, "Happened error when convert userId to int64")
	}
	err = h.service.UpdateCanExport(userId)
	if err != nil {
		log.Error("Happened error when update user can export. Error", err)
		pkg.PanicExeption(constant.UnknownError, "Happened error when update user can export")
	}
	c.JSON(http.StatusOK, pkg.BuildReponseSuccessNoData(http.StatusOK, constant.Success))
}
