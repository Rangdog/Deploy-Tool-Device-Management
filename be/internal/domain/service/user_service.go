package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/pkg/utils"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         repository.UserRepository
	emailService *EmailService
}

func NewUserService(repo repository.UserRepository, emailService *EmailService) *UserService {
	return &UserService{repo: repo, emailService: emailService}
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
		RoleId:    0,
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

func (service *UserService) ResetPassword(userId int64, newPassword, oldPassword string) error {
	if newPassword == oldPassword {
		return errors.New("new password equal old password")
	}
	user, err := service.repo.FindByUserId(userId)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	err = service.repo.UpdatePassword(user)
	return err
}
