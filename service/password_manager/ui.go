package password_manager

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utils"
	"github.com/hilaoyu/password-manager/config"
	"slices"
	"strings"
)

func (pm *PasswordManager) UiWelcome() (c *fyne.Container) {
	c = container.NewStack(widget.NewLabel("welcome"))
	return
}

func (pm *PasswordManager) UiMenuTree() (c *fyne.Container) {

	tree := widget.NewTree(
		func(id widget.TreeNodeID) []widget.TreeNodeID {
			var ids []widget.TreeNodeID
			switch id {
			case "":
				for _, po := range pm.PasswordObjects {
					if nil != po {
						ids = append(ids, utilEnc.Md5(po.SavePath))
					}
				}
				break
			default:
				break
			}
			return ids
		},
		func(id widget.TreeNodeID) bool {
			if "" == id {
				return true
			}
			return false
		},
		func(branch bool) fyne.CanvasObject {
			label := widget.NewLabel("密码本名称")
			removeIcon := config.UiDefault.IconRemove(func() {})
			if branch {
				label = widget.NewLabel("-密码本名称")
			}
			return container.NewBorder(nil, nil, nil, removeIcon, label)
		},
		func(id widget.TreeNodeID, branch bool, o fyne.CanvasObject) {
			text := "-"
			po := utils.SliceFind(pm.PasswordObjects, func(object *PasswordObject) bool {
				if nil == object {
					return false
				}
				return utilEnc.Md5(object.SavePath) == id
			})
			if nil != po {
				text = po.Name
			}
			if branch {
				text += " (分支)"
			}
			label := widget.NewLabel(text)
			removeIcon := config.UiDefault.IconRemove(func() {
				pm.HandleRemovePasswordObject(po)
			})
			content := container.NewBorder(nil, nil, nil, removeIcon, label)
			*(o.(*fyne.Container)) = *content

		})
	tree.OnSelected = func(id widget.TreeNodeID) {
		po := utils.SliceFind(pm.PasswordObjects, func(object *PasswordObject) bool {
			if nil == object {
				return false
			}
			return utilEnc.Md5(object.SavePath) == id
		})
		if nil == po {
			config.UiDefault.DialogError(fmt.Errorf("数据错误，请重新打开"))
			return
		}
		pm.HandleViewPasswordObject(po)
	}

	pm.menuTree = tree
	c = container.NewStack(tree)
	return
}

func (pm *PasswordManager) UiTop() (c *fyne.Container) {

	menu := container.NewGridWithColumns(2,
		widget.NewButton("打开...", pm.HandleOpenPasswordObject),
		widget.NewButton("新建", func() {
			pm.HandleNewPasswordObject()
		}),
	)
	title := widget.NewRichTextFromMarkdown("# 密码管理")
	split := container.NewHBox(title, menu)

	c = container.NewStack(split)
	return
}

func (pm *PasswordManager) UiPasswordObject(po *PasswordObject) (c *fyne.Container) {
	if nil == po {
		return container.NewStack(widget.NewLabel("empty"))
	}

	var viewPasswords []*PasswordItem
	po.searchKeyword = strings.TrimSpace(po.searchKeyword)
	if "" != po.searchKeyword {
		fmt.Println("po.searchKeyword:", po.searchKeyword)
		viewPasswords = utils.SliceFilter(po.Passwords, func(item *PasswordItem) bool {
			return strings.Contains(item.Name, po.searchKeyword) ||
				strings.Contains(item.Account, po.searchKeyword) ||
				strings.Contains(item.Uri, po.searchKeyword) ||
				strings.Contains(item.Description, po.searchKeyword)

		})
	} else {
		viewPasswords = po.Passwords
	}

	c = container.NewBorder(
		container.NewVBox(pm.UiPasswordObjectTitle(po), pm.UiPasswordObjectToolbar(po)), nil, nil, nil, pm.UiPasswordObjectPasswords(viewPasswords, po))

	return
}
func (pm *PasswordManager) UiPasswordObjectTitle(po *PasswordObject) (content *fyne.Container) {
	if nil == po {
		return container.NewStack(widget.NewLabel("empty"))
	}

	title := widget.NewRichTextFromMarkdown(fmt.Sprintf("# %s", po.Name))
	poEditButton := config.UiDefault.IconEdit(func() {
		pm.HandleEditPasswordObject(po)
	})

	intro := widget.NewAccordion(widget.NewAccordionItem(fmt.Sprintf("%s", po.SavePath), widget.NewRichTextFromMarkdown(po.Description)))
	/*saveButton := config.UiDefault.IconSave(func() {
		newPo, err := po.Clone()
		if nil != err {
			config.UiDefault.DialogError(fmt.Errorf("另存失败,复制数据失败:%v", err))
			return

		}
		newPo.Secret = ""
		newPo.SavePath = ""
		pm.HandleSavePasswordObject(newPo, false)
	})*/
	saveButton := widget.NewButton("另存...", func() {
		newPo, err := po.Clone()
		if nil != err {
			config.UiDefault.DialogError(fmt.Errorf("另存失败,复制数据失败:%v", err))
			return

		}
		newPo.Secret = ""
		newPo.SavePath = ""
		pm.HandleSavePasswordObject(newPo, false)
	})
	content = container.NewVBox(container.NewBorder(nil, nil, container.NewHBox(title, poEditButton, widget.NewSeparator()), saveButton, intro), widget.NewSeparator())
	return
}
func (pm *PasswordManager) UiPasswordObjectToolbar(po *PasswordObject) (content *fyne.Container) {
	content = container.NewVBox()
	searchInput := widget.NewEntry()
	searchInput.SetText(po.searchKeyword)
	searchButton := config.UiDefault.IconSearch(func() {
		po.searchKeyword = strings.TrimSpace(searchInput.Text)
		pm.HandleViewPasswordObject(po)
		return
	})

	addButton := config.UiDefault.IconAdd(func() {
		pm.HandleEditPasswordItem(nil, po)
	})

	content.Add(container.NewGridWithColumns(3, container.NewStack(searchInput), searchButton, container.NewBorder(nil, nil, nil, addButton, nil)))
	content.Add(widget.NewSeparator())

	return
}
func (pm *PasswordManager) UiPasswordObjectPasswords(passwords []*PasswordItem, po *PasswordObject) (content *fyne.Container) {
	if len(passwords) <= 0 {
		content = container.NewStack(widget.NewLabel("没有符合当前条件的密码项"))
		return
	}
	content = container.NewVBox()
	for _, pi := range passwords {
		content.Add(pm.UiPasswordItem(pi, po))
	}
	return
}

func (pm *PasswordManager) UiPasswordItem(pi *PasswordItem, po *PasswordObject) (content *fyne.Container) {
	if nil == pi {
		return container.NewStack(widget.NewLabel("empty"))
	}
	content = container.NewVBox()

	nameUi := widget.NewRichTextFromMarkdown(fmt.Sprintf(`## %s `, pi.Name))
	editButton := config.UiDefault.IconEdit(func() {
		pm.HandleEditPasswordItem(pi, po)
	})
	deleteButton := config.UiDefault.IconDelete(func() {
		po.Passwords = slices.DeleteFunc(po.Passwords, func(item *PasswordItem) bool {
			if nil == item {
				return false
			}
			return item.Id == pi.Id
		})
		pm.HandleSavePasswordObject(po, true)
	})
	content.Add(container.NewBorder(nil, nil, nil, deleteButton, container.NewHBox(nameUi, editButton)))

	uriUi := container.NewHBox(
		widget.NewRichTextFromMarkdown("### URI: "),
		widget.NewLabel(pi.Uri),
		config.UiDefault.IconCopy(pi.Uri),
	)

	DescUi := container.NewHBox(
		widget.NewLabel(pi.Description),
	)
	column1 := container.NewVBox(uriUi, DescUi)

	accountUi := container.NewHBox()
	accountUi.Add(widget.NewLabel(pi.Account))
	accountUi.Add(config.UiDefault.IconCopy(pi.Account))

	column2 := container.NewVBox(accountUi, pm.UiPassword(pi.Password))

	column3 := container.NewVBox()
	if len(pi.Extra) > 0 {
		for _, pe := range pi.Extra {
			peUi := container.NewHBox(
				widget.NewRichTextFromMarkdown(fmt.Sprintf("### %s: ", pe.Name)),
				widget.NewLabel(pe.Value),
				config.UiDefault.IconCopy(pe.Value),
			)
			column3.Add(peUi)
		}
	}

	content.Add(container.NewGridWithColumns(3, column1, column2, column3))
	content.Add(widget.NewSeparator())
	return
}
func (pm *PasswordManager) UiPassword(password string) (content *fyne.Container) {
	if "" == password {
		content = container.NewStack(widget.NewLabel(""))
		return
	}
	content = container.NewHBox()
	passwordView := widget.NewLabel("******      ")
	content.Add(passwordView)

	content.Add(config.UiDefault.IconCopy(password))

	passwordVisibilityButton := widget.NewToolbar()
	passwordVisibilityAction := widget.NewToolbarAction(theme.VisibilityIcon(), func() {})
	passwordVisibilityOffAction := widget.NewToolbarAction(theme.VisibilityOffIcon(), func() {})
	passwordVisibilityButton.Append(passwordVisibilityAction)
	passwordVisibilityAction.OnActivated = func() {
		passwordView.SetText(password)
		passwordVisibilityButton.Items = nil
		passwordVisibilityButton.Append(passwordVisibilityOffAction)
		passwordVisibilityButton.Refresh()

	}
	passwordVisibilityOffAction.OnActivated = func() {
		passwordView.SetText("******      ")
		passwordVisibilityButton.Items = nil
		passwordVisibilityButton.Append(passwordVisibilityAction)
		passwordVisibilityButton.Refresh()

	}
	content.Add(passwordVisibilityButton)

	return
}
