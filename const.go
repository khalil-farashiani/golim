package main

const (
	createLimiterOperation = "create"
	removeLimiterOperation = "remove"
)

const (
	addRoleOperation      = "add"
	removeRoleOperationID = "remove"
	getRolesOperationID   = "getRoles"
)

const (
	OperationGet    = "GET"
	OperationPost   = "POST"
	OperationPut    = "PUT"
	OperationPatch  = "PATCH"
	OperationDelete = "DELETE"
)

const (
	unknownLimiterRoleError      = "unknown limiter role operation"
	unknownLimiterError          = "unknown limiter operation"
	unsupportedOperationError    = "unsupported operation"
	requiredNameDestinationError = "name and destination is required"
	requiredLimiterIDError       = "limiter id is required"
	createProxyError             = "Error creating proxy request"
	sendingProxyError            = "Error sending proxy request"
	slowDownError                = "slow down"
)
