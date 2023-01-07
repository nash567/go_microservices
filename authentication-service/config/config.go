package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/authentication-service/internal/db"
	"github.com/jinzhu/configor"
	"github.com/joho/godotenv"
)

var ErrInvalidFileExtension = errors.New("invalid file extension")

type AppConfig struct {
	Database *db.Config
	URLHost  string
	URLPort  int
}

// LoadConfig load configs from file,supported formats: yml, json, toml, .env respectively.
// This will load configs only once in the application. Has a broader scope for using custom config files in the applicaion.

func LoadConfig(filenames ...string) (*AppConfig, error) {
	fmt.Printf("reached here")

	loadFiles := make([]string, 0, len(filenames))
	envFiles := make([]string, 0, len(filenames))

	for _, file := range filenames {
		fileParts := strings.Split(file, ".")
		ext := fileParts[len(fileParts)-1]

		switch ext {
		case "yml", "json", "yaml", "toml":
			loadFiles = append(loadFiles, file)
		case "env":
			envFiles = append(envFiles, file)
		default:
			return nil, ErrInvalidFileExtension
		}
	}

	if len(envFiles) > 0 {
		fmt.Printf("sadad %d", len(envFiles))
		// it will set the env variables from env file to environment and make it available for loading in config
		err := godotenv.Load(envFiles...)
		if err != nil {
			return nil, fmt.Errorf("error loading env files: (%s):%w", strings.Join(envFiles, ","), err)
		}

	}
	_cfg, err := loadConfig(loadFiles...)
	if err != nil {
		return nil, fmt.Errorf("error loading config file: (%s):%w", strings.Join(loadFiles, ","), err)
	}

	return _cfg, nil
}

func loadConfig(fileName ...string) (*AppConfig, error) {
	var appConfig AppConfig

	conf := newConf()
	if err := conf.Load(&appConfig, fileName...); err != nil {
		return nil, fmt.Errorf("failed to load config file: %w", err)
	}

	return &appConfig, nil

}

func newConf() *configor.Configor {
	conf := configor.Config{ENVPrefix: "AUTH_SERVICE"}
	config := configor.New(&conf)
	return config
}
