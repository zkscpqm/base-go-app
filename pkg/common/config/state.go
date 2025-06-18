package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	activeState = defaultState()
)

type serverState struct {
	Enabled    bool
	Port       int
	TlsEnabled bool
}

type apiState struct {
	REST        serverState
	Certificate string
	Key         string
}

type logState struct {
	Level string
	File  string
}

type dbState struct {
	Engine        string
	Username      string
	Password      string
	Host          string
	Port          int
	Database      string
	ScriptsFolder string
}

type state struct {
	DisableAutoSave bool
	API             apiState
	Log             logState
	Database        dbState
}

func (s state) asConfig() (rv Configuration) {
	return Configuration{
		DisableAutoSave: s.DisableAutoSave,
		API: ApiConfig{
			REST: Server{
				Enabled:    s.API.REST.Enabled,
				Port:       s.API.REST.Port,
				TlsEnabled: s.API.REST.TlsEnabled,
			},
			Certificate: s.API.Certificate,
			Key:         s.API.Key,
		},
		Log: LogConfig{
			Level: s.Log.Level,
			File:  s.Log.File,
		},
		Database: DbConfig{
			Engine:   s.Database.Engine,
			Username: s.Database.Username,
			Password: s.Database.Password,
			Host:     s.Database.Host,
			Port:     s.Database.Port,
			Database: s.Database.Database,
		},
	}
}

func defaultState() state {
	return state{}
}

func autoSaveDisabled() bool {
	return activeState.DisableAutoSave
}

func Init(config Configuration) error {
	return SetRunningConfig(config)
}

func SetRunningConfig(config Configuration) error {
	activeState = config.state()
	return autoSave()
}

func DisableAutoSave(disable bool) error {
	activeState.DisableAutoSave = disable
	return autoSave()
}

func autoSave() error {
	if autoSaveDisabled() {
		return nil
	}
	if err := CopyRunningConfigToStartupConfig(); err != nil {
		return fmt.Errorf("failed to auto save config: %v", err)
	}
	return nil
}

func GetRunningConfig() Configuration {
	return activeState.asConfig()
}

func CopyRunningConfigToStartupConfig() error {
	f, err := os.OpenFile(configFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("failed to open state log file %s: %v", configFilePath, err)
	}
	defer f.Close()
	cfg := activeState.asConfig()
	return json.NewEncoder(f).Encode(&cfg)
}
