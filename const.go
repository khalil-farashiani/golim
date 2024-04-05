package main

const (
	createLimiterOperation = "create"
	removeLimiterOperation = "remove"
)

const (
	addRoleOperation    = "add"
	removeRoleOperation = "remove"
	getRolesOperation   = "getRoles"
)

const (
	initialTokenForTheFirstUSerRequest = 1
)

const (
	limiterCacheMainKey         = "GOLIM_KEY"
	limiterCacheRegexPatternKey = "*GOLIM_KEY"
)

const (
	unknownLimiterRoleError      = "unknown limiter role operation"
	unknownLimiterError          = "unknown limiter operation"
	requiredNameDestinationError = "name and destination is required"
	requiredLimiterIDError       = "limiter id is required"
	slowDownError                = "slow down"
	notFoundSqlError             = "sql: no rows in result set"
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
	basePath       = "/.golim/"
	loggerFileName = "/logs/golim.log"
	dbPath         = "/db/"
)
