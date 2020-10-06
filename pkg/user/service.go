package user

import (
	"github.com/nozgurozturk/noo-analytics/entities"
	errors "github.com/nozgurozturk/noo-analytics/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
)

type Service interface {
	Create(u *entities.User) (user *entities.UserResponse, error *errors.ApplicationError)
	FindByEmail(email string) (user *entities.User, error *errors.ApplicationError)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func validate(user *entities.User) *errors.ApplicationError {

	user.Name = strings.TrimSpace(user.Name)
	if user.Name == "" {

		return errors.BadRequest("Name is required")
	}

	if len(user.Name) < 6 {
		return errors.BadRequest("Name must be minimum 6 character long")
	}

	user.Password = strings.TrimSpace(user.Password)
	if user.Password == "" {
		return errors.BadRequest("Password is required")
	}

	if len(user.Password) < 8 {
		return errors.BadRequest("Password must be minimum 8 character long")
	}

	user.Email = strings.TrimSpace(user.Email)
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if user.Email == "" {
		return errors.BadRequest("Email is required")
	}

	if !emailRegex.MatchString(user.Email) {
		return errors.BadRequest("Email is not valid")
	}

	return nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword string, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(userPassword))
}

func (s *service) Create(u *entities.User) (user *entities.UserResponse, error *errors.ApplicationError) {
	validationError := validate(u)

	if validationError != nil {
		return nil, validationError
	}

	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return nil, errors.InternalServer(err.Error())
	}

	u.Password = string(hashedPassword)

	created, appErr := s.repository.Create(u)

	if appErr != nil {
		return nil, appErr
	}
	createdUser := &entities.UserResponse{Email: created.Email, Name: created.Name}
	return createdUser, nil

}

func (s *service) FindByEmail(email string) (user *entities.User, error *errors.ApplicationError) {
	if email == "" {
		return nil, errors.BadRequest("Email is required")
	}

	found, err := s.repository.FindByEmail(email)

	if err != nil {
		return nil, err
	}

	return found, nil
}
