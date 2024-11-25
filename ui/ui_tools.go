package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilRandom"
	"image/color"
	"strconv"
)

func (u *Ui) ToolPasswordGenerate(callback ...func(value string)) (content *fyne.Container) {
	var passwords []string

	form := widget.NewForm()
	lengthValue := 0
	useSpecialCharacterValue := false
	lengthInput := NewNumericalEntry()
	lengthInput.SetText("16")
	useSpecialCharacterInput := widget.NewCheck("包含特殊字条", func(b bool) {
		useSpecialCharacterValue = b
	})

	numInput := NewNumericalEntry()
	numInput.SetText("1")
	numValue := 1

	form.AppendItem(widget.NewFormItem("长度", lengthInput))
	form.AppendItem(widget.NewFormItem("", useSpecialCharacterInput))
	form.AppendItem(widget.NewFormItem("数量", numInput))

	passwordView := container.NewVBox()
	passwordViewScroll := NewScrollWithSize(passwordView, 200, 210)

	passwordViewRefresh := func() {
		passwordView.RemoveAll()
		if len(passwords) > 0 {
			for _, password := range passwords {
				callbackButton := container.NewStack()
				if len(callback) > 0 {
					callbackButton.Add(widget.NewButton("使用", func() {
						callback[0](password)
					}))
				}
				passwordView.Add(container.NewBorder(nil, nil, nil, container.NewHBox(IconCopy(func() {
					u.UtilToClipboard(password)
				}), callbackButton), NewLabelWrap(password)))
			}
		}
		passwordView.Refresh()
		passwordViewScroll.Refresh()
	}

	submitButton := widget.NewButton("生成", func() {
		if "" == lengthInput.Text {
			u.DialogError(fmt.Errorf("长度不能为空"))
			return
		}
		var err error
		lengthValue, err = strconv.Atoi(lengthInput.Text)
		if nil != err {
			u.DialogError(fmt.Errorf("长度只能是数字"))
			return
		}
		numValue, err = strconv.Atoi(numInput.Text)
		if nil != err {
			u.DialogError(fmt.Errorf("数量只能是数字"))
			return
		}

		if numValue <= 0 {
			return
		}
		passwords = []string{}
		for i := 0; i < numValue; i++ {
			passwords = append(passwords, utilRandom.RandPassword(lengthValue, !useSpecialCharacterValue))
		}
		sizeWidth := float32(12*lengthValue) + 24
		if len(callback) > 0 {
			sizeWidth += 48
		}
		passwordViewScroll.SetMinSize(fyne.NewSize(sizeWidth, 210))
		passwordViewRefresh()

	})

	content = container.NewVBox()
	content.Add(form)
	content.Add(submitButton)
	content.Add(passwordViewScroll)
	content.Add(canvas.NewText("--------------------------------------------------", color.Transparent))

	return
}
