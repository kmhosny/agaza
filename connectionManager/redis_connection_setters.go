package connectionManager

import (
	"agaza/logger"
	"agaza/models"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mediocregopher/radix.v2/redis"
)

const queuedKeyword = "queued"

/*steps for adding a data layer function:
* 1- retrieve a connection from the connections pool
* 2- defer the connection to make sure it's released at the end of the function
* 3- write the data retrieval login
 */

//NewLeave is to create a new leave, it takes a leave object, then inserts it into
//leaves, user's leaves and the days in the leaves starting from => to
func (r *RedisConnection) NewLeave(l *models.Leave) (string, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return "-1", err
	}
	defer r.PutBackConnection(conn)

	strID, errIncr := conn.Cmd("INCR", leaveIDCounterPrefix).Int()
	newLeaveID := strconv.Itoa(strID)
	if errIncr != nil {
		logger.Error.Println("Failed to incremenet number of leaves")
		return "-1", errIncr
	}

	remainStr, err := conn.Cmd("HGET", userKeyPrefix+l.UserID, userRemainingAnnualLeaves).Str()
	if err != nil {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", err
	}

	remainingLeaves, _ := strconv.Atoi(remainStr)

	logger.Trace.Println(remainingLeaves, remainStr, userKeyPrefix+l.UserID)

	userName, err := conn.Cmd("HGET", userKeyPrefix+l.UserID, userName).Str()
	if err != nil {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", err
	}

	days := int((l.To.Sub(l.From) / 24).Hours()) + 1

	logger.Trace.Println(days)
	if remainingLeaves-days < 0 {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", errors.New("Your remaining leaves are less than days requested")
	}

	return r.createUpdateLeave(l, newLeaveID, days, userName, remainingLeaves, conn)
}

//EditLeave is to create a new leave, it takes a leave object, then inserts it into
//leaves, user's leaves and the days in the leaves starting from => to
func (r *RedisConnection) EditLeave(l *models.Leave) (string, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return "-1", err
	}
	defer r.PutBackConnection(conn)

	remainStr, err := conn.Cmd("HGET", userKeyPrefix+l.UserID, userRemainingAnnualLeaves).Str()
	if err != nil {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", err
	}

	remainingLeaves, _ := strconv.Atoi(remainStr)

	userName, err := conn.Cmd("HGET", userKeyPrefix+l.UserID, userName).Str()
	if err != nil {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", err
	}

	days := int((l.To.Sub(l.From) / 24).Hours()) + 1

	if remainingLeaves-days < 0 {
		logger.Error.Println("Failed to get remaining leaves")
		return "-1", errors.New("Your remaining leaves are less than days requested")
	}

	return r.createUpdateLeave(l, l.ID, days, userName, remainingLeaves, conn)
}

//TODO: bug, it inserts the same date in the from to date fields
func (r *RedisConnection) createUpdateLeave(l *models.Leave, newLeaveID string, days int, userName string, remainingLeaves int, conn *redis.Client) (string, error) {

	if ok, err1 := conn.Cmd("MULTI").Str(); strings.ToLower(ok) != "ok" {
		logger.Error.Println("Cannot execute commands now")
		return "-1", err1
	}

	if queued, err2 := conn.Cmd("SADD", leavesListName, leaveKeyPrefix+newLeaveID).Str(); strings.ToLower(queued) != queuedKeyword {
		logger.Error.Println("Error Queuing command SADD", leavesListName, leaveKeyPrefix+newLeaveID)
		return "-1", err2
	}

	logger.Trace.Println(l.DepartmentID)
	if queued, err3 := conn.Cmd("HMSET", leaveKeyPrefix+newLeaveID, leaveID, newLeaveID, leaveUserID, l.UserID, leaveDepartmentID, l.DepartmentID,
		leaveFrom, l.From.Format("2006-Jan-01"), leaveTo, l.To.Format("2006-Jan-01"), leaveType, l.Type, leaveReason, l.Reason, leaveStatus, l.Status, leaveUserName, userName).Str(); strings.ToLower(queued) != queuedKeyword {
		logger.Error.Println("Error Queuing command HMSET", leaveKeyPrefix+newLeaveID, leaveID, newLeaveID, leaveUserID, l.UserID, leaveDepartmentID, l.DepartmentID,
			leaveFrom, l.From, leaveTo, l.To, leaveType, l.Type, leaveReason, l.Reason, leaveStatus, l.Status)
		return "-1", err3
	}

	if queued, err4 := conn.Cmd("SADD", userTakenLeavesSetPrefix+l.UserID+redisNameSeparator+"leaves", newLeaveID).Str(); strings.ToLower(queued) != queuedKeyword {
		logger.Error.Println("Error Queuing command SADD", userTakenLeavesSetPrefix+l.UserID+redisNameSeparator+"leaves", newLeaveID)
		return "-1", err4
	}

	for d := l.From; d.Before(l.To) || d == l.To; d = d.Add(time.Duration(24) * time.Hour) {
		datePart := strings.TrimSuffix(d.String(), " 00:00:00 +0000 UTC")
		if queued, err4 := conn.Cmd("ZADD", daySortedSetPrefix+datePart, newLeaveID, leaveKeyPrefix+newLeaveID).Str(); strings.ToLower(queued) != queuedKeyword {
			logger.Error.Println("Error Queuing command ZADD", daySortedSetPrefix+datePart, newLeaveID, leaveKeyPrefix+newLeaveID)
			return "-1", err4
		}
	}

	result := conn.Cmd("EXEC")
	if result.Err != nil {
		logger.Trace.Println("Error while executing commands")
		return "-1", result.Err
	}

	remainingLeaves = remainingLeaves - days

	logger.Trace.Println(remainingLeaves)

	_, err := conn.Cmd("HSET", userKeyPrefix+l.UserID, userRemainingAnnualLeaves, remainingLeaves).Int()
	if err != nil {
		logger.Error.Println("Failed to set remaning leaves of user", userKeyPrefix+l.UserID)
		return "-1", err
	}

	return newLeaveID, nil
}

//NewUser create a new user, it takes a user object, generates an ID then insters the new ID to list of users
//adds the user object as hash
func (r *RedisConnection) NewUser(l *models.User) (int, error) {
	conn, err := r.GetOneConnection()
	if err != nil {
		logger.Error.Println("Failed to retrieve a connection")
		return -1, err
	}
	defer r.PutBackConnection(conn)

	newUserID, errIncr := conn.Cmd("INCR", userIDCounterPrefix).Str()
	if errIncr != nil {
		logger.Error.Println("Failed to incremenet number of users")
		return -1, errIncr
	}

	if ok, err1 := conn.Cmd("MULTI").Str(); strings.ToLower(ok) != "ok" {
		logger.Error.Println("Cannot execute commands now")
		return -1, err1
	}

	if queued, err2 := conn.Cmd("SADD", usersListName, userKeyPrefix+newUserID).Str(); strings.ToLower(queued) != queuedKeyword {
		logger.Error.Println("Error Queuing command SADD", usersListName, userKeyPrefix+newUserID)
		return -1, err2
	}

	if queued, err3 := conn.Cmd("HMSET", userKeyPrefix+newUserID, userID, newUserID, userName, l.Name, userDepartmentID, l.DepartmentID,
		userRemainingAnnualLeaves, l.RemainingAnnualLeaves).Str(); strings.ToLower(queued) != queuedKeyword {
		logger.Error.Println("Error Queuing command HMSET", userKeyPrefix+newUserID, userID, newUserID, userName, l.Name, userDepartmentID, l.DepartmentID,
			userRemainingAnnualLeaves, l.RemainingAnnualLeaves)
		return -1, err3
	}

	result := conn.Cmd("EXEC")
	if result.Err != nil {
		logger.Trace.Println("Error while executing commands")
		return -1, result.Err
	}
	return 1, nil
}
