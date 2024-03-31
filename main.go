package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil && !strings.Contains(err.Error(), "already exists") {
		log.Fatal(err)
	}
	return db
}

// everything start from main
func main() {
	ctx := context.Background()
	db := initDB(ctx)
	cache := initRedis()
	limiter, err := initFlags(ctx, db, cache)
	if err != nil {
		log.Fatalf("Error initializing limiter: %v", err)
	}

	data, err := limiter.ExecCMD(ctx)
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
	if data != nil {
		makeTable(toSlice(data))
	}
	fmt.Println("DONE")
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
	removeLimiterCMD := g.addRemoveLimiterCMD()
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
