package service

import (
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo            repository.UserRepository
	emailService    *EmailService
	userSessionRepo repository.UsersSessionRepository
}

func NewUserService(repo repository.UserRepository, emailService *EmailService, userSessionRepo repository.UsersSessionRepository) *UserService {
	return &UserService{repo: repo, emailService: emailService, userSessionRepo: userSessionRepo}
}

func (service *UserService) Register(firstName, lastName, password, email, redirectUrl string) (*entity.Users, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	token := uuid.New().String()
	users := &entity.Users{
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
		Email:     email,
		RoleId:    1,
		IsActive:  false,
		Token:     token,
	}
	err = service.repo.Create(users)
	if err != nil {
		return nil, err
	}
	go service.emailService.SendActivationEmail(email, token, redirectUrl)
	return users, nil
}

func (service *UserService) Login(email string, password string) (*entity.Users, string, string, error) {
	user, err := service.repo.FindByEmail(email)
	if err != nil {
		return nil, "", "", errors.New("email have")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.Id, email)
	if err != nil {
		return nil, "", "", err
	}
	if service.userSessionRepo.CheckUserInSession(user.Id) {
		userSession, err := service.userSessionRepo.FindByUserIdInSession(user.Id)
		if err != nil {
			return nil, "", "", err
		}
		err = service.userSessionRepo.UpdateIsRevoked(userSession)
		if err != nil {
			return nil, "", "", err
		}
	}
	userSession := entity.UsersSessions{
		UserId:       user.Id,
		CreatedAt:    time.Now(),
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
		ExpiresAt:    time.Now().Add(5 * time.Minute),
	}
	tx := service.repo.GetDB().Begin()
	err = service.userSessionRepo.Create(&userSession, tx)
	if err != nil {
		tx.Rollback()
		return nil, "", "", err
	}
	tx.Commit()
	return user, accessToken, refreshToken, nil
}

func (service *UserService) Activate(token string) error {
	users, err := service.repo.FindByToken(token)
	if err != nil {
		return err
	}
	err = service.repo.Update(users)
	return err
}

func (service *UserService) FindUserByEmail(email string) (*entity.Users, error) {
	user, err := service.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (service *UserService) ResetPassword(user *entity.Users, newPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword)); err == nil {
		return errors.New("the new password must not be the same as the old password")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = service.repo.UpdatePassword(user)
	return err
}

func (service *UserService) CheckPasswordReset(email string, redirectUrl string) error {

	tokenPW := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(2 * time.Minute).Unix(),
	})
	tokenPWstring, err := tokenPW.SignedString([]byte(config.PasswordSecret))
	if err != nil {
		return err
	}
	_, err = service.repo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("email don't register")
		}
		return err
	}
	body := "Click link to reset password account: <a href='" + redirectUrl + "?token=" + tokenPWstring + "'>reset</a>"
	go service.emailService.SendEmail(email, "Reset Password", body)
	return err
}

func (service *UserService) DeleteUser(email string) error {
	err := service.repo.DeleteUser(email)
	return err
}

func (service *UserService) CheckRefreshToken(token string) bool {
	userSession, err := service.userSessionRepo.FindByRefreshToken(token)
	if err != nil {
		return false
	}
	if userSession.IsRevoked {
		return false
	}
	return true
}

func (service *UserService) FindByUserId(userId int64) (*entity.Users, error) {
	user, err := service.repo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (service *UserService) FindSessionById(userId int64) (*entity.UsersSessions, error) {
	session, err := service.userSessionRepo.FindByUserIdInSession(userId)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (service *UserService) UpdateInvoked(userSession entity.UsersSessions) error {
	err := service.userSessionRepo.UpdateIsRevoked(&userSession)
	return err
}
