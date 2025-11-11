package auth_service

import (
	"errors"
	auth_dto "my-project/modul/auth/dto"
	user_model "my-project/modul/user/model"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx echo.Context, req auth_dto.RegisterRequest) (auth_dto.AuthResponse, error)
	Login(ctx echo.Context, req auth_dto.LoginRequest) (auth_dto.AuthResponse, error)
	Refresh(ctx echo.Context) (auth_dto.AuthResponse, error)
	Logout(ctx echo.Context) error
	Me(ctx echo.Context) (auth_dto.UserResponse, error)
}

type authService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) AuthService {
	return &authService{db: db, jwtSecret: jwtSecret}
}

func (service *authService) Register(ctx echo.Context, req auth_dto.RegisterRequest) (auth_dto.AuthResponse, error) {

	var existing user_model.User
	{
		if err := service.db.Where("email = ?", req.Email).First(&existing).Error; err == nil {
			return auth_dto.AuthResponse{}, errors.New("email already registered")
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	{
		if err != nil {
			return auth_dto.AuthResponse{}, err
		}
	}

	user := user_model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := service.db.Create(&user).Error; err != nil {
		return auth_dto.AuthResponse{}, err
	}

	token, err := service.generateToken(user)
	{
		if err != nil {
			return auth_dto.AuthResponse{}, err
		}
	}

	return auth_dto.ToAuthResponse(user, token), nil
}

func (service *authService) Login(ctx echo.Context, req auth_dto.LoginRequest) (auth_dto.AuthResponse, error) {
	var user user_model.User
	{
		if err := service.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return auth_dto.AuthResponse{}, errors.New("invalid credentials")
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return auth_dto.AuthResponse{}, errors.New("invalid credentials")
	}

	token, err := service.generateToken(user)
	{
		if err != nil {
			return auth_dto.AuthResponse{}, err
		}
	}

	return auth_dto.ToAuthResponse(user, token), nil
}

func (service *authService) Refresh(ctx echo.Context) (auth_dto.AuthResponse, error) {
	user := ctx.Get("user").(*user_model.User)
	token, err := service.generateToken(*user)
	{
		if err != nil {
			return auth_dto.AuthResponse{}, err
		}
	}
	return auth_dto.ToAuthResponse(*user, token), nil
}

func (service *authService) Logout(ctx echo.Context) error {

	sess, _ := session.Get("session", ctx)
	sess.Values["token"] = ""
	sess.Values["user_id"] = ""
	sess.Values["user_email"] = ""
	sess.Values["user_name"] = ""
	sess.Save(ctx.Request(), ctx.Response())
	return nil
}

func (service *authService) Me(ctx echo.Context) (auth_dto.UserResponse, error) {

	user, ok := ctx.Get("user").(*user_model.User)
	{
		if !ok {
			return auth_dto.UserResponse{}, errors.New("unauthorized")
		}
	}

	return auth_dto.ToUserResponse(*user), nil
}

func (service *authService) generateToken(user user_model.User) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(service.jwtSecret))
}
