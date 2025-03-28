package main

import (
	"context"
	"log/slog"
	"os"
	"strings"

	"tasks/internal/repository"
	"tasks/internal/service"

	"github.com/spf13/cobra"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		slog.Error("DATABASE_URL env var must be set.", "format", "postgresql://username:password@host:port/db?sslmode=disable")
		os.Exit(1)
	}

	db, err := repository.NewPostgresDB(dbURL)
	if err != nil {
		slog.Error("unable to create Postgres client", "error", err)
		os.Exit(1)
	}

	ctx := context.Background()

	s, err := service.NewService(ctx, db)
	if err != nil {
		slog.Error("unable to create service", "error", err)
		os.Exit(1)
	}

	defer func() {
		if err = s.DB.Close(ctx); err != nil {
			slog.Error("problem closing database connection", "error", err)
		}
	}()

	var addCmd = &cobra.Command{
		Use:   "add [task to add]",
		Short: "Add a task to the task list",
		Long:  "Tasks are added to a task list and then can be completed and removed later",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err = s.AddTask(ctx, strings.Join(args, " "))
			if err != nil {
				return err
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
			task := args[0]
			err = s.CompleteTask(ctx, task)
			if err != nil {
				return err
			}

			return nil
		},
	}

	var outstandingCmd = &cobra.Command{
		Use:   "outstanding",
		Short: "Lists the outstanding tasks in the task list",
		Long:  "Lists the tasks which are waiting to be actioned",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err = s.ListOutStandingTasks(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	}

	var completedCmd = &cobra.Command{
		Use:   "completed",
		Short: "Lists the tasks which have been completed today",
		Long:  "Lists the tasks which have been completed today",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			err = s.ListCompletedTasks(ctx)
			if err != nil {
				return err
			}

			return nil
		},
	}

	var rootCmd = &cobra.Command{Use: "task"}
	rootCmd.AddCommand(addCmd, doCmd, outstandingCmd, completedCmd)
	if err = rootCmd.Execute(); err != nil {
		slog.Error("task returned an error.", "error", err)
		os.Exit(1)
	}
}
