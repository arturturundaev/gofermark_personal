package user

import (
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

	err = service.compareHashAndPassword([]byte(user.Password), password)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (service *UserService) hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (service *UserService) compareHashAndPassword(hashedPassword []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
}
