package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/khalil-farashiani/golim/internal/contract"
	"github.com/khalil-farashiani/golim/internal/domain"
	"github.com/khalil-farashiani/golim/internal/entity"
	"github.com/khalil-farashiani/golim/internal/store"
	role2 "github.com/khalil-farashiani/golim/internal/store/role"
	"github.com/khalil-farashiani/golim/pkg/log"
	"github.com/peterbourgon/ff/v4"
)

const (
	helpMessageUsage = `
Golim help:
	- golim run -p{--port} <port> [run in the specific port default is 8080]
	- golim get -l{--limiter} <limiter id> [get roles of a rate limiter]
	- golim init -n{--name} foo -d{--destination} 8.8.8.8 [initial new rate limiter]
	- golim add -l{--limiter} <limiter id> -e{--endpoint} <endpoint> -b{--bsize} <bucket size> -a{--add_token} <add_token per minute> -i{--initial_token} <initial tokens> [add specific role to limiter]
	- golim remove -i{--id} <role id> [remove specific role]
	- golim remove-limiter -l{--limiter} <limiter id> [remove specific limiter]`
)

const (
	unknownLimiterRoleError      = "unknown limiter role operation"
	unknownLimiterError          = "unknown limiter operation"
	requiredNameDestinationError = "name and destination is required"
	requiredLimiterIDError       = "limiter id is required"
	slowDownError                = "slow down"
	notFoundSqlError             = "sql: no rows in result set"
)

type Golim struct {
	Cache  contract.Cache
	Logger contract.Logger
	DB     contract.DBStore
}

func New() Golim {
	return Golim{}
}

func (g *Golim) AddCache(cache contract.Cache) *Golim {
	g.Cache = cache
	return g
}

func (g *Golim) AddDB(db contract.DBStore) *Golim {
	g.DB = db
	return g
}

func (g *Golim) AddLogger(logger contract.Logger) *Golim {
	g.Logger = logger
	return g
}

func (g *Golim) GetRole(ctx context.Context, ID int64) (entity.Role, error) {
	data := g.Cache.GetRole(ctx, ID)
	if data != nil {
		return *data, nil
	}

	role, err := g.DB.GetRole(ctx, ID)
	if err != nil {
		return entity.Role{}, err
	}

	go g.Cache.SetRole(ctx, role)

	return role, nil
}

func (g *Golim) getRoles(ctx context.Context, ID int64) ([]entity.Role, error) {
	return g.DB.GetRoles(ctx, ID)
}

func (g *Golim) addRole(ctx context.Context, role entity.Role) error {
	return g.DB.CreateRole(ctx, role)
}

func (g *Golim) removeRole(ctx context.Context, ID int64) error {
	return g.DB.DeleteRole(ctx, ID)
}

func (g *Golim) createRateLimiter(ctx context.Context, limiter entity.Limiter) error {
	return g.DB.CrateRateLimiter(ctx, limiter)
}

func (g *Golim) removeRateLimiter(ctx context.Context, ID int64) error {
	return g.DB.DeleteRateLimiter(ctx, ID)
}

func (g *Golim) DO(ctx context.Context) (interface{}, error) {

	if g.port != 0 {
		go main.runCronTasks(ctx, g)
		return main.startServer(g)
	}
	if g.limiter != nil {
		return handleLimiterOperation(g, ctx)
	}
	if g.limiterRole != nil {
		return handleLimiterRoleOperation(g, ctx)
	}
	return nil, nil
}

func handleLimiterOperation(g *GolimCLI, ctx context.Context) (interface{}, error) {
	switch g.limiter.operation {
	case main.createLimiterOperation:
		return nil, g.createRateLimiter(ctx)
	case main.removeLimiterOperation:
		return nil, g.removeRateLimiter(ctx)
	}
	return nil, errors.New(main.unknownLimiterError)
}

func handleLimiterRoleOperation(g *GolimCLI, ctx context.Context) (interface{}, error) {
	switch g.limiterRole.operation {
	case main.addRoleOperation:
		return nil, g.addRole(ctx)
	case main.removeRoleOperation:
		return nil, g.removeRole(ctx)
	case main.getRolesOperation:
		return g.getRoles(ctx)
	}
	return nil, errors.New(main.unknownLimiterRoleError)
}

func NewLimiter(db *sql.DB, cache *store.Cache, logger *log.Logger) *Golim {
	return &domain.GolimCLI{
		Logger: logger,
		Store: Store{
			db:    role2.New(db),
			cache: cache,
		},
	}
}

func (g *GolimCLI) createHelpCMD() *ff.Command {
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

func (g *GolimCLI) createRunCMD() *ff.Command {
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

func (g *GolimCLI) createInitCMD() *ff.Command {
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
					operation:   main.createLimiterOperation,
				}
			} else {
				return errors.New(main.requiredNameDestinationError)
			}
			g.skip = true
			return nil
		},
	}
}

func (g *GolimCLI) addRemoveLimiterCMD() *ff.Command {
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
					operation: main.removeLimiterOperation,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *GolimCLI) createAddCMD() *ff.Command {
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
					operation:    main.addRoleOperation,
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

func (g *GolimCLI) createRemoveCMD() *ff.Command {
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
					operation: main.removeRoleOperation,
					limiterID: *roleID,
				}
			}
			g.skip = true
			return nil
		},
	}
}

func (g *GolimCLI) createGetRolesCMD() *ff.Command {
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
					operation: main.getRolesOperation,
					limiterID: *limiterID,
				}
			} else {
				return errors.New(main.requiredLimiterIDError)
			}
			g.skip = true
			return nil
		},
	}
}

func toCreateRoleParam(g *GolimCLI) role2.CreateRoleParams {
	return role2.CreateRoleParams{
		Endpoint:       g.limiterRole.endPoint,
		Operation:      g.limiterRole.method,
		BucketSize:     int64(g.limiterRole.bucketSize),
		AddTokenPerMin: g.limiterRole.addToken,
		InitialTokens:  int64(g.limiterRole.initialToken),
		RateLimiterID:  int64(g.limiterRole.limiterID),
	}
}

func toCreateRateLimiter(g *GolimCLI) role2.CrateRateLimiterParams {
	return role2.CrateRateLimiterParams{
		Name:        g.limiter.name,
		Destination: g.limiter.destination,
	}
}
