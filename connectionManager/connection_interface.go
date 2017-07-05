package connectionManager

import (
	"agaza/models"
	"sync"
)

type DBoperationsFactory interface {

	//users model
	NewUser(*models.User) (int, error)
	DeleteUser(string) (int, error)
	UserExists(userId string) (bool, error)

	//application model
	NewLeave(*models.Leave) (int, error)
	GetLeaveById(string) (*models.Leave, error)
	DeleteLeave(*models.Leave) (int, error)
	GetAllLeaves() ([]string, error)
}

var (
	redisConnection *RedisConnection
	once            sync.Once
)

func GetRedisConnection() DBoperationsFactory {
	once.Do(func() {
		redisConnection = new(RedisConnection)
	})
	return redisConnection
}
