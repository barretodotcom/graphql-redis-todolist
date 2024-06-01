package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/barretodotcom/graphql-redis-todolist/database"
	"github.com/barretodotcom/graphql-redis-todolist/graph"
	"github.com/barretodotcom/graphql-redis-todolist/internal/cache"
	"github.com/barretodotcom/graphql-redis-todolist/internal/repositories"
	"github.com/barretodotcom/graphql-redis-todolist/internal/services"
	"github.com/barretodotcom/graphql-redis-todolist/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

const defaultPort = "8080"

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(middleware.AuthMiddleware())

	db, err := database.InitDb()
	if err != nil {
		log.Fatalf("Error while initializing db: %v", err)
	}
	database.CreateTables()
	redisClient, err := cache.ConnectRedis()
	if err != nil {
		log.Fatalf("Error while initializing redis: %v", err)
	}

	redisService := cache.NewRedisService(redisClient)

	userRepository := repositories.NewUserRepository(db)
	todoRepository := repositories.NewTodoRepository(db)
	userService := services.NewUserService(userRepository, redisService)
	todoService := services.NewTodoService(todoRepository, redisService)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{UserService: userService, TodoService: todoService}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
