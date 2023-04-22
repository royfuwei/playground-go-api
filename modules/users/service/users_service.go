package service

import (
	"playground-go-api/domain"

	"golang.org/x/crypto/bcrypt"
)

type usersService struct {
}

func NewUsersService() domain.UsersService {
	return &usersService{}
}

// 檢查密碼是否正確
func (svc *usersService) CheckPassword(attemptPass string, user *domain.User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(attemptPass))
	if err != nil {
		return false
	}
	return true
}
