package connectionManager

import (
	"strconv"
	"sync"

	"agaza/config"
	"agaza/logger"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type RedisConnection struct {
	p    *pool.Pool
	once sync.Once
}

/* creating a connection to redis
** this connection will be shared with all running threads
** a singleton function to avoid having multiple pools
 */
func (r *RedisConnection) getConnectionPool() (*pool.Pool, error) {
	configuration := config.LoadConfiguration()
	r.once.Do(func() {
		tempP, err := pool.New(configuration.Redisprotocol, configuration.Redishost+":"+strconv.Itoa(configuration.Redisport), configuration.Redispoolsize)
		r.p = tempP
		if err != nil {
			logger.Error.Println("error constructing a connection")
		}
	})
	return r.p, nil
}

/* this method is used to get onc connection from the pool of connections
** it must be followed by defer r.PutBackConnection to avoid abusing the pool
** returns: pointer to redis.Client used to communicate with redis db
 */
func (r *RedisConnection) GetOneConnection() (*redis.Client, error) {
	configuration := config.LoadConfiguration()
	_, err := r.getConnectionPool()
	if err != nil {
		logger.Error.Println("Failed to retrieve connection pool")
		return nil, err
	}
	con, err := r.p.Get()
	if err != nil {
		logger.Error.Println("Failed to retrieve connection")
		return nil, err
	}
	if configuration.Environment == "testing" {
		con.Cmd("SELECT", configuration.Dbnumber)
	}
	return con, nil
}

//this method is used to put back a connection to the pool
//remember to defer PutBackConnection whenever you cann GetOneConnection
func (r *RedisConnection) PutBackConnection(con *redis.Client) {
	r.p.Put(con)
}
