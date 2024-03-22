package main

import (
	"context"
	"database/sql"
	"errors"
	"github.com/khalil-farashiani/golim/role"
	"github.com/peterbourgon/ff/v4"
)

type limiterRole struct {
	operation    int
	limiterID    int
	endPoint     string
	bucketSize   int
	initialToken int
	addToken     int
}

type limiter struct {
	id          interface{}
	name        string
	destination string
	operation   int
}

type golim struct {
	limiter     *limiter
	limiterRole *limiterRole
	port        int64
	skip        bool
}

var operationIdToString = map[int]string{
	1: "GET",
	2: "POST",
	3: "PUT",
	4: "PATCH",
	5: "DELETE",
}

func (g *golim) getRole(ctx context.Context, db *sql.DB, cache *cache) (role.GetRoleRow, error) {
	params := toGetRole(g)
	query := role.New(db)
	data := cache.getLimiter(ctx, params)
	if data != nil {
		return *data, nil
	}
	row, err := query.GetRole(ctx, params)
	if err != nil {
		return role.GetRoleRow{}, err
	}
	go func() {
		cache.setLimiter(ctx, &params, &row)
	}()
	return row, nil
}

func (g *golim) getRoles(ctx context.Context, db *sql.DB) ([]role.GetRolesRow, error) {
	query := role.New(db)
	return query.GetRoles(ctx, int64(g.limiterRole.limiterID))
}

func (g *golim) addRole(ctx context.Context, db *sql.DB) error {
	params := toCreateRoleParam(g)
	query := role.New(db)
	_, err := query.CreateRole(ctx, params)
	return err
}

func (g *golim) removeRole(ctx context.Context, db *sql.DB) error {
	query := role.New(db)
	return query.DeleteRole(ctx, int64(g.limiterRole.limiterID))
}

func (g *golim) createRateLimiter(ctx context.Context, db *sql.DB) error {
	params := toCreateRateLimiter(g)
	query := role.New(db)
	_, err := query.CrateRateLimiter(ctx, params)
	return err
}

func (g *golim) removeRateLimiter(ctx context.Context, db *sql.DB) error {
	query := role.New(db)
	return query.DeleteRateLimiter(ctx, g.limiter.id.(int64))
}

func (g *golim) ExecCMD(ctx context.Context, db *sql.DB, cache *cache) (interface{}, error) {

	if g.port != 0 {
		return startServer(g, db, cache)
	}
	if g.limiter != nil {
		switch g.limiter.operation {
		case createLimiterOperationID:
			return nil, g.createRateLimiter(ctx, db)
		case removeLimiterOperationID:
			return nil, g.removeRateLimiter(ctx, db)
		}
	}
	if g.limiterRole != nil {
		switch g.limiterRole.operation {
		case addRoleOperationID:
			return nil, g.addRole(ctx, db)
		case removeRoleOperationID:
			return nil, g.removeRole(ctx, db)
		case getRolesOperationID:
			return g.getRoles(ctx, db)
		}
	}
	return nil, errors.New("unsupported operation")
}

func newLimiter() *golim {
	return &golim{}
}

func (g *golim) createHelpCMD() *ff.Command {
	return &ff.Command{
		Name:      "help",
		Usage:     "golim help",
		ShortHelp: "Displays help information for golim",
		Flags:     ff.NewFlagSet("help"),
		Exec: func(ctx context.Context, args []string) error {
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
					operation:   createLimiterOperationID,
				}
			} else {
				return errors.New("name and destination is required")
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) addRemoveLimiterCMD() *ff.Command {
	initFlags := ff.NewFlagSet("remove-limiter")
	limiterID := initFlags.Int('l', "limiter", 0, "The name of the golim to initialize")
	return &ff.Command{
		Name:      "init",
		Usage:     "golim init -n <limiter_name>",
		ShortHelp: "Initializes a standalone rate golim",
		Flags:     initFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterID != 0 {
				g.limiter = &limiter{
					id:        *limiterID,
					operation: removeLimiterOperationID,
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
			if *limiterID == 0 || *endpoint == "" {
				g.limiterRole = &limiterRole{
					operation:    addRoleOperationID,
					limiterID:    *limiterID,
					endPoint:     *endpoint,
					bucketSize:   *bucketSize,
					addToken:     *addToken,
					initialToken: *initialToken,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *golim) createRemoveCMD() *ff.Command {
	removeFlags := ff.NewFlagSet("remove")
	limiterID := removeFlags.Int('l', "limiter", 0, "The limiter id")

	return &ff.Command{
		Name:      "add",
		Usage:     "golim add -e <endpoint> -b <bsize> -a <add_token>",
		ShortHelp: "Adds a new golim with the specified configuration",
		Flags:     removeFlags,
		Exec: func(ctx context.Context, args []string) error {
			if g.skip {
				return nil
			}
			if *limiterID == 0 {
				g.limiterRole = &limiterRole{
					operation: removeRoleOperationID,
					limiterID: *limiterID,
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
					operation: getRolesOperationID,
					limiterID: *limiterID,
				}
			} else {
				return errors.New("limiter id is required")
			}
			g.skip = true
			return nil
		},
	}
}

func toCreateRoleParam(g *golim) role.CreateRoleParams {
	return role.CreateRoleParams{
		Endpoint:       g.limiterRole.endPoint,
		Operation:      operationIdToString[g.limiterRole.operation],
		BucketSize:     int64(g.limiterRole.bucketSize),
		AddTokenPerMin: int64(g.limiterRole.addToken),
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
		Operation: operationIdToString[g.limiterRole.operation],
	}
}

func toRole(row role.GetRoleRow) role.Role {
	return role.Role{
		Endpoint:       row.Endpoint,
		Operation:      row.Operation,
		BucketSize:     row.BucketSize,
		AddTokenPerMin: row.AddTokenPerMin,
		InitialTokens:  row.InitialTokens,
	}
}
