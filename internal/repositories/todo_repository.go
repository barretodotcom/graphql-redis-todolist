package repositories

import (
	"database/sql"
	"errors"

	"github.com/barretodotcom/graphql-redis-todolist/internal/entities"
	"github.com/barretodotcom/graphql-redis-todolist/pkg/date"
)

type TodoRepository struct {
	DB *sql.DB
}

func NewTodoRepository(db *sql.DB) *TodoRepository {
	return &TodoRepository{
		DB: db,
	}
}

func (r *TodoRepository) FindTodosByUserId(userId string) ([]entities.Todo, error) {
	query := "SELECT id, title, startDate, endDate FROM todo WHERE userId = ?"

	rows, err := r.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	var todos []entities.Todo

	var id, title, startDateString, endDateString string
	for rows.Next() {

		err := rows.Scan(&id, &title, &startDateString, &endDateString)
		if err != nil {
			return nil, err
		}

		startDate, err := date.ParseStringToDate(startDateString)
		if err != nil {
			return nil, err
		}
		endDate, err := date.ParseStringToDate(startDateString)
		if err != nil {
			return nil, err
		}

		todos = append(todos, entities.Todo{ID: id, Title: title, StartDate: startDate, EndDate: endDate})

	}

	return todos, nil
}

func (r *TodoRepository) FindTodoByTitle(title string, userId string) (*entities.Todo, error) {
	query := "SELECT id, startDate, endDate FROM todo WHERE title = ? AND userId = ?"
	row := r.DB.QueryRow(query, title, userId)

	var id, startDateString, endDateString string

	err := row.Scan(&id, &startDateString, &endDateString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	startDate, err := date.ParseStringToDate(startDateString)
	if err != nil {
		return nil, err
	}
	endDate, err := date.ParseStringToDate(startDateString)
	if err != nil {
		return nil, err
	}

	return &entities.Todo{ID: id, Title: title, StartDate: startDate, EndDate: endDate, UserID: userId}, nil
}

func (r *TodoRepository) FindTodoById(id string) (*entities.Todo, error) {
	query := "SELECT title, startDate, endDate FROM todo WHERE id = ?"

	row := r.DB.QueryRow(query, id)

	var title, startDateString, endDateString string
	err := row.Scan(&title, &startDateString, &endDateString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	startDate, err := date.ParseStringToDate(startDateString)
	if err != nil {
		return nil, err
	}
	endDate, err := date.ParseStringToDate(endDateString)
	if err != nil {
		return nil, err
	}

	return &entities.Todo{ID: id, Title: title, StartDate: startDate, EndDate: endDate}, nil
}

func (r *TodoRepository) CreateTodo(todo entities.Todo) error {
	query := "INSERT INTO todo (id, title, startDate, endDate, userId) VALUES (?, ?, ?, ?, ?)"

	_, err := r.DB.Exec(query, todo.ID, todo.Title, todo.StartDate, todo.EndDate, todo.UserID)
	return err
}

func (r *TodoRepository) DeleteTodoById(id string) error {
	query := "DELETE FROM todo WHERE id = ?"

	_, err := r.DB.Exec(query, id)

	return err
}
