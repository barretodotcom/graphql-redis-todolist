package graph

import "github.com/barretodotcom/graphql-redis-todolist/internal/services"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserService *services.UserService
	TodoService *services.TodoService
}
