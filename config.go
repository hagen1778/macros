package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type config struct {
	Rules  []*rule `yaml:"rules"`
	Global global  `yaml:"global"`

	defaultDelay time.Duration

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

type global struct {
	ModeButton   string `yaml:"mode_button"`
	DefaultDelay string `yaml:"default_delay"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

type rule struct {
	Name     string `yaml:"name"`
	Activate string `yaml:"activate"`
	Scenario []step `yaml:"scenario"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

func loadFile(filename string) (*config, error) {
	if stat, err := os.Stat(filename); err != nil {
		return nil, fmt.Errorf("cannot get file info: %s", err)
	} else if stat.IsDir() {
		return nil, fmt.Errorf("is a directory")
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := load(string(content))
	if err != nil {
		return nil, err
	}

	err = cfg.validate()
	return cfg, nil
}

func load(s string) (*config, error) {
	cfg := &config{}
	err := yaml.Unmarshal([]byte(s), cfg)
	if err != nil {
		return nil, err
	}

	err = checkOverflow(cfg.XXX)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) validate() error {
	dur, err := time.ParseDuration(cfg.Global.DefaultDelay)
	if err != nil {
		return err
	}

	cfg.defaultDelay = dur
	return nil
}

func checkOverflow(m map[string]interface{}) error {
	if len(m) > 0 {
		var keys []string
		for k := range m {
			keys = append(keys, k)
		}
		return fmt.Errorf("unknown fields in config: %s", strings.Join(keys, ", "))
	}
	return nil
}
