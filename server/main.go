package main

import (
	"fmt"
	hook "github.com/robotn/gohook"
)

func main() {
	evChan := hook.Start()
	defer hook.End()

	holdMap := make(map[uint16]bool)

	for ev := range evChan {
		if ev.Kind >= 6 && ev.Kind <= 11 {
			continue
		}
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown {
			holdMap[ev.Rawcode] = true
		}
		if ev.Kind == hook.KeyUp {
			fmt.Println(holdMap)
			holdMap = make(map[uint16]bool)
		}
	}
}
