package static

const (
	EnvLocal       = "local"
	EnvDev         = "dev"
	EnvEnvironment = "ENV"
)

// server environment variable name
const (
	EnvServerAddress    = "SERVER_ADDRESS"
	EnvServerPort       = "SERVER_PORT"
	EnvServerTimeout    = "SERVER_TIMEOUT"
	EnvServerApiPath    = "SERVER_API_PATH"
	EnvServerSwaggerUrl = "SERVER_SWAGGER_URL"
)

// database environment variable name
const (
	EnvDbHost      = "DB_HOST"
	EnvDbPort      = "DB_PORT"
	EnvDbUser      = "DB_USER"
	EnvDbPassword  = "DB_PASSWORD"
	EnvDbName      = "DB_NAME"
	EnvDbParameter = "DB_PARAMETER"
)

// auth environment variable name
const (
	EnvAuthSubject      = "AUTH_SUBJECT"
	EnvAuthIssuer       = "AUTH_ISSUER"
	EnvAuthAudience     = "AUTH_AUDIENCE"
	EnvAuthAccessSecret = "AUTH_ACCESS_SECRET"
	EnvAuthAccessTime   = "AUTH_ACCESS_TIME"
)
