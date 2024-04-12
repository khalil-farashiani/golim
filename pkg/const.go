package pkg

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
	unknownLimiterRoleError      = "unknown limiter role operation"
	unknownLimiterError          = "unknown limiter operation"
	requiredNameDestinationError = "name and destination is required"
	requiredLimiterIDError       = "limiter id is required"
	slowDownError                = "slow down"
	notFoundSqlError             = "sql: no rows in result set"
)
