package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func IconCopy(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
		callback()
	})))
	return
}
func IconEdit(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DocumentCreateIcon(), callback)))
	return
}
func IconSearch(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.SearchIcon(), callback)))
	return
}
func IconAdd(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentAddIcon(), callback)))
	return
}
func IconDelete(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DeleteIcon(), callback)))
	return
}
func IconRemove(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentClearIcon(), callback)))
	return
}
func IconSave(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.DocumentSaveIcon(), callback)))
	return
}
func IconClear(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.ContentClearIcon(), callback)))
	return
}
func IconVisibility(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.VisibilityIcon(), callback)))
	return
}
func IconVisibilityOff(callback func()) (content *fyne.Container) {
	content = container.NewStack(widget.NewToolbar(widget.NewToolbarAction(theme.VisibilityOffIcon(), callback)))
	return
}
