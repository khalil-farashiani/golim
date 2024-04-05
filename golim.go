package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/khalil-farashiani/golim/role"
	"github.com/peterbourgon/ff/v4"
)

type limiterRole struct {
	operation    string
	limiterID    int
	endPoint     string
	method       string
	bucketSize   int
	initialToken int
	addToken     int64
}

type limiter struct {
	id          interface{}
	name        string
	destination string
	operation   string
}

type Store struct {
	db    *role.Queries
	cache *cache
}

type golim struct {
	limiter     *limiter
	limiterRole *limiterRole
	port        int64
	skip        bool
	*logger
	Store
}

func (g *golim) getRole(ctx context.Context) (role.GetRoleRow, bool, error) {
	params := toGetRole(g)
	data := g.cache.getLimiter(ctx, params)
	if data != nil {
		return *data, true, nil
	}

	row, err := g.db.GetRole(ctx, params)
	if err != nil {
		if strings.Contains(err.Error(), notFoundSqlError) {
			return role.GetRoleRow{}, false, nil
		}
		return role.GetRoleRow{}, false, err
	}

	if row.Endpoint == "" {
		return role.GetRoleRow{}, false, nil
	}

	go g.cache.setLimiter(ctx, &params, &row)

	return row, true, nil
}

func (g *golim) getRoles(ctx context.Context) ([]role.GetRolesRow, error) {
	return g.db.GetRoles(ctx, int64(g.limiterRole.limiterID))
}

func (g *golim) addRole(ctx context.Context) error {
	params := toCreateRoleParam(g)
	_, err := g.db.CreateRole(ctx, params)
	return err
}

func (g *golim) removeRole(ctx context.Context) error {
	return g.db.DeleteRole(ctx, int64(g.limiterRole.limiterID))
}

func (g *golim) createRateLimiter(ctx context.Context) error {
	params := toCreateRateLimiter(g)
	_, err := g.db.CrateRateLimiter(ctx, params)
	return err
}

func (g *golim) removeRateLimiter(ctx context.Context) error {
	return g.db.DeleteRateLimiter(ctx, int64(g.limiter.id.(int)))
}

func (g *golim) ExecCMD(ctx context.Context) (interface{}, error) {

	if g.port != 0 {
		go runCronTasks(ctx, g)
		return startServer(g)
	}
	if g.limiter != nil {
		return handleLimiterOperation(g, ctx)
	}
	if g.limiterRole != nil {
		return handleLimiterRoleOperation(g, ctx)
	}
	return nil, nil
}

func handleLimiterOperation(g *golim, ctx context.Context) (interface{}, error) {
	switch g.limiter.operation {
	case createLimiterOperation:
		return nil, g.createRateLimiter(ctx)
	case removeLimiterOperation:
		return nil, g.removeRateLimiter(ctx)
	}
	return nil, errors.New(unknownLimiterError)
}

func handleLimiterRoleOperation(g *golim, ctx context.Context) (interface{}, error) {
	switch g.limiterRole.operation {
	case addRoleOperation:
		return nil, g.addRole(ctx)
	case removeRoleOperation:
		return nil, g.removeRole(ctx)
	case getRolesOperation:
		return g.getRoles(ctx)
	}
	return nil, errors.New(unknownLimiterRoleError)
}

func newLimiter(db *sql.DB, cache *cache, logger *logger) *golim {
	return &golim{
		logger: logger,
		Store: Store{
			db:    role.New(db),
			cache: cache,
		},
	}
}

func (g *golim) createHelpCMD() *ff.Command {
	helpFlags := ff.NewFlagSet("help")
	return &ff.Command{
		Name:      "help",
		Usage:     "golim help",
		ShortHelp: "Displays help information for golim",
		Flags:     helpFlags,
		Exec: func(ctx context.Context, args []string) error {
			fmt.Println(helpMessageUsage)
			return nil
		},
	}
}

func (g *golim) createRunCMD() *ff.Command {
	runFlags := ff.NewFlagSet("run")
	portNumber := runFlags.Int('p', "port", 8080, "The name of the golim to initialize")
	return &ff.Command{
		Name:      "run",
		Usage:     "golim run -p <port number>",
		ShortHelp: "Initializes a standalone rate golim",
		Flags:     runFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			g.port = int64(*portNumber)
			g.skip = true
			return nil
		},
	}
}

func (g *golim) createInitCMD() *ff.Command {
	initFlags := ff.NewFlagSet("init")
	limiterName := initFlags.String('n', "name", "", "The name of the golim to initialize")
	destinationAddress := initFlags.String('d', "destination", "", "The name of the golim to initialize")
	return &ff.Command{
		Name:      "init",
		Usage:     "golim init -n <limiter_name>",
		ShortHelp: "Initializes a standalone rate golim",
		Flags:     initFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterName != "" && *destinationAddress != "" {
				g.limiter = &limiter{
					name:        *limiterName,
					destination: *destinationAddress,
					operation:   createLimiterOperation,
				}
			} else {
				return errors.New(requiredNameDestinationError)
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) addRemoveLimiterCMD() *ff.Command {
	removeFlags := ff.NewFlagSet("removel")
	limiterID := removeFlags.Int('l', "limiter", 0, "The name of the golim to initialize")
	return &ff.Command{
		Name:      "removel",
		Usage:     "golim removel -l <limiter_id>",
		ShortHelp: "Initializes a standalone rate golim",
		Flags:     removeFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterID != 0 {
				g.limiter = &limiter{
					id:        *limiterID,
					operation: removeLimiterOperation,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) createAddCMD() *ff.Command {
	addFlags := ff.NewFlagSet("add")
	limiterID := addFlags.Int('l', "limiter", 0, "The limiter id")
	endpoint := addFlags.String('e', "endpoint", "", "The endpoint address")
	method := addFlags.String('m', "method", "GET", "The endpoint method")
	bucketSize := addFlags.Int('b', "bsize", 100, "The initial bucket size")
	addToken := addFlags.Int('a', "add_token", 60, "The number of tokens to add per minute")
	initialToken := addFlags.Int('i', "initial_token", 100, "The number of tokens to add per minute")

	return &ff.Command{
		Name:      "add",
		Usage:     "golim add -e <endpoint> -b <bsize> -a <add_token>",
		ShortHelp: "Adds a new golim with the specified configuration",
		Flags:     addFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterID != 0 && *endpoint != "" {
				g.limiterRole = &limiterRole{
					operation:    addRoleOperation,
					limiterID:    *limiterID,
					endPoint:     *endpoint,
					bucketSize:   *bucketSize,
					addToken:     int64(*addToken),
					initialToken: *initialToken,
					method:       *method,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) createRemoveCMD() *ff.Command {
	removeFlags := ff.NewFlagSet("remove")
	roleID := removeFlags.Int('i', "role_id", 0, "the role id")

	return &ff.Command{
		Name:      "remove",
		Usage:     "golim remove -i <role id>",
		ShortHelp: "Adds a new golim with the specified configuration",
		Flags:     removeFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *roleID != 0 {
				g.limiterRole = &limiterRole{
					operation: removeRoleOperation,
					limiterID: *roleID,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) createGetRolesCMD() *ff.Command {
	getFlags := ff.NewFlagSet("get")
	limiterID := getFlags.Int('l', "limiter", 0, "The limiter id")

	return &ff.Command{
		Name:      "get",
		Usage:     "golim get -l <limiter id>",
		ShortHelp: "Adds a new golim with the specified configuration",
		Flags:     getFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterID != 0 {
				g.limiterRole = &limiterRole{
					operation: getRolesOperation,
					limiterID: *limiterID,
				}
			} else {
				return errors.New(requiredLimiterIDError)
			}
			g.skip = true
			return nil
		},
	}
}

func toCreateRoleParam(g *golim) role.CreateRoleParams {
	return role.CreateRoleParams{
		Endpoint:       g.limiterRole.endPoint,
		Operation:      g.limiterRole.method,
		BucketSize:     int64(g.limiterRole.bucketSize),
		AddTokenPerMin: g.limiterRole.addToken,
		InitialTokens:  int64(g.limiterRole.initialToken),
		RateLimiterID:  int64(g.limiterRole.limiterID),
	}
}

func toCreateRateLimiter(g *golim) role.CrateRateLimiterParams {
	return role.CrateRateLimiterParams{
		Name:        g.limiter.name,
		Destination: g.limiter.destination,
	}
}

func toGetRole(g *golim) role.GetRoleParams {
	return role.GetRoleParams{
		Endpoint:  g.limiterRole.endPoint,
		Operation: g.limiterRole.operation,
	}
}
