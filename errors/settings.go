package errors

type Settings struct {
	UserInfoContextKey any

	SystemInfo any
}

var settings = &Settings{
	UserInfoContextKey: nil,
	SystemInfo:         nil,
}

func Init(
	systemInfo any,
	userInfoContextKey any,
) {
	settings.UserInfoContextKey = userInfoContextKey
	settings.SystemInfo = systemInfo
}
