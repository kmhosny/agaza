package connectionManager

import (
	"agaza/models"
	"sync"
	"time"
)

//DBoperationsFactory is the factory class of db operations, all should be implemented by any db engine
//to be used.
type DBoperationsFactory interface {

	//users model
	NewUser(*models.User) (int, error)
	GetUserByID(string) (*models.User, error)
	//DeleteUser(string) (int, error)
	//UserExists(userId string) (bool, error)

	//application model
	NewLeave(*models.Leave) (string, error)
	EditLeave(*models.Leave) (string, error)
	GetLeaveByID(string) (*models.Leave, error)
	//DeleteLeave(*models.Leave) (int, error)
	GetLeavesInRange(time.Time, time.Time) ([]*models.ExposedLeave, error)
}

var (
	redisConnection *RedisConnection
	once            sync.Once
)

//GetRedisConnection get a db connection that connects to Redis
func GetRedisConnection() DBoperationsFactory {
	once.Do(func() {
		redisConnection = new(RedisConnection)
	})
	return redisConnection
}
