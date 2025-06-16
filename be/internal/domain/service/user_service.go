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
	repo               repository.UserRepository
	emailService       *EmailService
	userSessionRepo    repository.UsersSessionRepository
	roleRepository     repository.RoleRepository
	assetRepo          repository.AssetsRepository
	userRBACRepository repository.UserRBACRepository
}

func NewUserService(repo repository.UserRepository, emailService *EmailService, userSessionRepo repository.UsersSessionRepository, roleRepository repository.RoleRepository, assetRepo repository.AssetsRepository, userRBACRepository repository.UserRBACRepository) *UserService {
	return &UserService{repo: repo, emailService: emailService, userSessionRepo: userSessionRepo, roleRepository: roleRepository, assetRepo: assetRepo, userRBACRepository: userRBACRepository}
}

func (service *UserService) Register(firstName, lastName, password, email, redirectUrl string) (*entity.Users, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role := service.roleRepository.GetRoleBySlug("viewer")
	token := uuid.New().String()
	users := &entity.Users{
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
		Email:     email,
		RoleId:    role.Id,
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
		return nil, "", "", errors.New("email dont; have")
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
	go service.SetRole(users.Id, users.RoleId)
	return err
}

func (service *UserService) SetRole(userId int64, roleId int64) {
	tx := service.repo.GetDB().Begin()
	assets, _ := service.assetRepo.GetAllAsset()
	for _, asset := range assets {
		userRbac := entity.UserRbac{
			AssetId: asset.Id,
			UserId:  userId,
			RoleId:  roleId,
		}
		err := service.userRBACRepository.Create(&userRbac, tx)
		if err != nil {
			tx.Rollback()
		}
	}
	var err error
	if err = tx.Commit().Error; err != nil {
		return
	}
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
		"exp":   time.Now().Add(10 * time.Minute).Unix(),
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

func (service *UserService) GetAllUser() []*entity.Users {
	users := service.repo.GetAllUser()
	return users
}

func (service *UserService) UpdateInformation(id int64, firstName, lastName string) (*entity.Users, error) {
	user := entity.Users{Id: id, FirstName: firstName, LastName: lastName}
	userUpdated, err := service.repo.UpdateUser(&user)
	if err != nil {
		return nil, err
	}
	return userUpdated, nil
}

func (service *UserService) UpdateRole(userId int64, setRoleUserId int64, slug string) (*entity.Users, error) {
	roles := service.roleRepository.GetRoleBySlug(slug)
	if roles == nil {
		return nil, errors.New("something went wrong ")
	}
	user := entity.Users{Id: setRoleUserId, RoleId: roles.Id}
	userUpdated, err := service.repo.UpdateUser(&user)
	if err != nil {
		return nil, err
	}
	return userUpdated, nil
}

func (service *UserService) GetAllUserOfDepartment(departmentId int64) ([]*entity.Users, error) {
	user, err := service.repo.GetAllUserOfDepartment(departmentId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *UserService) UpdateDepartment(userId int64, departmentId int64) error {
	err := service.repo.UpdateDepartment(userId, departmentId)
	if err != nil {
		return err
	}
	go service.assetRepo.DeleteOwnerAssetOfOwnerId(userId)
	return err
}

func (service *UserService) UpdateHeadDep(id int64) error {
	user, err := service.repo.FindByUserId(id)
	if err != nil {
		return err
	}
	if user.IsHeadDepartment {
		err := service.repo.UpdateHeadDep(id, !user.IsHeadDepartment)
		return err
	} else {
		if user.DepartmentId == nil {
			return errors.New("user don't have department")
		}
		err := service.repo.CheckHeadDep(*user.DepartmentId)
		if err != nil {
			return err
		}
		err = service.repo.UpdateHeadDep(id, !user.IsHeadDepartment)
		return err
	}
}

func (service *UserService) UpdateManagerDep(id int64) error {
	user, err := service.repo.FindByUserId(id)
	if err != nil {
		return err
	}
	if user.IsAssetManager {
		err := service.repo.UpdateManagerDep(id, !user.IsAssetManager)
		return err
	} else {
		if user.DepartmentId == nil {
			return errors.New("user don't have department")
		}
		err := service.repo.CheckManagerDep(*user.DepartmentId)
		if err != nil {
			return err
		}
		err = service.repo.UpdateManagerDep(id, !user.IsAssetManager)
		return err
	}
}

func (service *UserService) UpdateCanExport(id int64) error {
	user, err := service.repo.FindByUserId(id)
	if err != nil {
		return err
	}
	roleViewer := service.roleRepository.GetRoleBySlug("viewer")
	if user.RoleId != roleViewer.Id {
		return nil
	}
	err = service.repo.UpdateCanExport(id, !user.CanExport)
	return err
}
