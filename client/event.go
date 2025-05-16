package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
	"time"
)

var keyCode = map[uint16]string{}

type EventManager struct {
	hotKey   []string
	window   fyne.Window
	keyDown  chan int
	KeyOpen  chan int
	KeyCheck chan []uint16
}

func NewEvent(hotKey []string, window fyne.Window) *EventManager {
	return &EventManager{
		hotKey:   hotKey,
		window:   window,
		keyDown:  make(chan int),
		KeyOpen:  make(chan int),
		KeyCheck: make(chan []uint16),
	}
}

func (em EventManager) GetKeys() {
	evChan := hook.Start()
	defer hook.End()

	holdMap := make(map[uint16]bool)

	for ev := range evChan {
		if ev.Kind != hook.KeyHold && ev.Kind != hook.KeyDown && ev.Kind != hook.KeyUp {
			continue
		}
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown {
			holdMap[ev.Rawcode] = true
		}
		if ev.Kind == hook.KeyUp {
			fmt.Println(holdMap)
			keys := []uint16{}
			for k, _ := range holdMap {
				keys = append(keys, k)
			}
			em.KeyCheck <- keys
			return
		}
	}
}

func (em EventManager) Listen(keys []uint16) {
	evChan := hook.Start()
	defer hook.End()

	holdMap := make(map[uint16]bool)

	for ev := range evChan {
		if ev.Kind != hook.KeyHold && ev.Kind != hook.KeyDown && ev.Kind != hook.KeyUp {
			continue
		}
		if ev.Kind == hook.KeyHold || ev.Kind == hook.KeyDown {
			holdMap[ev.Rawcode] = true
		}
		if ev.Kind == hook.KeyUp {
			fmt.Println(holdMap)
			for _, key := range keys {
				if _, ok := holdMap[key]; !ok {
					return
				}
			}
			em.KeyOpen <- 1
			return
		}
	}
}

func (em *EventManager) Send() {
	time.Sleep(time.Millisecond * 500)
	err := robotgo.KeyTap("v", "cmd")
	if err != nil {
		fmt.Println(err)
	}
}
