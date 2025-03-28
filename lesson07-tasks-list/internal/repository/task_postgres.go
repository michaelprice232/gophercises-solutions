package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"tasks/internal/model"
)

const (
	dbTimeout = time.Second * 5
)

type PostgresDB struct {
	dbURL string
	conn  *pgx.Conn
}

func NewPostgresDB(dbURL string) (*PostgresDB, error) {
	return &PostgresDB{dbURL: dbURL}, nil
}

func (db *PostgresDB) Connect(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	conn, err := pgx.Connect(ctx, db.dbURL)
	if err != nil {
		return fmt.Errorf("creating postgres connection: %w", err)
	}
	db.conn = conn

	return nil
}

func (db *PostgresDB) Close(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if db.conn != nil {
		err := db.conn.Close(ctx)
		if err != nil {
			return fmt.Errorf("closing postgres connection: %w", err)
		}
	}
	return nil
}

func (db *PostgresDB) AddTask(ctx context.Context, name string) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if db.conn == nil {
		return fmt.Errorf("db connection not initialized")
	}

	tag, err := db.conn.Exec(ctx, `INSERT INTO tasks (name) VALUES ($1)`, name)
	if err != nil {
		return fmt.Errorf("inserting task '%s' into tasks table: %w", name, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows inserted whilst adding task '%s' into tasks table", name)
	}

	return nil
}

func (db *PostgresDB) CompleteTask(ctx context.Context, taskID int) error {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	now := time.Now()

	tag, err := db.conn.Exec(ctx, `
		UPDATE tasks 
		SET completed = true, completed_at = $1 
		WHERE id = $2
		`, now, taskID)
	if err != nil {
		return fmt.Errorf("setting taskID %d to completed: %w", taskID, err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("taskID %d not found", taskID)
	}

	return nil
}

func (db *PostgresDB) ListOutstandingTasks(ctx context.Context) (model.Tasks, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if db.conn == nil {
		return nil, fmt.Errorf("db connection not initialized")
	}

	// Check for errors at the end
	rows, _ := db.conn.Query(ctx, `
		SELECT id, name FROM tasks 
		WHERE completed = false`)
	defer rows.Close()

	results := make(model.Tasks, 0)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("scanning postgres row: %w", err)
		}

		results = append(results, model.Task{
			ID:   id,
			Name: name,
		})
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("scanning postgres rows: %w", rows.Err())
	}

	return results, nil
}

func (db *PostgresDB) ListCompletedTasks(ctx context.Context) (model.Tasks, error) {
	ctx, cancel := context.WithTimeout(ctx, dbTimeout)
	defer cancel()

	if db.conn == nil {
		return nil, fmt.Errorf("db connection not initialized")
	}

	// Check for errors at the end
	rows, _ := db.conn.Query(ctx, `
		SELECT id, name, completed_at FROM tasks 
	 	WHERE completed = true AND DATE(completed_at) = $1
		ORDER BY completed_at DESC 
		`, time.Now())
	defer rows.Close()

	results := make(model.Tasks, 0)
	for rows.Next() {
		var id int
		var name string
		var completedAt time.Time

		err := rows.Scan(&id, &name, &completedAt)
		if err != nil {
			return nil, fmt.Errorf("scanning postgres row: %w", err)
		}

		results = append(results, model.Task{
			ID:          id,
			Name:        name,
			CompletedAt: completedAt,
		})
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("scanning postgres rows: %w", rows.Err())
	}

	return results, nil
}
