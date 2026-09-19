package blueprints

var SQLInterface = `
	package databases

	import (
		"fmt"
		"github.com/spf13/viper"
		"database/sql"

		_ "{{ .DBUrl }}"
	)
	
	func ConnectDB() *sql.DB {

		dsn := fmt.Sprintf("%s:%s@%s:%s?database=%s",
			viper.Get("DB_USERNAME"),
			viper.Get("DB_PASSWORD"),
			viper.Get("DB_HOST"),
			viper.Get("DB_PORT"),
			viper.Get("DB_NAME"))

		db, err := sql.Open("{{ .DBName }}", dsn)
		if err != nil {
			panic(err)
		}

		return db
	}
`

var ORMInterface = `
	package databases

	import (
		"fmt"
		"github.com/spf13/viper"
		"{{ .ModuleName }}/pkg/logger"

		"gorm.io/driver/{{ .DBName }}"
		"gorm.io/gorm"
	)

	func ConnectDB() *gorm.DB {
		var dsn string
		{{- if eq .DBName "postgres" }}
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			viper.GetString("DB_HOST"), viper.GetString("DB_USERNAME"), viper.GetString("DB_PASSWORD"),
			viper.GetString("DB_NAME"), viper.GetInt("DB_PORT"), viper.GetString("DB_SSL_MODE"))
		{{- else if eq .DBName "mysql" }}
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			viper.GetString("DB_USERNAME"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_HOST"),
			viper.GetInt("DB_PORT"), viper.GetString("DB_NAME"))
		{{- else }}
		dsn = fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			viper.GetString("DB_USERNAME"), viper.GetString("DB_PASSWORD"), viper.GetString("DB_HOST"),
			viper.GetInt("DB_PORT"), viper.GetString("DB_NAME"))
		{{- end }}

		conn, err := gorm.Open({{ .DBName }}.Open(dsn), &gorm.Config{
			TranslateError: false,
			Logger:         logger.NewGORMLogger(viper.GetDuration("DB_SLOW_QUERY_THRESHOLD")),
		})
		if err != nil {
			panic(err)
		}

		return conn
	}
`
