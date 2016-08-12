package main

import (
	"fmt"
	"github.com/hagen1778/macros/keyboard"
	"time"
)

type step map[string]string

func (s step) convertToAction() (a runner, err error) {
	for key, value := range s {
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

type runner interface {
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
