package connectionManager

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"agaza/config"
	"agaza/logger"
	"agaza/models"

	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

//RedisConnection is a struct wrapper to the pool and singleton object of redis connector
type RedisConnection struct {
	p    *pool.Pool
	once sync.Once
}

/* getConnectionPool creating a connection to redis
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

/*GetOneConnection this method is used to get onc connection from the pool of connections
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

//PutBackConnection this method is used to put back a connection to the pool
//remember to defer PutBackConnection whenever you cann GetOneConnection
func (r *RedisConnection) PutBackConnection(con *redis.Client) {
	r.p.Put(con)
}

//GetLeaveByID get the leave given by this ID as leave model
func (r *RedisConnection) GetLeaveByID(dayLeaveID string) (*models.Leave, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return nil, err
	}
	defer r.PutBackConnection(conn)

	result, err := conn.Cmd("HGETALL", leaveKeyPrefix+dayLeaveID).List()
	if err != nil {
		logger.Error.Println("Failed execute command HGETALL", leaveKeyPrefix+dayLeaveID)
		return nil, err
	}
	if len(result) == 0 {
		logger.Error.Println(leaveKeyPrefix, dayLeaveID, " doesn't exist")
		return nil, errors.New(leaveKeyPrefix + dayLeaveID + " doesn't exist")
	}
	leaveObject := new(models.Leave)
	for i, value := range result {
		switch value {
		case leaveID:
			leaveObject.ID = result[i+1]
			break
		case leaveDepartmentID:
			leaveObject.DepartmentID = result[i+1]
			break
		case leaveFrom:
			const shortForm = "2006-Jan-02"
			t, _ := time.Parse(shortForm, result[i+1])
			leaveObject.From = t
			break
		case leaveTo:
			const shortForm = "2006-Jan-02"
			t, _ := time.Parse(shortForm, result[i+1])
			leaveObject.To = t
			break
		case leaveReason:
			leaveObject.Reason = result[i+1]
			break
		case leaveStatus:
			leaveObject.Status = result[i+1]
			break
		case leaveUserName:
			leaveObject.UserName = result[i+1]
			break
		case leaveType:
			leaveObject.Type, _ = strconv.Atoi(result[i+1])
			break
		}

	}
	return leaveObject, nil
}

//GetLeavesInRange returns all leaves as ExposedLeave object in range of from to to
func (r *RedisConnection) GetLeavesInRange(from time.Time, to time.Time) ([]*models.ExposedLeave, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return nil, err
	}
	defer r.PutBackConnection(conn)

	leaveIds := make([]string, 0)

	for d := from; d.Before(to) || d == to; d = d.Add(time.Duration(24) * time.Hour) {
		datePart := strings.TrimSuffix(d.String(), " 00:00:00 +0000 UTC")

		result, err := conn.Cmd("ZRANGE", daySortedSetPrefix+datePart, 0, -1).List()
		if err != nil {
			logger.Error.Println("Error Queuing command ZRANGE", daySortedSetPrefix+datePart, 0, -1)
			return nil, err
		}
		for _, v := range result {
			leaveIds = append(leaveIds, v)
		}
	}
	leaves := make([]*models.ExposedLeave, 0)
	for _, dayLeaveID := range leaveIds {
		result, err := conn.Cmd("HGETALL", dayLeaveID, 0, -1).List()
		if err != nil {
			logger.Error.Println("Failed execute command HGETALL", dayLeaveID)
			return nil, err
		}
		if len(result) == 0 {
			logger.Error.Println(dayLeaveID, " doesn't exist")
			return nil, errors.New(dayLeaveID + " doesn't exist")
		}
		leaveObject := new(models.ExposedLeave)
		for i, value := range result {
			switch value {
			case leaveID:
				leaveObject.ID = result[i+1]
				break
			case leaveDepartmentID:
				leaveObject.DepartmentID = result[i+1]
				break
			case leaveFrom:
				const shortForm = "2006-Jan-02"
				t, _ := time.Parse(shortForm, result[i+1])
				leaveObject.From = t
				break
			case leaveTo:
				const shortForm = "2006-Jan-02"
				t, _ := time.Parse(shortForm, result[i+1])
				leaveObject.To = t
				break
			case leaveUserName:
				leaveObject.UserName = result[i+1]
				break
			}
		}
		leaves = append(leaves, leaveObject)
	}
	return leaves, nil
}

//GetUserByID return the user of this ID
func (r *RedisConnection) GetUserByID(ID string) (*models.User, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return nil, err
	}
	defer r.PutBackConnection(conn)

	result, err := conn.Cmd("HGETALL", userKeyPrefix+ID).List()
	if err != nil {
		logger.Error.Println("Failed execute command HGETALL", userKeyPrefix+ID)
		return nil, err
	}
	if len(result) == 0 {
		logger.Error.Println(userKeyPrefix+ID, " doesn't exist")
		return nil, errors.New(userKeyPrefix + ID + " doesn't exist")
	}
	userObject := new(models.User)
	for i, value := range result {
		switch value {
		case userID:
			userObject.ID = result[i+1]
			break
		case userDepartmentID:
			userObject.DepartmentID = result[i+1]
			break
		case userName:
			userObject.Name = result[i+1]
			break
		case userRemainingAnnualLeaves:
			userObject.RemainingAnnualLeaves, _ = strconv.Atoi(result[i+1])
			break
		}

	}
	return userObject, nil
}
