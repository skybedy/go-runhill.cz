package config

const (
	HTTP_PORT   = "1305"
	DbDriver    = "mysql"
	DbUser      = "skybedy"
	DbPass      = "mk1313life"
	DbName      = "runhill.cz"
	ViperConfig = "./config"
)

// Configurations exported
type Configurations struct {
	Server         ServerConfigurations
	Database       DatabaseConfigurations
	Google         GoogleConfigurations
	Facebook       FacebookConfigurations
	Authentication AuthenticationConfigurations
	Email          EmailConfigurations
	EXAMPLE_PATH   string
	EXAMPLE_VAR    string
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port         string
	Name         string
	Webname      string
	StaticPath   string
	TemplatePath string
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBDriver   string
}

type GoogleConfigurations struct {
	ClientId      string
	ClientSecret  string
	RedirectRoute string
}

type FacebookConfigurations struct {
	ClientId      string
	ClientSecret  string
	RedirectRoute string
}

type AuthenticationConfigurations struct {
	SessionName string
}

type EmailConfigurations struct {
	SmtpServer    string
	SmtpPort      int
	EmailCharset  string
	EmailFrom     string
	EmailFromName string
}
