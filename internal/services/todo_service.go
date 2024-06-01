package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/barretodotcom/graphql-redis-todolist/graph/model"
	"github.com/barretodotcom/graphql-redis-todolist/internal/cache"
	"github.com/barretodotcom/graphql-redis-todolist/internal/entities"
	"github.com/barretodotcom/graphql-redis-todolist/internal/repositories"
	"github.com/barretodotcom/graphql-redis-todolist/pkg/date"
	"github.com/google/uuid"
)

type TodoService struct {
	repo         *repositories.TodoRepository
	redisService *cache.RedisService
}

func NewTodoService(repo *repositories.TodoRepository, redisService *cache.RedisService) *TodoService {
	return &TodoService{
		repo:         repo,
		redisService: redisService,
	}
}

func (s *TodoService) FindTodosByUserId(userId string) ([]*model.Todo, error) {
	userTodosRedisKey := s.getUserTodosRedisKey(userId)
	cachedValue, err := s.redisService.Get(userTodosRedisKey)
	if err != nil {
		return nil, err
	}

	var todosModel []*model.Todo
	if cachedValue != "" {
		err := json.Unmarshal([]byte(cachedValue), &todosModel)

		return todosModel, err
	}

	todos, err := s.repo.FindTodosByUserId(userId)
	if err != nil {
		return nil, err
	}

	for _, todo := range todos {
		todosModel = append(todosModel, &model.Todo{ID: todo.ID, Title: todo.Title, StartDate: todo.StartDate.String(), EndDate: todo.EndDate.String()})
	}

	todosModelJson, _ := json.Marshal(todosModel)

	return todosModel, s.redisService.Set(todosModelJson, userTodosRedisKey)
}

func (s *TodoService) CreateTodo(input model.NewTodo, userID string) (*model.Todo, error) {
	existingTodo, err := s.repo.FindTodoByTitle(input.Title, userID)
	if err != nil {
		return nil, err
	}

	if existingTodo != nil {
		return nil, errors.New("Todo already exists")
	}

	id := uuid.New().String()

	startDate, err := date.ParseStringToDate(input.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := date.ParseStringToDate(input.EndDate)
	if err != nil {
		return nil, err
	}

	todo := entities.Todo{ID: id, Title: input.Title, StartDate: startDate, EndDate: endDate, UserID: userID}
	err = s.repo.CreateTodo(todo)
	if err != nil {
		return nil, err
	}

	err = s.redisService.Delete(s.getUserTodosRedisKey(userID))
	if err != nil {
		return nil, err
	}

	return &model.Todo{ID: todo.ID, Title: input.Title, StartDate: input.StartDate, EndDate: input.EndDate}, err
}
func (s *TodoService) DeleteTodoById(id string, userID string) (bool, error) {
	todo, err := s.repo.FindTodoById(id)
	if err != nil {
		return false, err
	}

	if todo == nil {
		return false, errors.New("Todo does not exists")
	}

	err = s.repo.DeleteTodoById(id)
	if err != nil {
		return false, err
	}

	err = s.redisService.Delete(s.getUserTodosRedisKey(userID))
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *TodoService) getUserTodosRedisKey(userId string) string {
	return fmt.Sprintf("user_todos:%s", userId)
}
