package main

import (
	"flag"
	"fmt"
	"github.com/hagen1778/macros/keyboard"
	"log"
)

var (
	confFile = flag.String("confFile", "conf.yml", "Path to yml conf file")
)

func Init() {
	cfg, err := loadFile(*confFile)
	if err != nil {
		log.Fatalf("Error while loading config >> %s", err)
	}

	macro := Macro{
		cfg: cfg,
	}
	macro.applyRules()

	kb, err := keyboard.Init()
	if err != nil {
		fmt.Println(err)
		return
	}
	macro.kb = kb
	macro.listen()
}

type Macro struct {
	cfg   *config
	kb    *keyboard.InputDevice
	mList []*macros
}

type macrosFunc func(d *keyboard.InputDevice)

func (m *Macro) checkEvent(e *keyboard.InputEvent) (macrosFunc, bool) {
	if e.Type == keyboard.EV_KEY && e.Value == 1 {
		for _, macros := range m.mList {
			if macros.activate == e.String() {
				return macros.Run, true
			}
		}
	}

	return nil, false
}

func (m *Macro) listen() {
	kb_events, err := m.kb.Listen()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen for %s events!\n", m.kb.Name)

	for e := range kb_events {
		if f, happen := m.checkEvent(&e); happen {
			f(m.kb)
		}
	}
}

func (m *Macro) applyRules() error {
	for _, rule := range m.cfg.Rules {
		macros := macros{
			name:     rule.Name,
			activate: rule.Activate,
		}
		for _, rule := range rule.Scenario {
			action, err := rule.convertToAction()
			if err != nil {
				return fmt.Errorf("invalid rule %#v", rule)
			}

			macros.actions = append(macros.actions, action)
		}
		m.mList = append(m.mList, &macros)
		log.Printf("Macro %q registred. Listen: %s", macros.name, macros.activate)
	}

	return nil
}

type macros struct {
	name     string
	activate string
	actions  []runner
}

func (m *macros) Run(d *keyboard.InputDevice) {
	for _, a := range m.actions {
		a.Run(d)
	}
}

func main() {
	Init()
}

/*
func htop(kb *keyboard.InputDevice, search string) {
	kb.Press("L_CTRL L_ALT T")
	time.Sleep(500 * time.Millisecond)

	kb.Print("htop")
	time.Sleep(100 * time.Millisecond)

	kb.Press("ENTER")
	time.Sleep(100 * time.Millisecond)

	kb.Press("F4")
	time.Sleep(100 * time.Millisecond)

	kb.Print(search)
}*/
