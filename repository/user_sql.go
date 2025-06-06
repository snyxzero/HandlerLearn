package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"pasha/models"
)

type UserSQLRepository struct {
	conn *pgx.Conn
}

func NewUserSQLRepository(ctx context.Context) *UserSQLRepository {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	return &UserSQLRepository{
		conn: conn,
	}
}

func (obj *UserSQLRepository) GetUser(ctx context.Context, id int) (*models.User, error) {
	user := &models.User{}
	err := obj.conn.QueryRow(ctx, `
SELECT id, email, name, age
FROM users
WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.Name, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id %d from db: %v ", id, err)
	}
	return user, nil
}

func (obj *UserSQLRepository) AddUser(ctx context.Context, user models.User) (int, error) {
	err := obj.conn.QueryRow(ctx, `
INSERT INTO users (email, name, age)
VALUES ($1, $2, $3)
RETURNING id`, user.Email, user.Name, user.Age).Scan(&user.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to add user to db: %v", err)
	}
	return user.ID, nil
}

func (obj *UserSQLRepository) UpdateUser(ctx context.Context, user models.User) error {
	_, err := obj.conn.Exec(ctx, `
UPDATE users
SET email = $2, name = $3, age = $4
WHERE id = $1`, user.ID, user.Email, user.Name, user.Age)
	if err != nil {
		return fmt.Errorf("failed to update user in db: %v", err)
	}
	return nil
}

func (obj *UserSQLRepository) DeleteUser(ctx context.Context, id int) error {
	_, err := obj.conn.Exec(ctx, `
DELETE FROM users
WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user in db: %v", err)
	}
	return nil
}

func (obj *UserSQLRepository) Close(ctx context.Context) {
	obj.conn.Close(ctx)
}
