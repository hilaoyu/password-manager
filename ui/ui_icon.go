package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func (u *Ui) IconCopy(callback func(copyFunc func(value string))) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		callback(u.UtilToClipboard)
	})))
	return
}
func (u *Ui) IconEdit(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DocumentCreateIcon(), callback)))
	return
}
func (u *Ui) IconSearch(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.SearchIcon(), callback)))
	return
}
func (u *Ui) IconAdd(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), callback)))
	return
}
func (u *Ui) IconDelete(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DeleteIcon(), callback)))
	return
}
func (u *Ui) IconRemove(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentClearIcon(), callback)))
	return
}
func (u *Ui) IconSave(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DocumentSaveIcon(), callback)))
	return
}
