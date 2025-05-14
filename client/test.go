package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

func test() {
	a := app.New()
	w := a.NewWindow("SysTray")

	if desk, ok := a.(desktop.App); ok {
		m := fyne.NewMenu("MyApp",
			fyne.NewMenuItem("Show", func() {
				w.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	w.SetContent(widget.NewLabel("Fyne System Tray"))
	w.SetCloseIntercept(func() {
		w.Hide()
	})

	go func() {
		add(w)
		//low()
		//event()

		//t := time.NewTicker(time.Second * 10)
		//for range t.C {
		//	fyne.Do(func() {
		//		robotgo.KeyTap("v", "cmd")
		//	})
		//}
	}()

	w.ShowAndRun()
}

func add(window fyne.Window) {
	fmt.Println("--- Please press ctrl + shift + q to stop hook ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})

	fmt.Println("--- Please press w---")
	hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
		fmt.Println("w")
		window.Show()
	})

	s := hook.Start()
	<-hook.Process(s)
}

func low() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
	}
}

func event() {
	ok := hook.AddEvents("q", "ctrl", "shift")
	if ok {
		fmt.Println("add events...")
	}

	keve := hook.AddEvent("k")
	if keve {
		fmt.Println("you press... ", "k")
	}

	mleft := hook.AddEvent("mleft")
	if mleft {
		fmt.Println("you press... ", "mouse left button")
	}
}

func test1() {
	robotgo.TypeStr("Hello World")
	robotgo.TypeStr("だんしゃり", 0, 1)
	// robotgo.TypeStr("テストする")

	robotgo.TypeStr("Hi, Seattle space needle, Golden gate bridge, One world trade center.")
	robotgo.TypeStr("Hi galaxy, hi stars, hi MT.Rainier, hi sea. こんにちは世界.")
	robotgo.Sleep(1)

	// ustr := uint32(robotgo.CharCodeAt("Test", 0))
	// robotgo.UnicodeType(ustr)

	robotgo.KeySleep = 100
	robotgo.KeyTap("enter")
	// robotgo.TypeStr("en")
	robotgo.KeyTap("i", "alt", "cmd")

	arr := []string{"alt", "cmd"}
	robotgo.KeyTap("i", arr)

	robotgo.MilliSleep(100)
	robotgo.KeyToggle("a")
	robotgo.KeyToggle("a", "up")

	robotgo.WriteAll("Test")
	text, err := robotgo.ReadAll()
	if err == nil {
		fmt.Println(text)
	}
}
