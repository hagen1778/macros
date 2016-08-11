package main

import (
	"fmt"
	"time"
	"github.com/hagen1778/macros/keyboard"
)

type rule map[string]string

func (r rule) convertToAction() (a actioner, err error) {
	for key, value := range r {
		switch key {
		case "sleep":
			var dur time.Duration
			dur, err = time.ParseDuration(value)
			if err != nil {
				return
			}
			a = sleep{dur}
		case "print":
			a = print{value}
		case "press":
			a = press{value}
		default:
			err = fmt.Errorf("Wrong rule type")
		}
	}

	return
}

type actioner interface {
	Run(d *keyboard.InputDevice)
}

type press struct {
	value string
}

type print struct {
	value string
}

type sleep struct {
	value time.Duration
}

func (a press) Run(d *keyboard.InputDevice) {
	d.Press(a.value)
}

func (a print) Run(d *keyboard.InputDevice) {
	d.Print(a.value)
}

func (a sleep) Run(d *keyboard.InputDevice) {
	time.Sleep(a.value)
}
