package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type MyList struct {
	widget.BaseWidget
	items      []string
	texts      []*canvas.Text
	container  *fyne.Container
	selected   int
	onSelected func(index int, value string)
}

func NewMyList(items []string, onSelected func(index int, value string)) *MyList {
	list := &MyList{
		items:      items,
		selected:   0,
		onSelected: onSelected,
	}

	list.texts = make([]*canvas.Text, len(items))
	for i, item := range items {
		txt := canvas.NewText(item, color.Black)
		txt.TextSize = 16
		list.texts[i] = txt
	}

	list.container = container.NewVBox(make([]fyne.CanvasObject, len(list.texts))...)
	for i, txt := range list.texts {
		list.container.Objects[i] = txt
	}

	list.ExtendBaseWidget(list)
	list.updateSelection()

	return list
}

func (l *MyList) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.container)
}

func (l *MyList) FocusGained() {}
func (l *MyList) FocusLost()   {}
func (l *MyList) Focused() bool {
	return true
}
func (l *MyList) TypedRune(r rune) {}
func (l *MyList) FocusedStyle() fyne.Focusable {
	return l
}

func (l *MyList) TypedKey(event *fyne.KeyEvent) {
	if event.Name == fyne.KeyDown {
		if l.selected < len(l.items)-1 {
			l.selected++
			l.updateSelection()
		}
	} else if event.Name == fyne.KeyUp {
		if l.selected > 0 {
			l.selected--
			l.updateSelection()
		}
	}
}

func (l *MyList) updateSelection() {
	for i, txt := range l.texts {
		if i == l.selected {
			txt.Color = color.NRGBA{R: 0x00, G: 0x88, B: 0xff, A: 0xff}
		} else {
			txt.Color = color.Black
		}
		txt.Refresh()
	}

	if l.onSelected != nil {
		l.onSelected(l.selected, l.items[l.selected])
	}
}
