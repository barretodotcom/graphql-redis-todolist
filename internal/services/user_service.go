package services

import (
	"errors"

	"github.com/barretodotcom/graphql-redis-todolist/graph/model"
	"github.com/barretodotcom/graphql-redis-todolist/internal/cache"
	"github.com/barretodotcom/graphql-redis-todolist/internal/entities"
	"github.com/barretodotcom/graphql-redis-todolist/internal/repositories"
	"github.com/barretodotcom/graphql-redis-todolist/pkg/hash"
	"github.com/barretodotcom/graphql-redis-todolist/pkg/jwt"
	"github.com/google/uuid"
)

type UserService struct {
	repo         *repositories.UserRepository
	redisService *cache.RedisService
}

func NewUserService(repo *repositories.UserRepository, redisService *cache.RedisService) *UserService {
	return &UserService{
		repo:         repo,
		redisService: redisService,
	}
}

func (service *UserService) CreateUser(username string, password string) (*model.User, error) {
	user, err := service.repo.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, errors.New("User already exists.")
	}

	id := uuid.New().String()

	hashedPassword, err := hash.HashPassword(password)
	if err != nil {
		return nil, err
	}

	newUser := entities.User{ID: id, Username: username, Password: string(hashedPassword)}

	err = service.repo.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return &model.User{ID: id, Username: username, Password: password}, nil
}

func (service *UserService) FindUsers() ([]*model.User, error) {
	users, err := service.repo.FindUsers()
	if err != nil {
		return nil, err
	}

	var usersModel []*model.User

	for _, user := range users {
		usersModel = append(usersModel, &model.User{ID: user.ID, Username: user.Username, Password: user.Password})
	}

	return usersModel, nil
}

func (service *UserService) FindUserById(id string) (*model.User, error) {
	user, err := service.repo.FindUserById(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, nil
	}

	return &model.User{ID: user.ID, Username: user.Username, Password: user.Password}, nil
}

func (service *UserService) AuthUser(login model.Login) (string, error) {
	user, err := service.repo.FindUserByUsername(login.Username)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("Invalid data")
	}

	validPassword := hash.CheckPasswordHash(user.Password, login.Password)
	if !validPassword {
		return "", errors.New("Invalid data")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return "", nil
	}

	return token, nil
}
