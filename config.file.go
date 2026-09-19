package blueprints

var ConfigViper = `
package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func envPrefix(moduleName string) string {
	parts := strings.Split(moduleName, "/")
	prefix := parts[len(parts)-1]
	prefix = strings.NewReplacer("-", "_", ".", "_").Replace(prefix)
	prefix = strings.ToUpper(strings.TrimSpace(prefix))
	if prefix == "" {
		return "APP"
	}

	return prefix
}

func InitViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	for _, configPath := range []string{".", "./configs", "../", "../configs"} {
		viper.AddConfigPath(configPath)
	}
	viper.SetEnvPrefix(envPrefix("{{ .ModuleName }}"))
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("cannot read config.yml from project root or configs directory: %w", err)
	}

	return nil
}
`

var ConfigYml = `
APP_NAME: "{{ .ModuleName }}"
APP_ENV: "development"
APP_PORT: 3000
LOG_DIR: "logs"
HTTP_READ_TIMEOUT: "15s"
HTTP_WRITE_TIMEOUT: "15s"
HTTP_IDLE_TIMEOUT: "60s"
HTTP_BODY_LIMIT: 12582912
SHUTDOWN_TIMEOUT: "10s"

DB_DRIVER: "{{ .DBName }}"
DB_HOST: "{{ .DBHost }}"
DB_PORT: {{ .DBPort }}
DB_NAME: "{{ .DatabaseName }}"
DB_USERNAME: "{{ .DBUser }}"
DB_PASSWORD: ""
DB_SSL_MODE: "disable"
DB_SLOW_QUERY_THRESHOLD: "200ms"
`
