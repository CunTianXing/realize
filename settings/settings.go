package settings

import (
	yaml "gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"time"
)

// settings const
const (
	permission = 0775
	directory  = ".realize"
)

// Settings defines a group of general settings
type Settings struct {
	Config    `yaml:",inline" json:"config"`
	Resources `yaml:"resources" json:"resources"`
	Server    `yaml:"server,omitempty" json:"server,omitempty"`
}

// Config defines structural options
type Config struct {
	Create bool  `yaml:"-" json:"-"`
	Flimit int64 `yaml:"flimit,omitempty" json:"flimit,omitempty"`
	Legacy `yaml:"legacy,omitempty" json:"legacy,omitempty"`
}

// Legacy configuration
type Legacy struct {
	Status   bool          `yaml:"status" json:"status"`
	Interval time.Duration `yaml:"interval" json:"interval"`
}

// Server settings, used for the web panel
type Server struct {
	Status bool   `yaml:"status" json:"status"`
	Open   bool   `yaml:"open" json:"open"`
	Host   string `yaml:"host" json:"host"`
	Port   int    `yaml:"port" json:"port"`
}

// Resources defines the files generated by realize
type Resources struct {
	Config  string `yaml:"-" json:"-"`
	Outputs string `yaml:"outputs" json:"outputs"`
	Logs    string `yaml:"logs" json:"log"`
	Errors  string `yaml:"errors" json:"error"`
}

// Read from config file
func (s *Settings) Read(out interface{}) error {
	localConfigPath := s.Resources.Config
	// backward compatibility
	path := filepath.Join(directory, s.Resources.Config)
	if _, err := os.Stat(path); err == nil {
		localConfigPath = path
	}
	content, err := s.Stream(localConfigPath)
	if err == nil {
		err = yaml.Unmarshal(content, out)
		return err
	}
	return err
}

// Record create and unmarshal the yaml config file
func (s *Settings) Record(out interface{}) error {
	if s.Config.Create {
		y, err := yaml.Marshal(out)
		if err != nil {
			return err
		}
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			if err = os.Mkdir(directory, permission); err != nil {
				return s.Write(s.Resources.Config, y)
			}
		}
		return s.Write(filepath.Join(directory, s.Resources.Config), y)
	}
	return nil
}

// Remove realize folder
func (s *Settings) Remove(d string) error {
	_, err := os.Stat(d)
	if !os.IsNotExist(err) {
		return os.RemoveAll(d)
	}
	return err
}
