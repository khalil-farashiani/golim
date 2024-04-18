package contract

import (
	"context"
)

type Job interface {
	RunTasks(ctx context.Context, cmd func()) error
}
