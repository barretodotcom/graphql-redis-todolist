package repositories

import (
	"database/sql"
	"errors"

	"github.com/barretodotcom/graphql-redis-todolist/internal/entities"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (repo *UserRepository) CreateUser(user entities.User) error {
	query := "INSERT INTO user (id, username, password) VALUES (?, ?, ?)"

	_, err := repo.DB.Exec(query, user.ID, user.Username, user.Password)

	return err
}

func (repo *UserRepository) GetUserByUsername(username string) (*entities.User, error) {
	query := "SELECT id, password FROM user WHERE username = ?"

	var id, password string
	row := repo.DB.QueryRow(query, username)

	err := row.Scan(&id, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{ID: id, Username: username, Password: password}, nil
}

func (repo *UserRepository) FindUsers() ([]entities.User, error) {
	query := "SELECT * FROM user"

	rows, err := repo.DB.Query(query)

	if err != nil {
		return nil, err
	}

	var users []entities.User

	for rows.Next() {
		var id, username, password string

		err := rows.Scan(&id, &username, &password)
		if err != nil {
			return nil, err
		}

		users = append(users, entities.User{ID: id, Username: username, Password: password})
	}

	return users, nil
}

func (repo *UserRepository) FindUserById(userId string) (*entities.User, error) {
	query := "SELECT * FROM user WHERE id = ?"

	row := repo.DB.QueryRow(query, userId)

	var id, username, password string

	err := row.Scan(&id, &username, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{ID: id, Username: username, Password: password}, nil
}

func (repo *UserRepository) FindUserByUsername(username string) (*entities.User, error) {
	query := "SELECT id, password FROM user WHERE username = ?"

	row := repo.DB.QueryRow(query, username)

	var id, password string

	err := row.Scan(&id, &password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &entities.User{ID: id, Username: username, Password: password}, nil
}
