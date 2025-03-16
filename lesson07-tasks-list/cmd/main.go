package main

import (
	"fmt"
	"os"
	"strings"

	"tasks/internal/repository"
	"tasks/internal/service"

	"github.com/spf13/cobra"
	"log/slog"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL env var must be set.", "format", "postgresql://username:password@host:port/db?sslmode=disable")
		os.Exit(1)
	}

	var addCmd = &cobra.Command{
		Use:   "add [task to add]",
		Short: "Add a task to the task list",
		Long:  "Tasks are added and then can be completed and removed later",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := repository.NewPostgresDB(dbURL)
			if err != nil {
				return fmt.Errorf("unable to create Postgres client: %w", err)
			}

			s, err := service.NewService(db)
			if err != nil {
				return fmt.Errorf("unable to create service: %w", err)
			}

			err = s.AddTask(strings.Join(args, " "))
			if err != nil {
				return fmt.Errorf("unable to add task: %w", err)
			}

			return nil
		},
	}

	var doCmd = &cobra.Command{
		Use:   "do [task number to complete]",
		Short: "Complete a task from the task list",
		Long:  "Tasks are completed from the task list once they have been actioned",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			input := strings.Join(args, " ")
			slog.Info("Completing task.", "task", input)
			return nil
		},
	}

	var outstandingCmd = &cobra.Command{
		Use:   "outstanding",
		Short: "List the outstanding tasks in the task list",
		Long:  "Lists the tasks which are waiting to be actioned",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Listing outstanding tasks.")
			return nil
		},
	}

	var rootCmd = &cobra.Command{Use: "task"}
	rootCmd.AddCommand(addCmd, doCmd, outstandingCmd)
	if err := rootCmd.Execute(); err != nil {
		slog.Error("task returned an error.", "error", err)
		os.Exit(1)
	}
}
