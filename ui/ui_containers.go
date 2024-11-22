package ui

import (
	"fyne.io/fyne/v2"
)

func (u *Ui) RefreshTop(items ...fyne.CanvasObject) {

	u.containerTop.RemoveAll()
	for _, item := range items {
		u.containerTop.Add(item)
	}
	u.containerTop.Refresh()

	return
}
func (u *Ui) RefreshMainLeft(items ...fyne.CanvasObject) {

	u.containerMainLeft.RemoveAll()
	for _, item := range items {
		u.containerMainLeft.Add(item)
	}

	u.containerMainLeft.Refresh()

	return
}
func (u *Ui) RefreshMainContent(items ...fyne.CanvasObject) {
	u.containerMainPrevItems = u.containerMainContent.Objects
	u.containerMainContent.RemoveAll()
	for _, item := range items {
		u.containerMainContent.Add(item)
	}
	u.containerMainContent.Refresh()

	return
}
func (u *Ui) PrevMainContent() {
	u.containerMainContent.RemoveAll()
	for _, item := range u.containerMainPrevItems {
		u.containerMainContent.Add(item)
	}
	u.containerMainContent.Refresh()

	return
}

func (u *Ui) RefreshBottom(items ...fyne.CanvasObject) {

	u.containerBottom.RemoveAll()
	for _, item := range items {
		u.containerBottom.Add(item)
	}

	u.containerBottom.Refresh()

	return
}
