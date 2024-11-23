package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/mobile"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strconv"
)

type NumericalEntry struct {
	widget.Entry
}

func NewNumericalEntry() *NumericalEntry {
	entry := &NumericalEntry{}
	entry.ExtendBaseWidget(entry)
	return entry
}

func (e *NumericalEntry) TypedRune(r rune) {
	if (r >= '0' && r <= '9') || r == '.' || r == ',' {
		e.Entry.TypedRune(r)
	}
}

func (e *NumericalEntry) TypedShortcut(shortcut fyne.Shortcut) {
	paste, ok := shortcut.(*fyne.ShortcutPaste)
	if !ok {
		e.Entry.TypedShortcut(shortcut)
		return
	}

	content := paste.Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err == nil {
		e.Entry.TypedShortcut(shortcut)
	}
}

func (e *NumericalEntry) Keyboard() mobile.KeyboardType {
	return mobile.NumberKeyboard
}

func NewLabelWrap(text string) (label *widget.Label) {
	label = widget.NewLabel(text)
	label.Wrapping = fyne.TextWrapBreak
	return
}
func NewRectangleWithSize(fillColor color.Color, w float32, h float32) (rectangle *canvas.Rectangle) {
	rectangle = canvas.NewRectangle(fillColor)
	rectangle.SetMinSize(fyne.NewSize(w, h))
	return
}
func NewScrollWithSize(content fyne.CanvasObject, w float32, h float32) (scroll *container.Scroll) {
	scroll = container.NewScroll(content)
	scroll.SetMinSize(fyne.NewSize(w, h))
	return
}

func NewContainerWithSize(w float32, h float32, objects ...fyne.CanvasObject) (c *fyne.Container) {
	var wFill fyne.CanvasObject
	var hFill fyne.CanvasObject
	if w > 0 {
		wFill = NewRectangleWithSize(color.Transparent, w, 0)
	} else {
		wFill = nil
	}
	if h > 0 {
		hFill = NewRectangleWithSize(color.Transparent, 0, h)
	} else {
		hFill = nil
	}

	c = container.NewBorder(nil, wFill, nil, hFill, objects...)
	return
}

func NewRichTextFromMarkdownWrap(text string) (richText *widget.RichText) {
	richText = widget.NewRichTextFromMarkdown(text)
	richText.Wrapping = fyne.TextWrapBreak
	return
}
