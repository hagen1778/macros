package main

import (
	"flag"
	"github.com/vharitonsky/iniflags"
	"lib/util"
)

var (
	confFile = flag.String("confFile", "conf.yml", "Path to yml conf file")
)

func main() {
	iniflags.Parse()
	util.LogAllFlags()

	Init(*confFile)
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
