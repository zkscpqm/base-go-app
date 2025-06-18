package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

var (
	configFilePath, _ = filepath.Abs("config.json")
)

func init() {
	val, ok := syscall.Getenv("LEDGER_CONFIG_FILE")
	if ok {
		configFilePath = val
	}
}

type Server struct {
	Enabled    bool `json:"enabled,omitempty"`
	Port       int  `json:"port,omitempty"`
	TlsEnabled bool `json:"tls_enabled,omitempty"`
}

type ApiConfig struct {
	REST        Server `json:"rest,omitempty"`
	Certificate string `json:"certificate,omitempty"`
	Key         string `json:"key,omitempty"`
}

type LogConfig struct {
	Level string `json:"level,omitempty"`
	File  string `json:"file,omitempty"`
}

type DbConfig struct {
	Engine   string `json:"engine"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
	Database string `json:"database"`
}

type Configuration struct {
	DisableAutoSave bool      `json:"disable_auto_save,omitempty"`
	API             ApiConfig `json:"api,omitempty"`
	Log             LogConfig `json:"log,omitempty"`
	Database        DbConfig  `json:"database,omitempty"`
}

func Empty() Configuration {
	return Configuration{}
}

func (c Configuration) state() state {
	return state{
		DisableAutoSave: c.DisableAutoSave,
		API: apiState{
			REST: serverState{
				Enabled:    c.API.REST.Enabled,
				Port:       c.API.REST.Port,
				TlsEnabled: c.API.REST.TlsEnabled,
			},
			Certificate: c.API.Certificate,
			Key:         c.API.Key,
		},
		Log: logState{
			Level: c.Log.Level,
			File:  c.Log.File,
		},
		Database: dbState{
			Engine:   c.Database.Engine,
			Username: c.Database.Username,
			Password: c.Database.Password,
			Host:     c.Database.Host,
			Port:     c.Database.Port,
			Database: c.Database.Database,
		},
	}
}

func FromFile() (*Configuration, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file `%s`: %v", configFilePath, err)
	}
	defer file.Close()

	var config Configuration
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file `%s`: %v", configFilePath, err)
	}

	return &config, nil
}

func (c Configuration) Copy() Configuration {
	return Configuration{
		DisableAutoSave: c.DisableAutoSave,
		API:             c.API,
		Log:             c.Log,
		Database:        c.Database,
	}
}

func (c Configuration) String() string {
	b, _ := json.MarshalIndent(c, "", "    ")
	return string(b)
}

func GetStartupConfig() (Configuration, error) {
	cfg, err := FromFile()
	if err != nil {
		return Empty(), fmt.Errorf("failed to load startup config: %s", err)
	}
	return *cfg, nil
}
