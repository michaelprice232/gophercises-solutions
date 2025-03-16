package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"tasks/internal/model"
)

type PostgresDB struct {
	dbURL string
	conn  *pgx.Conn
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	return &PostgresDB{dbURL: dbURL}, nil
}

func (db *PostgresDB) Connect() error {
	conn, err := pgx.Connect(context.Background(), db.dbURL)
	if err != nil {
		return fmt.Errorf("creating postgres connection: %w", err)
	}
	db.conn = conn

	return nil
}

func (db *PostgresDB) Close() error {
	if db.conn != nil {
		err := db.conn.Close(context.Background())
		if err != nil {
			return fmt.Errorf("closing postgres connection: %w", err)
		}
	}
	return nil
}

func (db *PostgresDB) AddTask(name string) error {
	if db.conn == nil {
		return fmt.Errorf("db connection not initialized")
	}

	_, err := db.conn.Exec(context.Background(), `INSERT INTO tasks (name) VALUES ($1)`, name)
	if err != nil {
		return fmt.Errorf("inserting task into tasks table: %w", err)
	}

	return nil
}

func (db *PostgresDB) CompleteTask(_ int) error {

	return nil
}

func (db *PostgresDB) ListOutstandingTasks() (model.Tasks, error) {

	return model.Tasks{}, nil
}
