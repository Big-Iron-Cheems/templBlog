package config

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"os"
)

// Config represents the data from the `config.json` file.
type Config struct {
	Port           int    `json:"port"`           // Port to run the server on (default: 8080)
	PostsDirectory string `json:"postsDirectory"` // Path to the posts directory containing MD files (default: posts)
	LogFileName    string `json:"logFileName"`    // Name of the log file (default: templBlog.log )
	Verbose        bool   `json:"verbose"`        // Whether to print debug messages
	EnableLogFile  bool   `json:"enableLogFile"`  // Whether to write logs to a file
}

var GlobalConfig Config

// InitConfig populates the GlobalConfig variable with the configuration from the `config.json` file.
func InitConfig() {
	file, err := os.Open("configs/config.json")
	if err != nil {
		log.Panic().Msgf("Cannot open config file: %v", err)
	}
	defer file.Close()

	if err = json.NewDecoder(file).Decode(&GlobalConfig); err != nil {
		log.Panic().Msgf("Error decoding config file: %v", err)
	}

	if GlobalConfig.Port == 0 {
		GlobalConfig.Port = 8080
	}

	if GlobalConfig.PostsDirectory == "" {
		GlobalConfig.PostsDirectory = "posts"
	}

	if GlobalConfig.LogFileName == "" {
		GlobalConfig.LogFileName = "templBlog.log"
	}
}
