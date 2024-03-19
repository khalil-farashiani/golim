package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"os"
	"strings"
)

//go:embed schema.sql
var ddl string

func initDB(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite3", "golim.sqlite")
	if err != nil {
		panic(err)
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil && !strings.Contains(err.Error(), "already exists") {
		panic(err)
	}
	return db
}

func main() {
	ctx := context.Background()

	db := initDB(ctx)
	limiter, err := initFlags(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	data, err := limiter.ExecCMD(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if data != nil {
		fmt.Fprintf(os.Stdout, "%v", data)
		return
	}
	fmt.Printf("%s", "DONE")
}

// initFlags get command and flags from std input to create a golim or role
func initFlags(ctx context.Context) (*golim, error) {
	golim := newLimiter()

	rootFlags := ff.NewFlagSet("golim")
	rootCmd := &ff.Command{
		Name:  "golim",
		Usage: "golim [COMMANDS] <FLAGS>",
		Flags: rootFlags,
		Exec: func(ctx context.Context, args []string) error {
			return nil
		},
	}

	helpCMD := golim.createHelpCMD()
	initCMD := golim.createInitCMD()
	addCMD := golim.createAddCMD()
	removeCMD := golim.createRemoveCMD()
	getCMD := golim.createGetRolesCMD()
	removeLimiterCMD := golim.createRemoveCMD()

	rootCmd.Subcommands = []*ff.Command{helpCMD, initCMD, addCMD, removeCMD, getCMD, removeLimiterCMD}
	if err := rootCmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		return nil, fmt.Errorf("%s\n%s", ffhelp.Command(rootCmd), err)
	}

	return golim, nil
}
