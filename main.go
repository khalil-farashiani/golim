package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/peterbourgon/ff/v4"
	"github.com/peterbourgon/ff/v4/ffhelp"
	"os"
)

//go:embed schema.sql
var ddl string

func initDB(ctx context.Context) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		panic(err)
	}
	return db
}

func main() {
	ctx := context.Background()

	db := initDB(ctx)
	limiter, err := initFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	err = limiter.ExecCMD(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// initFlags get command and flags from std input to create a golim or role
func initFlags() (*golim, error) {
	golim := newLimiter()

	rootFlags := ff.NewFlagSet("golim")
	rootCmd := &ff.Command{
		Name:  "golim",
		Usage: "golim [COMMANDS] <FLAGS>",
		Flags: rootFlags,
	}

	helpCMD := golim.createHelpCMD()
	initCMD := golim.createInitCMD()
	addCMD := golim.createAddCMD()
	removeCMD := golim.createRemoveCMD()
	getCMD := golim.createGetRolesCMD()
	removeLimiterCMD := golim.createRemoveCMD()

	rootCmd.Subcommands = []*ff.Command{helpCMD, initCMD, addCMD, removeCMD, getCMD, removeLimiterCMD}
	if err := rootCmd.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("%s\n%s", ffhelp.Command(rootCmd), err)
	}

	return golim, nil
}
