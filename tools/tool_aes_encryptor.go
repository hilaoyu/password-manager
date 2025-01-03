package tools

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/ui"
	"strings"
)

func ToolAesEncrypt() (content *fyne.Container) {

	secretInput := widget.NewEntry()
	plainDataInput := widget.NewMultiLineEntry()
	plainDataInput.Wrapping = fyne.TextWrapBreak
	encDataInput := widget.NewMultiLineEntry()
	encDataInput.Wrapping = fyne.TextWrapBreak

	enButton := widget.NewButton("加密>>", func() {
		secretInput.Text = strings.TrimSpace(secretInput.Text)
		if "" == secretInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("密钥不能为空"))
			return
		}

		plainDataInput.Text = strings.TrimSpace(plainDataInput.Text)
		if "" == plainDataInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("明文不能为空"))
			return
		}

		encryptor := utilEnc.NewAesEncryptor(secretInput.Text)

		var err error
		encDataInput.Text, err = encryptor.EncryptString(plainDataInput.Text)
		encDataInput.Refresh()
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("加密失败"))
			return
		}

	})

	deButton := widget.NewButton("<<解密", func() {
		secretInput.Text = strings.TrimSpace(secretInput.Text)
		if "" == secretInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("密钥不能为空"))
			return
		}

		encDataInput.Text = strings.TrimSpace(encDataInput.Text)
		if "" == encDataInput.Text {
			config.UiDefault().WindowError(fmt.Errorf("密文不能为空"))
			return
		}

		encryptor := utilEnc.NewAesEncryptor(secretInput.Text)

		var err error
		plainDataInput.Text, err = encryptor.DecryptString(encDataInput.Text)
		plainDataInput.Refresh()
		if nil != err {
			config.UiDefault().WindowError(fmt.Errorf("解密密失败"))
			return
		}

	})

	content = container.NewVBox()
	content.Add(widget.NewForm(widget.NewFormItem("密钥:", secretInput)))
	content.Add(container.NewHBox(
		container.NewVBox(
			widget.NewLabel("明文"),
			ui.NewContainerWithSize(400, 300, plainDataInput),
		),
		container.NewVBox(
			ui.NewContainerWithSize(0, 60),
			enButton,
			ui.NewContainerWithSize(0, 60),
			deButton,
		),
		container.NewVBox(
			widget.NewLabel("密文"),
			ui.NewContainerWithSize(400, 300, encDataInput),
		),
	))

	return
}
