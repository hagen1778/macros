package main

import (
	"github.com/hagen1778/macros/keyboard"
	"flag"
	"fmt"
	"log"
	"time"
)

var (
	confFile = flag.String("confFile", "conf.yml", "Path to yml conf file")
)

func main() {
	cfg, err := loadFile(*confFile)
	if err != nil {
		log.Fatalf("Error while loading config >> %s", err)
	}
	cfg.Validate()

	kb, err := keyboard.Init()
	if err != nil {
		fmt.Println(err)
		return
	}

	kb_events, err := kb.Listen()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Listen for %s events!\n", kb.Name)

	for e := range kb_events {
		if f, e := cfg.CheckEvent(&e); e {
			f(kb)
		}
		////listen only key stroke event
		//if e.Type == keyboard.EV_KEY {
		//	if e.String() == "Z" && e.Value == 1 {
		//		//kb.Print("kiss my shiny metal but")
		//		htop(kb, "skype")
		//	}
		//}
	}
}

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
}
