package auth_service

import (
	"errors"
	"my-project/helper"
	user_model "my-project/modul/user/model"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(reqRegister *RegisterDto) (*user_model.User, error)
	VerifyEmail(token string) error
	Login(email, password string) (string, error)
	FindById(id uint) (*user_model.User, error)
	UpdateProfile(id uint, name string) (*user_model.User, error)
	RequestPasswordReset(email string) error
	ResetPassword(token, newPassword string) error
}

type authService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthService {
	return &authService{db: db}
}

type RegisterDto struct {
	Name     string
	Email    string
	Password string
}

func (s *authService) Register(req *RegisterDto) (*user_model.User, error) {
	// check email exists
	var u user_model.User
	if err := s.db.Where("email = ?", req.Email).First(&u).Error; err == nil {
		return nil, errors.New("email already registered")
	}
	passHash, _ := helper.HashPassword(req.Password)
	vt, _ := helper.RandToken(16)

	user := user_model.User{
		Name:              req.Name,
		Email:             req.Email,
		PasswordHash:      passHash,
		IsVerified:        false,
		VerificationToken: vt,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// send verification email (async in prod, here immediate)
	verifyLink := helper.ENV("APP_URL") + "/api/v1/admin/auth/verify?token=" + vt
	body := "Please verify your email: " + verifyLink
	_ = helper.SendEmail(user.Email, "Verify your email", body)

	return &user, nil
}

func (s *authService) VerifyEmail(token string) error {
	var user user_model.User
	if err := s.db.Where("verification_token = ?", token).First(&user).Error; err != nil {
		return errors.New("invalid token")
	}
	user.IsVerified = true
	user.VerificationToken = ""
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *authService) Login(email, password string) (string, error) {
	var user user_model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return "", errors.New("invalid credentials")
	}
	if !helper.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	if !user.IsVerified {
		return "", errors.New("email not verified")
	}

	token, err := helper.JwtGenerate(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) FindById(id uint) (*user_model.User, error) {
	var user user_model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *authService) UpdateProfile(id uint, name string) (*user_model.User, error) {
	user, err := s.FindById(id)
	if err != nil {
		return nil, err
	}
	user.Name = name
	if err := s.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) RequestPasswordReset(email string) error {
	var user user_model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		// do not reveal email existence — silently return nil
		return nil
	}
	token, _ := helper.RandToken(16)
	exp := time.Now().Add(1 * time.Hour)
	user.PasswordResetToken = token
	user.PasswordResetExpiry = &exp
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	resetLink := helper.ENV("APP_URL") + "/api/v1/auth/password-reset?token=" + token
	body := "Reset your password: " + resetLink
	_ = helper.SendEmail(user.Email, "Password reset", body)
	return nil
}

func (s *authService) ResetPassword(token, newPassword string) error {
	var user user_model.User
	if err := s.db.Where("password_reset_token = ?", token).First(&user).Error; err != nil {
		return errors.New("invalid token")
	}
	if user.PasswordResetExpiry == nil || user.PasswordResetExpiry.Before(time.Now()) {
		return errors.New("token expired")
	}
	passHash, _ := helper.HashPassword(newPassword)
	user.PasswordHash = passHash
	user.PasswordResetToken = ""
	user.PasswordResetExpiry = nil
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}
	return nil
}
