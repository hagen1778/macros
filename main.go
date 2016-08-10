package main

import (
	"github.com/hagen1778/macros/keyboard"
	"fmt"
)

func main() {
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
		//listen only key stroke event
		if e.Type == keyboard.EV_KEY {
			//fmt.Println(i.KeyString())
			if e.String() == "F12" {
				//Run(Htop)
				fmt.Println("f12 was pressed!")
			}
		}
	}
}