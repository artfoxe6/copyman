package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"strings"
	"time"
)

type EventManager struct {
	hotKey []string
	window fyne.Window
}

func NewEvent(hotKey []string, window fyne.Window) *EventManager {
	return &EventManager{
		hotKey: hotKey,
		window: window,
	}
}

func (em *EventManager) Start() {
	go func() {
		fmt.Println("listen " + strings.Join(em.hotKey, "-"))
		hook.Register(hook.KeyDown, em.hotKey, func(e hook.Event) {
			fyne.Do(func() {
				em.window.Show()
			})
		})
		s := hook.Start()
		<-hook.Process(s)
	}()
}

func (em *EventManager) Send() {
	time.Sleep(time.Millisecond * 300)
	err := robotgo.KeyTap("v", "cmd")
	if err != nil {
		fmt.Println(err)
	}
}
