package conf

import (
	"io"
	"os"

	"github.com/hugefiver/fakessh/modules/gitserver"
	"github.com/pelletier/go-toml/v2"
)

type AppConfig struct {
	BaseConfig

	Modules ModulesConfig `toml:"modules"`
}

type BaseConfig struct {
	Server struct {
		ServPort   string `toml:"bind"`
		SSHVersion string `toml:"version"`

		MaxTry    int `toml:"max_try"`
		Delay     int `toml:"delay"`
		Deviation int `toml:"deviation"`

		AntiScan bool `toml:"anti_scan"`
	} `toml:"server"`

	Log struct {
		LogFile     string `toml:"file"`
		LogLevel    string `toml:"level"`
		LogFormat   string `toml:"format"`
		IsLogPasswd bool   `toml:"log_passwd"`
	} `toml:"log"`

	Key struct {
		KeyFiles []string `toml:"key"`
		KeyType  string   `toml:"type"`
	} `toml:"key"`
}

type ModulesConfig struct {
	GitServer gitserver.Config `toml:"gitserver"`
}

func (c *BaseConfig) FillDefault() error {
	c.Server.ServPort = DefaultBind
	c.Server.SSHVersion = DefaultSSHVersion
	c.Server.Delay = DefaultDelay
	c.Server.Deviation = DefaultDeviation
	c.Server.AntiScan = DefaultEnableAntiScan

	c.Log.LogLevel = DefaultLogLevel
	c.Log.LogFormat = DefaultLogFormat
	c.Log.IsLogPasswd = false

	c.Key.KeyType = DefaultKeyType

	return nil
}

// func (c *AppConfig) FillDefault() error {
// 	if err := c.BaseConfig.FillDefault(); err != nil {
// 		return err
// 	}

// 	if err := c.Modules.GitServer.FillDefault(); err != nil {
// 		return err
// 	}
// 	return nil
// }

func NewDefaultAppConfig() *AppConfig {
	c := &AppConfig{}

	c.BaseConfig.FillDefault()

	return c
}

func ParseConfig(s []byte) (*AppConfig, error) {
	var config AppConfig
	config.FillDefault()

	if err := toml.Unmarshal(s, &config); err != nil {
		return nil, err
	}

	// Fill default values of Modules.GitServer
	if err := config.Modules.GitServer.FillDefault(); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadFromFile(file string) (*AppConfig, error) {
	r, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	s, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseConfig(s)
}

func MergeConfig(c *AppConfig, f *FlagArgsStruct, set StringSet) error {
	var enableAnti, disableAnti bool

	set.ForEach(func(s string) error {
		switch s {
		case FlagBind:
			c.Server.ServPort = f.ServPort
		case FlagSSHVersion:
			c.Server.SSHVersion = f.SSHVersion
		case FlagMaxTry:
			c.Server.MaxTry = f.MaxTry
		case FlagDelay:
			c.Server.Delay = f.Delay
		case FlagDeviation:
			c.Server.Deviation = f.Deviation

		case FlagLogFile:
			c.Log.LogFile = f.LogFile
		case FlagLogLevel:
			c.Log.LogLevel = f.LogLevel
		case FlagLogFormat:
			c.Log.LogFormat = f.LogFormat
		case FlagLogPasswd:
			c.Log.IsLogPasswd = f.IsLogPasswd

		case FlagKeyPaths:
			c.Key.KeyFiles = f.KeyFiles
		case FlagKeyType:
			c.Key.KeyType = f.KeyType
		case FlagEnableAntiScan:
			enableAnti = true
		case FlagDisableAntiScan:
			disableAnti = true
		}
		return nil
	})

	if enableAnti || disableAnti {
		c.Server.AntiScan = enableAnti
	}
	return nil
}
