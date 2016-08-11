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
			//fmt.Printf("%#v\n",e)
			if e.String() == "Z" && e.Value == 1  {
				kb.Print("kiss my shiny metal but")
				kb.Press("L_CTRL L_ALT T")
			}
		}
	}
}