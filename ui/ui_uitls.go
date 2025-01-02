package ui

func (u *Ui) UtilToClipboard(value string) {
	u.WindowMain.Clipboard().SetContent(value)
}
