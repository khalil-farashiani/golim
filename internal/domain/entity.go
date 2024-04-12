package domain

import (
	"github.com/khalil-farashiani/golim/internal/store"
	"github.com/khalil-farashiani/golim/pkg/log"
)

type LimiterRole struct {
	Operation    string
	LimiterID    int
	EndPoint     string
	Method       string
	BucketSize   int
	InitialToken int
	AddToken     int64
}

type Limiter struct {
	ID          interface{}
	Name        string
	Destination string
	Operation   string
}

type Golim struct {
	Limiter     *Limiter
	LimiterRole *LimiterRole
	Port        int64
	Skip        bool
	*log.Logger
	store.Store
}
