package database

import "context"

type Pinger interface {
	Ping() error
	PingContext(ctx context.Context) error
}
