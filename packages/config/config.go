package config

const (
	SessionTimeDefault = 300
	SessionTimeExpireTime = -1
)

type TasksConfig struct {
	Route string
	TitleParseNaming string
	HTMLPageName string
	RedirectPath string
}

type ParseKeys struct {
	UsernameParseKey string
	PasswordParseKey string
	RePasswordParseKey string
}

type UserRouteConfig struct {
	GetTask TasksConfig
	DeleteTask TasksConfig
	CreateTask TasksConfig
	Route string
}

type AuthPageConfig struct {
	PageName       string
	Path           string
	RedirectPath   string
	EmptyPathString string
	ParseKeys
	SessionTime int // 300 default
	SessionTimeOut int
}

type Cookie struct {
	Naming string
	UserInfoKey string
}

func NewCookie (name string, key string) Cookie {
	return Cookie{name, key};
}

func NewParseKeys () ParseKeys {
	return ParseKeys{
		UsernameParseKey: "uname",
		PasswordParseKey: "pword",
		RePasswordParseKey: "re-pword",
	}
}

type RouteConfig struct {
	UserConfig UserRouteConfig
	MainLoginConfig AuthPageConfig
	MainRegisterConfig AuthPageConfig
	MainLogoutConfig AuthPageConfig
	Authentication AuthPageConfig
	Cookie
}

var Routes *RouteConfig = &RouteConfig{
	UserConfig: UserRouteConfig{
		GetTask: TasksConfig{
			Route: "/user/tasks",
			HTMLPageName: "todoMain.html",
			RedirectPath: "/user/logout",
		},

		DeleteTask: TasksConfig{
			Route: "/user/deleteTask",
			RedirectPath: "/user/tasks",
		},
		
		CreateTask: TasksConfig{
			Route: "/user/addTask",
			RedirectPath: "/user/tasks",
		},

		Route: "/user",
	},

	MainLoginConfig: AuthPageConfig{
		PageName: "login.html",
		Path: "/login",
		EmptyPathString: "/",
		ParseKeys: NewParseKeys(),
		RedirectPath: "/user/tasks",
	},

	MainRegisterConfig: AuthPageConfig{
		PageName: "register.html",
		Path: "/register",
		ParseKeys: NewParseKeys(),
	},

	MainLogoutConfig: AuthPageConfig{
		RedirectPath: "/login",
	},

	Authentication: AuthPageConfig{
		RedirectPath: "/logout",
		SessionTime: SessionTimeDefault,
		SessionTimeOut: SessionTimeExpireTime,
		ParseKeys: NewParseKeys(),
	},

	Cookie: NewCookie("loginSession", "user"),
}

// hard coded part should be improved by more readible coding