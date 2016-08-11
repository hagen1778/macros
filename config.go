package main

import (
	"gopkg.in/yaml.v2"
	"fmt"
	"github.com/hagen1778/macros/keyboard"
	"io/ioutil"
	"os"
	"strings"
	"log"
)

type config struct {
	Macroses []*Macros `yaml:"macros"`
	Global   global    `yaml:"global"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

type global struct {
	ModeButton   string `yaml:"mode_button"`
	DefaultDelay string `yaml:"default_delay"`

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

type Macros struct {
	Name     string `yaml:"name"`
	Activate string `yaml:"activate"`
	Scenario []rule `yaml:"scenario"`

	actions []actioner

	// Catches all undefined fields and must be empty after parsing.
	XXX map[string]interface{} `yaml:",inline"`
}

func (m *Macros) Run(d *keyboard.InputDevice) {
	for _, a := range m.actions {
		a.Run(d)
	}
}

// Validate checks whether all names and values in the action list
// are valid.
func (cfg *config) Validate() error {
	for _, macros := range cfg.Macroses {
		for _, rule := range macros.Scenario {
			action, err := rule.convertToAction()
			if err != nil {
				return fmt.Errorf("invalid rule %#v", rule)
			}

			macros.actions = append(macros.actions, action)
		}
		log.Printf("Macros %q registred. Listen: %s", macros.Name, macros.Activate)
	}

	return nil
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

type MacrosFunc func (d *keyboard.InputDevice)

func (cfg *config) CheckEvent(e *keyboard.InputEvent) (MacrosFunc, bool) {
	if e.Type == keyboard.EV_KEY && e.Value == 1 {
		for _, macros := range cfg.Macroses {
			if macros.Activate == e.String() {
				return macros.Run, true
			}
		}
	}

	return nil, false
}
