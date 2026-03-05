package services

import (
	"errors"

	"GrowEasy/config"
	"GrowEasy/dto"
	"GrowEasy/models"
	"GrowEasy/utils"

	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(req dto.RegisterRequest) error {

	var existing models.User
	err := config.DB.Where("email = ?", req.Email).First(&existing).Error
	if err == nil {
		return errors.New("email already registered")
	}
	if err != gorm.ErrRecordNotFound {
		return err
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: hashedPassword,
	}

	return config.DB.Create(&user).Error
}

func (s *AuthService) Login(req dto.LoginRequest) (*models.User, error) {

	var user models.User

	err := config.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	err = utils.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}
