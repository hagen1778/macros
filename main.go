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
			if e.String() == ";" && e.Value == 1 && kb.R_SHIFT {
				fmt.Println("R_SHIFT + ; was pressed!")
				kb.Execute("kiss my shiny metal but")

			}
		}
	}
}