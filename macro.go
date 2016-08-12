package main

import (
	"fmt"
	"github.com/hagen1778/macros/keyboard"
	"log"
	"strings"
	"time"
)

type Macro struct {
	cfg   *config
	kb    *keyboard.InputDevice
	mList []*macros
}

func (m *Macro) checkEvent(e *keyboard.InputEvent) macrosFunc {
	if e.Type == keyboard.EV_KEY && e.Value == 1 {
		for _, macros := range m.mList {
			if macros.isModifiersPressed(m.kb) && macros.key == e.String() {
				return macros.Run
			}
		}
	}

	return nil
}

func (m *Macro) listen() {
	kb_events, err := m.kb.Listen()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen for %s events!\n", m.kb.Name)

	for e := range kb_events {
		if f := m.checkEvent(&e); f != nil {
			f(m.kb, m.cfg.defaultDelay)
		}
	}
}

func (m *Macro) applyRules() error {
	for _, rule := range m.cfg.Rules {
		macros := macros{
			name: rule.Name,
		}

		for _, v := range strings.Fields(rule.Activate) {
			if keyboard.IsModifier(v) {
				macros.modifiers = append(macros.modifiers, v)
				continue
			}

			macros.key = v
		}

		for _, rule := range rule.Scenario {
			action, err := rule.convertToAction()
			if err != nil {
				return fmt.Errorf("invalid rule %#v", rule)
			}

			macros.actions = append(macros.actions, action)
		}
		m.mList = append(m.mList, &macros)
		log.Printf("Macro %q registred", macros)
	}

	return nil
}

func Init(confPath string) {
	cfg, err := loadFile(confPath)
	if err != nil {
		log.Fatalf("Error while loading config >> %s", err)
	}

	macro := Macro{
		cfg: cfg,
	}
	macro.applyRules()

	kb, err := keyboard.Init()
	if err != nil {
		log.Fatalf("Error while initing keyboard: %s", err)
	}
	macro.kb = kb
	macro.listen()
}

type macrosFunc func(d *keyboard.InputDevice, dur time.Duration)

type macros struct {
	name      string
	key       string
	modifiers []string
	actions   []runner
}

func (m *macros) Run(d *keyboard.InputDevice, dur time.Duration) {
	for _, a := range m.actions {
		a.Run(d)
		time.Sleep(dur)
	}
}

func (m *macros) isModifiersPressed(kb *keyboard.InputDevice) bool {
	result := true
	for _, mod := range m.modifiers {
		if pressed, _ := kb.Modifiers[mod]; !pressed {
			result = false
		}
	}

	return result
}

func (m macros) String() string {
	modifiers := ""
	for _, mod := range m.modifiers {
		modifiers += mod + " + "
	}

	return fmt.Sprintf("%s: %s%s", m.name, modifiers, m.key)
}
