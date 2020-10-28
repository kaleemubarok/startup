package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUser) (User, error)
	IsEmailAvailable(input EmailCheckInput) (bool, error)
	SaveAvatar(ID int, filePathLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	email := input.Email
	isEmailAvailable, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if isEmailAvailable.ID != 0 {
		return user, errors.New("email " + email + " sudah digunakan")
	}

	passwordBcrypt, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordBcrypt)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginUser) (User, error) {
	email := input.Email
	password := input.Password

	loginUser, err := s.repository.FindByEmail(email)
	if err != nil {
		return loginUser, err
	}

	if loginUser.ID == 0 {
		return loginUser, errors.New("tidak ada user terdaftar dengan email" + email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginUser.PasswordHash), []byte(password))
	if err != nil {
		return loginUser, err
	}

	return loginUser, nil

}

func (s *service) IsEmailAvailable(input EmailCheckInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}

	return !strings.EqualFold(email, user.Email), nil

}

func (s *service) SaveAvatar(ID int, filePathLocation string) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID != 0 {
		user.AvatarFileName = filePathLocation
	}

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {
	user, err := s.repository.FindByID(ID)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("cannot find user with ID given")
	}

	return user, nil
}
