package config

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config contains application configuration.
type Config struct {
	AppEnv   string
	LogLevel string

	DepartmentServicePort int
	EmployeeServicePort   int
	AuthServicePort       int
	APIGatewayPort        int

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBSSLMode  string

	DepartmentDBName string
	EmployeeDBName   string

	ConsulAddress    string
	ConsulScheme     string
	ConsulDatacenter string

	DepartmentServiceURL string
	EmployeeServiceURL   string
	AuthServiceURL       string

	HTTPClientTimeoutSeconds   int
	HTTPClientMaxRetries       int
	HTTPClientRetryDelayMillis int

	CircuitBreakerFailureThreshold    int
	CircuitBreakerResetTimeoutSeconds int

	JWTSecret            string
	JWTExpirationMinutes int

	AuthAdminUsername    string
	AuthAdminPassword    string
	AuthEmployeeUsername string
	AuthEmployeePassword string
}

// Load reads configuration from the .env file.
func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		AppEnv:   getEnv("APP_ENV"),
		LogLevel: strings.ToUpper(getEnv("LOG_LEVEL")),

		DepartmentServicePort: getIntEnvWithDefault("DEPARTMENT_SERVICE_PORT", 8081),
		EmployeeServicePort:   getIntEnvWithDefault("EMPLOYEE_SERVICE_PORT", 8082),
		AuthServicePort:       getIntEnvWithDefault("AUTH_SERVICE_PORT", 8083),
		APIGatewayPort:        getIntEnvWithDefault("API_GATEWAY_PORT", 8080),

		DBHost:     getEnv("DB_HOST"),
		DBPort:     getIntEnvWithDefault("DB_PORT", 5432),
		DBUser:     getEnv("DB_USER"),
		DBPassword: getEnv("DB_PASSWORD"),
		DBSSLMode:  getEnv("DB_SSLMODE"),

		DepartmentDBName: getEnv("DEPARTMENT_DB_NAME"),
		EmployeeDBName:   getEnv("EMPLOYEE_DB_NAME"),

		ConsulAddress:    getEnv("CONSUL_ADDRESS"),
		ConsulScheme:     getEnv("CONSUL_SCHEME"),
		ConsulDatacenter: getEnv("CONSUL_DATACENTER"),

		DepartmentServiceURL: getEnv("DEPARTMENT_SERVICE_URL"),
		EmployeeServiceURL:   getEnv("EMPLOYEE_SERVICE_URL"),
		AuthServiceURL:       getEnv("AUTH_SERVICE_URL"),

		HTTPClientTimeoutSeconds:   getIntEnv("HTTP_CLIENT_TIMEOUT_SECONDS"),
		HTTPClientMaxRetries:       getIntEnv("HTTP_CLIENT_MAX_RETRIES"),
		HTTPClientRetryDelayMillis: getIntEnv("HTTP_CLIENT_RETRY_DELAY_MILLIS"),

		CircuitBreakerFailureThreshold:    getIntEnv("CIRCUIT_BREAKER_FAILURE_THRESHOLD"),
		CircuitBreakerResetTimeoutSeconds: getIntEnv("CIRCUIT_BREAKER_RESET_TIMEOUT_SECONDS"),

		JWTSecret:            getEnv("JWT_SECRET"),
		JWTExpirationMinutes: getIntEnv("JWT_EXPIRATION_MINUTES"),

		AuthAdminUsername:    getEnv("AUTH_ADMIN_USERNAME"),
		AuthAdminPassword:    getEnv("AUTH_ADMIN_PASSWORD"),
		AuthEmployeeUsername: getEnv("AUTH_EMPLOYEE_USERNAME"),
		AuthEmployeePassword: getEnv("AUTH_EMPLOYEE_PASSWORD"),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate checks required configuration values.
func (c Config) Validate() error {
	required := map[string]string{
		"APP_ENV":                c.AppEnv,
		"LOG_LEVEL":              c.LogLevel,
		"DB_HOST":                c.DBHost,
		"DB_USER":                c.DBUser,
		"DB_PASSWORD":            c.DBPassword,
		"DB_SSLMODE":             c.DBSSLMode,
		"DEPARTMENT_DB_NAME":     c.DepartmentDBName,
		"EMPLOYEE_DB_NAME":       c.EmployeeDBName,
		"CONSUL_ADDRESS":         c.ConsulAddress,
		"CONSUL_SCHEME":          c.ConsulScheme,
		"CONSUL_DATACENTER":      c.ConsulDatacenter,
		"DEPARTMENT_SERVICE_URL": c.DepartmentServiceURL,
		"EMPLOYEE_SERVICE_URL":   c.EmployeeServiceURL,
		"AUTH_SERVICE_URL":       c.AuthServiceURL,
		"JWT_SECRET":             c.JWTSecret,
		"AUTH_ADMIN_USERNAME":    c.AuthAdminUsername,
		"AUTH_ADMIN_PASSWORD":    c.AuthAdminPassword,
		"AUTH_EMPLOYEE_USERNAME": c.AuthEmployeeUsername,
		"AUTH_EMPLOYEE_PASSWORD": c.AuthEmployeePassword,
	}

	for key, value := range required {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("missing required configuration: %s", key)
		}
	}

	if c.DBPort <= 0 || c.DBPort > 65535 {
		return fmt.Errorf("invalid DB_PORT")
	}

	ports := []int{
		c.DepartmentServicePort,
		c.EmployeeServicePort,
		c.AuthServicePort,
		c.APIGatewayPort,
	}

	for _, port := range ports {
		if port <= 0 || port > 65535 {
			return fmt.Errorf("invalid service port: %d", port)
		}
	}

	if c.HTTPClientTimeoutSeconds <= 0 ||
		c.HTTPClientRetryDelayMillis <= 0 ||
		c.HTTPClientMaxRetries < 0 {
		return fmt.Errorf("invalid HTTP client configuration")
	}

	if c.CircuitBreakerFailureThreshold <= 0 ||
		c.CircuitBreakerResetTimeoutSeconds <= 0 {
		return fmt.Errorf("invalid circuit breaker configuration")
	}

	if c.JWTExpirationMinutes <= 0 {
		return fmt.Errorf("JWT expiration must be positive")
	}

	if len(c.JWTSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must contain at least 32 characters")
	}

	switch c.LogLevel {
	case "DEBUG", "INFO", "WARN", "ERROR":
	default:
		return fmt.Errorf("invalid LOG_LEVEL: %s", c.LogLevel)
	}

	if c.ConsulScheme != "http" && c.ConsulScheme != "https" {
		return fmt.Errorf("CONSUL_SCHEME must be http or https")
	}

	return nil
}

// DatabaseURL builds a PostgreSQL connection URL.
func (c Config) DatabaseURL(databaseName string) string {
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, strconv.Itoa(c.DBPort)),
		Path:   databaseName,
	}

	query := connectionURL.Query()
	query.Set("sslmode", c.DBSSLMode)
	connectionURL.RawQuery = query.Encode()

	return connectionURL.String()
}

// HTTPTimeout returns the configured HTTP timeout.
func (c Config) HTTPTimeout() time.Duration {
	return time.Duration(c.HTTPClientTimeoutSeconds) * time.Second
}

// HTTPRetryDelay returns the configured HTTP retry delay.
func (c Config) HTTPRetryDelay() time.Duration {
	return time.Duration(c.HTTPClientRetryDelayMillis) * time.Millisecond
}

// CircuitBreakerResetTimeout returns the configured reset duration.
func (c Config) CircuitBreakerResetTimeout() time.Duration {
	return time.Duration(c.CircuitBreakerResetTimeoutSeconds) * time.Second
}

// getEnv reads a string environment variable.
func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

// getIntEnv reads an integer environment variable.
func getIntEnv(key string) int {
	value, _ := strconv.Atoi(getEnv(key))
	return value
}

// getIntEnvWithDefault reads an integer or returns its default value.
func getIntEnvWithDefault(key string, defaultValue int) int {
	value := getEnv(key)
	if value == "" {
		return defaultValue
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}

	return parsedValue
}
