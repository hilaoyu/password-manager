package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
)

func (u *Ui) DialogError(err error) {
	dialog.NewError(err, u.WindowMain).Show()

	return
}
func (u *Ui) DialogInfo(title string, msg string) {
	dialog.NewInformation(title, msg, u.WindowMain).Show()

	return
}

func (u *Ui) DialogSaveFile(callback func(writer fyne.URIWriteCloser), defaultDir string, defaultName string) {

	d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if nil != err {
			u.DialogError(err)
			return
		}
		callback(writer)
	}, u.WindowMain)

	if "" != defaultDir {
		defaultDirUri, _ := storage.ListerForURI(storage.NewFileURI(defaultDir))
		d.SetLocation(defaultDirUri)
	}
	if "" != defaultName {
		d.SetFileName(defaultName)
	}

	d.Show()

	return
}
func (u *Ui) DialogOpenFile(callback func(reader fyne.URIReadCloser), filter []string, defaultDir string) {

	d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if nil != err {
			u.DialogError(err)
			return
		}
		if nil == reader {
			//pm.UiDialogError(fmt.Errorf("reaber error"))
			return
		}

		callback(reader)

	}, u.WindowMain)

	if len(filter) > 0 {
		d.SetFilter(storage.NewExtensionFileFilter(filter))
	}

	if "" != defaultDir {
		defaultDirUri, _ := storage.ListerForURI(storage.NewFileURI(defaultDir))
		d.SetLocation(defaultDirUri)
	}

	d.Show()

	return
}

func (u *Ui) Dialog(title string, content *fyne.Container) dialog.Dialog {
	d := dialog.NewCustom(title, "取消", content, u.WindowMain)
	d.Show()

	return d
}
