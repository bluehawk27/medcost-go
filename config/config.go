package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/joho/godotenv"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/spf13/viper"
)

func getProjectPath() string {
	dir, err := filepath.Abs(filepath.Dir("."))
	if err != nil {
		log.Println("Warning, cannot get current path")
		return ""
	}
	// Traverse back from current directory until service base dir is reached and add to config path
	for !strings.HasSuffix(dir, "medcost-go") && dir != "/" {
		dir, err = filepath.Abs(dir + "/..")
		if err != nil {
			break
		}
	}

	return dir
}

func getHomePath() string {
	home, _ := homedir.Dir()
	return home
}

// Init : Initialize all the things
func Init() {

	// Find home directory.
	viper.SetEnvPrefix("medcost")
	viper.BindEnv("configFile")
	viper.BindEnv("configPath")

	viper.SetDefault("configFile", "config")
	// viper.SetDefault("configPath", "../medcost-go")

	viper.SetDefault("logging.level", "DEBUG")
	viper.SetDefault("logging.errorLogFile", "error.log")

	// Search config in home directory with name ".rs-collabs-brand-test" (without extension).
	viper.AddConfigPath(getHomePath())
	viper.AddConfigPath(getProjectPath())
	viper.AddConfigPath(viper.GetString("configPath"))
	viper.SetConfigName(viper.GetString("configFile"))
	fmt.Println(viper.ConfigFileUsed())
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error using config file:", viper.ConfigFileUsed())
		log.Println(err.Error())
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

// GetDatabaseConfig : Gets the MYSQL Configuration
func GetDatabaseConfig() Database {

	r := Database{
		Host:     viper.GetString("database.host"),
		Port:     viper.GetString("database.port"),
		Password: viper.GetString("database.password"),
		Username: viper.GetString("database.username"),
		DBname:   viper.GetString("database.name"),
	}

	return r
}

// DBConnectionString : Gets the MYSQL Configuration
func DBConnectionString() (string, error) {
	db := GetDatabaseConfig()
	if db.DBname == "" || db.Host == "" || db.Password == "" || db.Username == "" || db.Port == "" {
		return "", errors.New("error reading in from config")
	}

	address := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db.Username, db.Password, db.Host, db.Port, db.DBname)

	return address, nil
}
