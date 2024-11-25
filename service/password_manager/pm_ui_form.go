package password_manager

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilUuid"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/ui"
	"image/color"
	"slices"
	"strings"
)

type passwordItemExtraItemEntry struct {
	Id    string
	Name  *widget.Entry
	Value *widget.Entry
}

func (pm *PasswordManager) UiPasswordObjectForm(po *PasswordObject) (content *fyne.Container) {
	var formItems []*widget.FormItem
	nameInput := widget.NewEntry()
	descriptionInput := widget.NewEntry()
	descriptionInput.MultiLine = true
	if nil != po {
		nameInput.Text = po.Name
		descriptionInput.Text = po.Description
	}
	submitButton := widget.NewButton("保存", func() {
		name := strings.TrimSpace(nameInput.Text)
		if "" == name {
			config.UiDefault().DialogError(fmt.Errorf("名称不能为空"))
			return
		}
		pm.HandleVerifyPOPassword(po, func() {
			if nil == po {
				po = &PasswordObject{}
			}
			po.Name = nameInput.Text
			po.Description = descriptionInput.Text
			pm.HandleSavePasswordObject(po, true)
		})
	})
	formItems = append(formItems,
		widget.NewFormItem("名称", nameInput),
		widget.NewFormItem("描述", descriptionInput),
	)

	form := widget.NewForm(formItems...)

	content = container.NewVBox(form, submitButton)
	content.Add(canvas.NewText("---------------------------------------------------------------------------------", color.Transparent))
	cancelButton := widget.NewButton("取消", func() {
		config.UiDefault().PrevMainContent()
	})
	content.Add(container.NewCenter(cancelButton))
	return
}

func (pm *PasswordManager) UiPasswordItemForm(pi *PasswordItem, po *PasswordObject, callback func()) (content *fyne.Container) {
	isNew := false
	if nil == pi {
		pi = &PasswordItem{}
		isNew = true
	}
	if "" == pi.Id {
		pi.Id = utilUuid.UuidGenerate()
	}

	var formItems []*widget.FormItem

	nameInput := widget.NewEntry()
	descriptionInput := widget.NewEntry()
	descriptionInput.MultiLine = true
	uriInput := widget.NewEntry()
	accountInput := widget.NewEntry()
	passwordInput := widget.NewPasswordEntry()
	passwordGenerate := widget.NewButton("生成", func() {
		var d dialog.Dialog
		d = config.UiDefault().Dialog("生成密码", config.UiDefault().ToolPasswordGenerate(func(password string) {
			passwordInput.SetText(password)
			d.Hide()
		}))

	})
	if nil != pi {
		nameInput.SetText(pi.Name)
		descriptionInput.SetText(pi.Description)
		uriInput.SetText(pi.Uri)
		accountInput.SetText(pi.Account)
		passwordInput.SetText(pi.Password)
	}

	formItems = append(formItems,
		widget.NewFormItem("名称", nameInput),
		widget.NewFormItem("描述", descriptionInput),
		widget.NewFormItem("Uri", uriInput),
		widget.NewFormItem("账号", accountInput),
		widget.NewFormItem("密码", container.NewBorder(nil, nil, nil, passwordGenerate, passwordInput)),
	)
	form := widget.NewForm(formItems...)

	extraItems := []*passwordItemExtraItemEntry{}
	extraUi := container.NewVBox()
	extraUiItems := container.NewVBox()

	extraUiItemsRefresh := func() {}
	extraUiItemsRefresh = func() {
		extraUiItems.RemoveAll()
		for _, extraItemEntry := range extraItems {
			extraItemDelete := ui.IconDelete(func() {
				extraItems = slices.DeleteFunc(extraItems, func(entry *passwordItemExtraItemEntry) bool {
					if nil == entry {
						return false
					}
					return entry.Id == extraItemEntry.Id
				})
				extraUiItemsRefresh()
			})
			extraItemUi := container.NewBorder(nil, nil, nil, extraItemDelete, container.NewGridWithColumns(2,
				container.NewBorder(nil, nil, widget.NewLabel("名称"), nil, extraItemEntry.Name), container.NewBorder(nil, nil, widget.NewLabel("内容"), nil, extraItemEntry.Value)))
			extraUiItems.Add(extraItemUi)
		}
		extraUiItems.Refresh()
	}

	extraAdd := ui.IconAdd(func() {
		extraItems = append(extraItems, &passwordItemExtraItemEntry{
			Id:    utilUuid.UuidGenerate(),
			Name:  widget.NewEntry(),
			Value: widget.NewEntry(),
		})
		extraUiItemsRefresh()
	})
	extraUi.Add(container.NewBorder(nil, nil, nil, extraAdd, widget.NewLabel("扩展信息")))
	extraUi.Add(extraUiItems)

	if len(pi.Extra) > 0 {
		for _, extraItem := range pi.Extra {
			extraNameInput := widget.NewEntry()
			extraNameInput.SetText(extraItem.Name)
			extraValueInput := widget.NewEntry()
			extraValueInput.SetText(extraItem.Value)
			extraItems = append(extraItems, &passwordItemExtraItemEntry{
				Id:    utilUuid.UuidGenerate(),
				Name:  extraNameInput,
				Value: extraValueInput,
			})
		}
	}
	extraUiItemsRefresh()

	submitButton := widget.NewButton("保存", func() {
		pm.HandleVerifyPOPassword(po, func() {
			pi.Name = nameInput.Text
			pi.Description = descriptionInput.Text
			pi.Uri = uriInput.Text
			pi.Account = accountInput.Text
			pi.Password = passwordInput.Text
			pi.Extra = []*PasswordItemExtra{}
			for _, extraItemEntry := range extraItems {
				pi.Extra = append(pi.Extra, &PasswordItemExtra{
					Name:  extraItemEntry.Name.Text,
					Value: extraItemEntry.Value.Text,
				})
			}

			if isNew {
				po.Passwords = append(po.Passwords, pi)
			}
			callback()
		})
	})

	content = container.NewVBox(form, extraUi, submitButton, canvas.NewText("----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----|----", color.Transparent))

	return
}

func (pm *PasswordManager) UiInputPasswordConfirmForm(callback func(value string)) (content *fyne.Container) {

	form := widget.NewForm()
	passwordInput := widget.NewPasswordEntry()
	passwordConfirmInput := widget.NewPasswordEntry()
	form.AppendItem(widget.NewFormItem("输入密码", passwordInput))
	form.AppendItem(widget.NewFormItem("确认密码", passwordConfirmInput))

	content = container.NewVBox()
	submitButton := widget.NewButton("确定", func() {
		if "" == passwordInput.Text {
			config.UiDefault().DialogError(fmt.Errorf("密码不能为空"))
			return
		}
		if passwordConfirmInput.Text != passwordInput.Text {
			config.UiDefault().DialogError(fmt.Errorf("两次密码不一样"))
			return
		}
		callback(passwordInput.Text)
	})

	content.Add(form)
	content.Add(canvas.NewText("此密码没有任何方式可以找回，请一定要牢记!!!", color.RGBA{
		R: 200,
		G: 0,
		B: 0,
		A: 255,
	}))
	content.Add(submitButton)
	content.Add(canvas.NewText("--------------------------------------------------", color.Transparent))

	return
}

func (pm *PasswordManager) UiInputPasswordForm(callback func(value string)) (content *fyne.Container) {

	form := widget.NewForm()
	passwordInput := widget.NewPasswordEntry()
	form.AppendItem(widget.NewFormItem("输入密码", passwordInput))

	content = container.NewVBox()
	submitButton := widget.NewButton("确定", func() {
		if "" == passwordInput.Text {
			config.UiDefault().DialogError(fmt.Errorf("密码不能为空"))
			return
		}
		callback(passwordInput.Text)
	})
	content.Add(form)
	content.Add(canvas.NewText("--------------------------------------------------", color.Transparent))
	content.Add(submitButton)

	return
}
