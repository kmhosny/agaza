package connectionManager

const redisNameSeparator = ":"

//redis keys and relations names
const (
	servicePrefix            = "agaza"
	userKeyPrefix            = servicePrefix + redisNameSeparator + "user"   //redis hash
	usersListName            = servicePrefix + redisNameSeparator + "users"  //redis list
	leavesListName           = servicePrefix + redisNameSeparator + "leaves" // redis list
	leaveKeyPrefix           = servicePrefix + redisNameSeparator + "leave"  //redis hash
	leaveIDCounterPrefix     = servicePrefix + redisNameSeparator + "last_leave_id"
	daySortedSetPrefix       = servicePrefix + redisNameSeparator + "day" + redisNameSeparator
	userTakenLeavesSetPrefix = servicePrefix + redisNameSeparator + "user" + redisNameSeparator
	userIDCounterPrefix      = servicePrefix + redisNameSeparator + "last_user_id"
)

//representing field values for application model
const (
	leaveID           = "ID"
	leaveUserID       = "user_id"
	leaveReason       = "reason"
	leaveFrom         = "from"
	leaveTo           = "to"
	leaveStatus       = "status"
	leaveDepartmentID = "department_id"
	leaveUserName     = "user_name"
	leaveType         = "type"
)

//representing field values for user models
const (
	userID                    = "ID"
	userName                  = "name"
	userDepartmentID          = "department_id"
	userRemainingAnnualLeaves = "remaining_annual_leave"
)
