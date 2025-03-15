package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"gophercises-ex7/internal/model"
)

type PostgresDB struct {
	conn *pgx.Conn
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	// todo: move this into the db package or keep simple?
	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("creating postgres connection: %w", err)
	}

	return &PostgresDB{conn: conn}, nil
}

func (db *PostgresDB) Connect() error {
	return nil
}

func (db *PostgresDB) AddTask(name string) error {
	_, err := db.conn.Exec(context.Background(), `INSERT INTO tasks (name) VALUES ($1)`, name)
	if err != nil {
		return fmt.Errorf("inserting task into tasks table: %w", err)
	}

	defer db.conn.Close(context.Background())

	return nil
}

func (db *PostgresDB) CompleteTask(_ int) error {

	return nil
}

func (db *PostgresDB) ListOutstandingTasks() (model.Tasks, error) {

	return model.Tasks{}, nil
}
