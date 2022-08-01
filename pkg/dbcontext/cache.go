package dbcontext

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
)

// Cache represents a Cache connection that can be used to run SQL queries.
type Cache struct {
	cache *redis.Client
}

// NewCache returns a new Cache connection that wraps the given redis.client instance.
func NewCache(cache *redis.Client) *Cache {
	return &Cache{cache}
}

// Cache returns the redis.Client wrapped by this object.
func (cache *Cache) Cache() *redis.Client {
	return cache.cache
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a Cache connection associated with the context.
func (cache *Cache) With(c *gin.Context) *redis.Client {
	return cache.cache.WithContext(c)
}

func ConnectToCache(dsn string) (*redis.Client) {
	rdb := redis.NewClient(&redis.Options{
        Addr:     dsn,
        Password: "", // no password set
        DB:       0,  // use default DB
    })
	return rdb
}