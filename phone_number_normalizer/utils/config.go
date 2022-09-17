package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost             string `mapstructure:"DB_HOST"`
	DBPort             int    `mapstructure:"DB_PORT"`
	DBUsername         string `mapstructure:"DB_USERNAME"`
	DBPassword         string `mapstructure:"DB_PASSWORD"`
	DBDatabaseNameRaw  string `mapstructure:"DB_DATABASE_NAME_RAW"`
	DBDatabaseNameSqlx string `mapstructure:"DB_DATABASE_NAME_SQLX"`
	DBDatabaseNameGorm string `mapstructure:"DB_DATABASE_NAME_GORM"`
}

type DatabaseNameType int

const (
	RawDB DatabaseNameType = iota
	SqlxDB
	GormDB
)


// GetPGConnectionString: returns a postgres connection string withouth the database name
func (c *Config) GetPGConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUsername, c.DBPassword)
}


// GetDBConnectionString: returns a postgres connection string with the database name
func (c *Config) GetDBConnectionString(dbType DatabaseNameType) string {

	var dbName string
	switch dbType {
	case RawDB:
		dbName = c.DBDatabaseNameRaw
	case SqlxDB:
		dbName = c.DBDatabaseNameSqlx
	case GormDB:
		dbName = c.DBDatabaseNameGorm
	default:
		panic("invalid database type")
	}

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, dbName)
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
