package connectionManager

//common constants between models and keys names
const (
	modified_at       = "Modified_at"
	created_at        = "Created_at"
	languageKeyPrefix = "language"
	app_id            = "App_id"
	app_version       = "App_version"
	app_name          = "App_name"
)

//redis keys and relations names
const (
	servicePrefix                 = "notification"
	outgoingNotificationKeyPrefix = servicePrefix + ":outgoing" //redis queue RPUSH
	requestKeyPrefix              = servicePrefix + ":request"
	requestsHashName              = servicePrefix + ":requests" //redis hash
	scheduleKeyPrefix             = servicePrefix + ":schedule" //redis ordered set
	messageKeyPrefix              = servicePrefix + ":message"
	messagesHashName              = servicePrefix + ":messages"     //redis hash
	userKeyPrefix                 = servicePrefix + ":user"         //redis hash
	usersListName                 = servicePrefix + ":users"        //redis list
	deviceKeyPrefix               = servicePrefix + ":device"       //redis hash
	devicesListName               = servicePrefix + ":devices"      //redis list
	device_countKeyPrefix         = servicePrefix + ":device_count" //redis counter
	appVersionKeyPrefix           = "version"
	channelKeyPrefix              = servicePrefix + ":channel"      //redis hash
	channelsListName              = servicePrefix + ":channels"     //redis list
	subscribersKeyPrefix          = "subscribers"                   //redis list
	subscribedToKeyPrefix         = "subscribedTo"                  //redis list
	applicationsListName          = servicePrefix + ":applications" // redis list
	applicationKeyPrefix          = servicePrefix + ":application"  //redis hash
	versionsKeyPrefix             = "versions"
	packageNameListPrefix         = servicePrefix + ":package"          //redis list
	successDeliveryKeyPrefix      = servicePrefix + ":success_delivery" //redis list
	failDeliveryKeyPrefix         = servicePrefix + ":fail_delivery"    //redis list
	requestAnnouncementKeyPrefix  = servicePrefix + ":request:announcement"
)

const (
	APIKey = "AIzaSyA3Sv-l7bpl1wbpRw1BjB20Fqe2dke9LiE"
)

//representing field values names for user model
const (
	eventtus_uid = "Eventtus_uid"
	timezone     = "Timezone"
	isNotifiable = "isNotifiable"
)

//representing field values for device model
const (
	device_id    = "Device_id"
	user_uid     = "User_uid"
	device_token = "Device_token"
	platform     = "Platform"
	badge        = "Badge"
	metadata     = "Metadata"
	package_name = "package_name"
)

//representing field values for application model
const (
	sender_id        = "Sender_id"
	api_server_key   = "API_server_key"
	certificate_path = "Certificate_path"
)
