package leaderboard

import (
	"github.com/redis/go-redis/v9"
)

// A struct to hold dependencies for your HTTP handler
type API struct {
	RedisClient *redis.Client
	// other dependencies like database connections, loggers, etc.
}

// NewAPI creates a new API instance with its dependencies
func NewAPI(rdb *redis.Client) *API {
	return &API{
		RedisClient: rdb,
	}
}
