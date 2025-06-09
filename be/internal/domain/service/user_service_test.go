package service

import (
	"BE_Manage_device/internal/domain/entity"
	"BE_Manage_device/internal/domain/repository"
	"BE_Manage_device/internal/domain/repository/mocks"
	"BE_Manage_device/pkg/utils"
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	tests := []struct {
		name string
		args args
		want *UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.repo, tt.args.emailService, tt.args.userSessionRepo, tt.args.roleRepository, tt.args.assetRepo, tt.args.userRBACRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Register(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		firstName   string
		lastName    string
		password    string
		email       string
		redirectUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.Register(tt.args.firstName, tt.args.lastName, tt.args.password, tt.args.email, tt.args.redirectUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		want1   string
		want2   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, got1, got2, err := service.Login(tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.Login() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserService.Login() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("UserService.Login() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}

func TestUserService_Activate(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.Activate(tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("UserService.Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_SetRole(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userId int64
		roleId int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			service.SetRole(tt.args.userId, tt.args.roleId)
		})
	}
}

func TestUserService_FindUserByEmail(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.FindUserByEmail(tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.FindUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.FindUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_ResetPassword(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		user        *entity.Users
		newPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.ResetPassword(tt.args.user, tt.args.newPassword); (err != nil) != tt.wantErr {
				t.Errorf("UserService.ResetPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_CheckPasswordReset(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		email       string
		redirectUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.CheckPasswordReset(tt.args.email, tt.args.redirectUrl); (err != nil) != tt.wantErr {
				t.Errorf("UserService.CheckPasswordReset() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.DeleteUser(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("UserService.DeleteUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_CheckRefreshToken(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		token string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if got := service.CheckRefreshToken(tt.args.token); got != tt.want {
				t.Errorf("UserService.CheckRefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_FindByUserId(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.FindByUserId(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.FindByUserId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.FindByUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_FindSessionById(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.UsersSessions
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.FindSessionById(tt.args.userId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.FindSessionById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.FindSessionById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateInvoked(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userSession entity.UsersSessions
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.UpdateInvoked(tt.args.userSession); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateInvoked() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_GetAllUser(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	tests := []struct {
		name   string
		fields fields
		want   []*entity.Users
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if got := service.GetAllUser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetAllUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateInformation(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		id        int64
		firstName string
		lastName  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.UpdateInformation(tt.args.id, tt.args.firstName, tt.args.lastName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateInformation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.UpdateInformation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateRole(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userId        int64
		setRoleUserId int64
		slug          string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.UpdateRole(tt.args.userId, tt.args.setRoleUserId, tt.args.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.UpdateRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_GetAllUserOfDepartment(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		departmentId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entity.Users
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			got, err := service.GetAllUserOfDepartment(tt.args.departmentId)
			if (err != nil) != tt.wantErr {
				t.Errorf("UserService.GetAllUserOfDepartment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserService.GetAllUserOfDepartment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserService_UpdateDepartment(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		userId       int64
		departmentId int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.UpdateDepartment(tt.args.userId, tt.args.departmentId); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateDepartment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateHeadDep(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.UpdateHeadDep(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateHeadDep() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateManagerDep(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update manager dep",
			fields: fields{
				repo: func() repository.UserRepository {
					m := new(mocks.UserRepository)
					m.On("UpdateCanExport", int64(1)).Return(nil)
					m.On("FindByUserId", int64(1)).Return(&entity.Users{Id: 1,
						IsAssetManager: false,
						DepartmentId:   utils.PtrInt64(2)}, nil)
					m.On("UpdateManagerDep", int64(1), true).Return(nil)
					m.On("CheckManagerDep", int64(2)).Return(nil)
					return m
				}(),
			},
			args:    args{id: 1},
			wantErr: false,
		},
		{
			name: "Update user don't have department",
			fields: fields{
				repo: func() repository.UserRepository {
					m := new(mocks.UserRepository)
					m.On("UpdateCanExport", int64(1)).Return(nil)
					m.On("FindByUserId", int64(1)).Return(&entity.Users{Id: 1,
						IsAssetManager: false,
					}, nil)
					m.On("UpdateManagerDep", int64(1), true).Return(nil)
					m.On("CheckManagerDep", int64(2)).Return(errors.New("Update user don't have department"))
					return m
				}(),
			},
			args:    args{id: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.UpdateManagerDep(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateManagerDep() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_UpdateCanExport(t *testing.T) {
	type fields struct {
		repo               repository.UserRepository
		emailService       *EmailService
		userSessionRepo    repository.UsersSessionRepository
		roleRepository     repository.RoleRepository
		assetRepo          repository.AssetsRepository
		userRBACRepository repository.UserRBACRepository
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Update thành công",
			fields: fields{
				repo: func() repository.UserRepository {
					m := new(mocks.UserRepository)
					m.On("UpdateCanExport", int64(1)).Return(nil)
					m.On("FindByUserId", int64(1)).Return(&entity.Users{ /*...*/ }, nil)
					return m
				}(),
				roleRepository: func() repository.RoleRepository {
					m := new(mocks.RoleRepository)
					m.On("GetRoleBySlug", "viewer").Return(&entity.Roles{Id: 4, Slug: "viewer"}, nil)
					return m
				}(),
			},
			args:    args{id: 1},
			wantErr: false,
		},
		{
			name: "Update lỗi",
			fields: fields{
				repo: func() repository.UserRepository {
					m := new(mocks.UserRepository)
					m.On("UpdateCanExport", int64(2), mock.Anything).Return(errors.New("db error"))
					m.On("FindByUserId", mock.Anything).Return(&entity.Users{
						RoleId:    4,
						CanExport: false,
					}, nil)
					return m
				}(),
				roleRepository: func() repository.RoleRepository {
					m := new(mocks.RoleRepository)
					m.On("GetRoleBySlug", "viewer").Return(&entity.Roles{Id: 4, Slug: "viewer"}, nil)
					return m
				}(),
			},
			args:    args{id: 2},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &UserService{
				repo:               tt.fields.repo,
				emailService:       tt.fields.emailService,
				userSessionRepo:    tt.fields.userSessionRepo,
				roleRepository:     tt.fields.roleRepository,
				assetRepo:          tt.fields.assetRepo,
				userRBACRepository: tt.fields.userRBACRepository,
			}
			if err := service.UpdateCanExport(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("UserService.UpdateCanExport() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
