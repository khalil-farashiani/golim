package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
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
	cache := initRedis()
	limiter, err := initFlags(ctx, db, cache)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	data, err := limiter.ExecCMD(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	if data != nil {
		makeTable(toSlice(data))
		fmt.Fprintf(os.Stdout, "DONE")
		return
	}
	fmt.Printf("DONE")
}

// initFlags get command and flags from std input to create a golim or role
func initFlags(ctx context.Context, db *sql.DB, cache *cache) (*golim, error) {
	golim := newLimiter(db, cache)

	rootCmd := createRootCommand(golim)
	if err := rootCmd.ParseAndRun(ctx, os.Args[1:]); err != nil {
		return nil, fmt.Errorf("%s\n%s", ffhelp.Command(rootCmd), err)
	}

	return golim, nil
}

func createRootCommand(g *golim) *ff.Command {
	rootFlags := ff.NewFlagSet("golim")

	helpCMD := g.createHelpCMD()
	initCMD := g.createInitCMD()
	addCMD := g.createAddCMD()
	removeCMD := g.createRemoveCMD()
	getCMD := g.createGetRolesCMD()
	removeLimiterCMD := g.createRemoveCMD()
	runCMD := g.createRunCMD()

	rootCmd := &ff.Command{
		Name:        "golim",
		Usage:       "golim [COMMANDS] <FLAGS>",
		Flags:       rootFlags,
		Subcommands: []*ff.Command{helpCMD, initCMD, addCMD, removeCMD, getCMD, removeLimiterCMD, runCMD},
		Exec: func(ctx context.Context, args []string) error {
			return nil
		},
	}

	return rootCmd
}
