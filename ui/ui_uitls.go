package ui

func (u *Ui) UtilToClipboard(value string) {
	u.Window.Clipboard().SetContent(value)
}
