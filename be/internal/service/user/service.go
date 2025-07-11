package service

import (
	"BE_Manage_device/config"
	"BE_Manage_device/internal/domain/entity"
	asset "BE_Manage_device/internal/repository/assets"
	company "BE_Manage_device/internal/repository/company"
	role "BE_Manage_device/internal/repository/role"
	user "BE_Manage_device/internal/repository/user"
	userRBAC "BE_Manage_device/internal/repository/user_rbac"
	userSession "BE_Manage_device/internal/repository/user_session"
	emailS "BE_Manage_device/internal/service/email"
	"BE_Manage_device/pkg/utils"

	"errors"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo               user.UserRepository
	emailService       *emailS.EmailService
	userSessionRepo    userSession.UsersSessionRepository
	roleRepository     role.RoleRepository
	assetRepo          asset.AssetsRepository
	userRBACRepository userRBAC.UserRBACRepository
	CompanyRepo        company.CompanyRepository
}

func NewUserService(repo user.UserRepository, emailService *emailS.EmailService, userSessionRepo userSession.UsersSessionRepository, roleRepository role.RoleRepository, assetRepo asset.AssetsRepository, userRBACRepository userRBAC.UserRBACRepository, CompanyRepo company.CompanyRepository) *UserService {
	return &UserService{repo: repo, emailService: emailService, userSessionRepo: userSessionRepo, roleRepository: roleRepository, assetRepo: assetRepo, userRBACRepository: userRBACRepository, CompanyRepo: CompanyRepo}
}

func (service *UserService) Register(firstName, lastName, password, email, redirectUrl string) (*entity.Users, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role := service.roleRepository.GetRoleBySlug("employee")
	token := uuid.New().String()
	company, err := service.CompanyRepo.GetCompanyBySuffixEmail(utils.GetSuffixEmail(email))
	if err != nil {
		return nil, errors.New("this email company don't register")
	}
	users := &entity.Users{
		FirstName: firstName,
		LastName:  lastName,
		Password:  string(hashedPassword),
		Email:     email,
		RoleId:    role.Id,
		IsActive:  false,
		Token:     token,
		CompanyId: company.Id,
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
	var err error
	user, err := service.repo.FindByUserId(userId)
	if err != nil {
		return
	}

	assets, _ := service.assetRepo.GetAllAsset(user.CompanyId)
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

func (service *UserService) GetAllUser(userId int64) []*entity.Users {
	user, err := service.repo.FindByUserId(userId)
	if err != nil {
		return nil
	}
	users := service.repo.GetAllUser(user.CompanyId)
	return users
}

func (service *UserService) UpdateInformation(id int64, firstName, lastName string, image *multipart.FileHeader) (*entity.Users, error) {
	imgFile, err := image.Open()
	if err != nil {
		return nil, fmt.Errorf("cannot open image: %w", err)
	}
	defer imgFile.Close()
	uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), image.Filename)
	imagePath := "avatar/" + uniqueName
	uploader := utils.NewSupabaseUploader()
	imageUrl, err := uploader.Upload(imagePath, imgFile, image.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	user := entity.Users{Id: id, FirstName: firstName, LastName: lastName, Avatar: imageUrl}
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
	if slug == "assetManager" {
		service.UpdateManagerDep(setRoleUserId)
	}
	if slug == "employee" {
		if userUpdated.IsAssetManager {
			service.UpdateManagerDep(setRoleUserId)
		}
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
		err = service.repo.UpdateManagerDep(id, !user.IsAssetManager)
		return err
	}
}

func (service *UserService) UpdateCanExport(id int64) error {
	user, err := service.repo.FindByUserId(id)
	if err != nil {
		return err
	}
	roleEmployee := service.roleRepository.GetRoleBySlug("employee")
	if user.RoleId != roleEmployee.Id {
		return nil
	}
	err = service.repo.UpdateCanExport(id, !user.CanExport)
	return err
}

func (service *UserService) GetUserNotHaveDep() ([]*entity.Users, error) {
	users, err := service.repo.GetUserNotHaveDep()

	if err != nil {
		return nil, err
	}
	return users, nil
}
