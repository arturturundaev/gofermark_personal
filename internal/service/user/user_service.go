package user

import (
	"fmt"
	"github.com/google/uuid"
	"gofermark_personal/internal/model"
	"gofermark_personal/internal/service"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository service.IUserRepository
}

func NewUserService(repository service.IUserRepository) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) Register(login string, password string) (*uuid.UUID, error) {
	id := uuid.New()
	passwordBytes, err := service.hashPassword(password)

	if err != nil {
		return nil, err
	}
	password = string(passwordBytes)

	return &id, service.repository.Save(id, login, password)
}

func (service *UserService) UserExists(login string) (bool, error) {
	return service.repository.UserExistsByLogin(login)
}

func (service *UserService) Auth(login string, password string) (*model.User, error) {
	user, err := service.repository.GetByLogin(login)

	if err != nil {
		return nil, err
	}

	if user != nil {
		err = service.compareHashAndPassword([]byte(user.Password), password)

		if err != nil {
			return nil, err
		}

		return user, nil
	}

	return nil, fmt.Errorf("something error")
}

func (service *UserService) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (service *UserService) compareHashAndPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}

func (service *UserService) GetBalance(userID uuid.UUID) (*model.UserBalance, error) {
	return service.repository.GetBalance(userID)
}

func (service *UserService) Withdraw(userID uuid.UUID, number string, sum float64) error {
	return service.repository.Withdraw(userID, number, sum)
}

func (service *UserService) GetWithdrawals(userID uuid.UUID) ([]model.UserWithdrawals, error) {
	return service.repository.GetWithdrawals(userID)
}
