package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	App                    fyne.App
	Window                 fyne.Window
	containerTop           *fyne.Container
	containerMainContent   *fyne.Container
	containerMainPrevItems []fyne.CanvasObject
	containerMainLeft      *fyne.Container
	containerBottom        *fyne.Container
	MenuTreeMainLeft       *widget.Tree
}

func NewUi(title string, w float32, h float32) (view *Ui) {
	fyneApp := app.New()
	window := fyneApp.NewWindow(title)

	if w > 0 && h > 0 {
		window.Resize(fyne.NewSize(w, h))
	} else {
		window.FixedSize()
	}

	view = &Ui{
		App:    fyneApp,
		Window: window,
	}

	view.Init()
	return
}

func (u *Ui) Init() {

	boxTop := container.NewStack()
	boxMainLeft := container.NewStack()
	boxMainContent := container.NewStack()
	split := container.NewHSplit(boxMainLeft, boxMainContent)
	split.Offset = 0.2
	boxMain := container.NewBorder(nil, nil, nil, nil, split)

	boxBottom := container.New(layout.NewHBoxLayout())
	u.containerTop = boxTop
	u.containerMainLeft = boxMainLeft
	u.containerMainContent = boxMainContent
	u.containerBottom = boxBottom
	content := container.NewBorder(boxTop, boxBottom, nil, nil, boxMain)
	//content := split
	u.Window.SetContent(content)

}

func (u *Ui) ShowAndRun() {
	if nil != u.Window {
		u.Window.ShowAndRun()
		return
	}
	panic("程序未初始化")
}
