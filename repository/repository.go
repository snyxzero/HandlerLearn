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

type PostgresSQLRepository struct {
	conn *pgx.Conn
}

func NewPostgresSQLRepository() *PostgresSQLRepository {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("ping failed: %v", err)
	}

	return &PostgresSQLRepository{
		conn: conn,
	}
}

func (obj *PostgresSQLRepository) GetUser(id int) (*models.User, error) {
	user := &models.User{}

	err := obj.conn.QueryRow(context.Background(), `
SELECT id, email, name, age
FROM users
WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.Name, &user.Age)
	if err != nil {
		return nil, fmt.Errorf("failed to get user id %d from db: %v ", id, err)
	}
	return user, nil
}

func (obj *PostgresSQLRepository) AddUser(user models.User) error {
	_, err := obj.conn.Exec(context.Background(), `
INSERT INTO users (email, name, age)
VALUES ($1, $2, $3)`, user.Email, user.Name, user.Age)
	if err != nil {
		return fmt.Errorf("failed to add user to db: %v", err)
	}
	return nil
}

func (obj *PostgresSQLRepository) UpdateUser(user models.User) error {
	_, err := obj.conn.Exec(context.Background(), `
UPDATE users
SET email = $2, name = $3, age = $4
WHERE id = $1`, user.ID, user.Email, user.Name, user.Age)
	if err != nil {
		return fmt.Errorf("failed to update user in db: %v", err)
	}
	return nil
}

func (obj *PostgresSQLRepository) DeleteUser(id int) error {
	_, err := obj.conn.Exec(context.Background(), `
DELETE FROM users
WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user in db: %v", err)
	}
	return nil
}

func (obj *PostgresSQLRepository) Close() {
	obj.conn.Close(context.Background())
}
