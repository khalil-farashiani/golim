package entity

type Role struct {
	Operation    string
	LimiterID    int
	EndPoint     string
	Method       string
	BucketSize   int
	InitialToken int
	AddToken     int64
}
