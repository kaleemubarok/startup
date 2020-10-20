package user

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginUser) (User, error)
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
		return loginUser, errors.New("Tidak ada user terdaftar dengan email" + email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(loginUser.PasswordHash), []byte(password))
	if err != nil {
		return loginUser, err
	}

	return loginUser, nil

}
