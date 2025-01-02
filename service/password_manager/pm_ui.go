package password_manager

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utils"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/tools"
	"github.com/hilaoyu/password-manager/ui"
	"slices"
	"strings"
	"time"
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
			removeIcon := ui.IconRemove(func() {})
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
			removeIcon := ui.IconRemove(func() {
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
			config.UiDefault().DialogError(fmt.Errorf("数据错误，请重新打开"))
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
	head := container.NewHBox(title, menu)

	toolBar := container.NewHBox()
	toolBar.Add(widget.NewButton("生成密码", func() {
		config.UiDefault().NweWindowAndShow("生成密码", tools.ToolPasswordGenerate())
	}))

	c = container.NewBorder(nil, nil, nil, toolBar, head)
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
	poEditButton := ui.IconEdit(func() {
		pm.HandleVerifyPOPassword(po, func() {
			pm.HandleEditPasswordObject(po)
		})
	})

	intro := widget.NewAccordion(widget.NewAccordionItem(fmt.Sprintf("%s", po.SavePath), ui.NewRichTextFromMarkdownWrap(po.Description)))
	saveButton := ui.IconSave(func() {
		pm.HandleVerifyPOPasswordByInput(po, func() {
			newPo, err := po.Clone()
			if nil != err {
				config.UiDefault().DialogError(fmt.Errorf("另存失败,复制数据失败:%v", err))
				return

			}
			newPo.Secret = ""
			newPo.SavePath = ""
			pm.HandleSavePasswordObject(newPo, false)
		})
	})
	content = container.NewVBox(container.NewBorder(nil, nil, container.NewHBox(title, poEditButton, widget.NewSeparator()), saveButton, intro), widget.NewSeparator())
	return
}
func (pm *PasswordManager) UiPasswordObjectToolbar(po *PasswordObject) (content *fyne.Container) {
	content = container.NewVBox()
	searchInput := widget.NewEntry()
	searchInput.SetText(po.searchKeyword)
	searchButton := ui.IconSearch(func() {
		po.searchKeyword = strings.TrimSpace(searchInput.Text)
		pm.HandleViewPasswordObject(po)
		return
	})
	searchClearButton := ui.IconClear(func() {
		po.searchKeyword = ""
		pm.HandleViewPasswordObject(po)
		return
	})

	addButton := ui.IconAdd(func() {
		pm.HandleEditPasswordItem(nil, po)
	})

	content.Add(container.NewGridWithColumns(3, container.NewStack(searchInput), container.NewHBox(searchButton, searchClearButton), container.NewBorder(nil, nil, nil, addButton, nil)))
	content.Add(widget.NewSeparator())

	return
}
func (pm *PasswordManager) UiPasswordObjectPasswords(passwords []*PasswordItem, po *PasswordObject) (content *fyne.Container) {
	if len(passwords) <= 0 {
		content = container.NewStack(widget.NewLabel("没有符合当前条件的密码项"))
		return
	}
	passwordsView := container.NewVBox()
	for _, pi := range passwords {
		passwordsView.Add(pm.UiPasswordItem(pi, po))
	}
	if nil != po.pwScroll {
		offset := po.pwScroll.Offset
		po.pwScroll.Content = passwordsView
		po.pwScroll.Refresh()
		if offset.X > 0 {
			po.pwScroll.Offset.X = offset.X
		}
		if offset.Y > 0 {
			po.pwScroll.Offset.Y = offset.Y
		}

	} else {
		po.pwScroll = ui.NewScrollWithSize(passwordsView, 1000, 400)
		po.pwScroll.Direction = container.ScrollVerticalOnly
	}

	content = container.NewStack(po.pwScroll)
	return
}

func (pm *PasswordManager) UiPasswordItem(pi *PasswordItem, po *PasswordObject) (content *fyne.Container) {
	if nil == pi {
		return container.NewStack(widget.NewLabel("empty"))
	}
	content = container.NewVBox()

	nameUi := widget.NewRichTextFromMarkdown(fmt.Sprintf(`## %s `, pi.Name))
	editButton := ui.IconEdit(func() {
		pm.HandleVerifyPOPassword(po, func() {
			pm.HandleEditPasswordItem(pi, po)
		})
	})
	deleteButton := ui.IconDelete(func() {
		pm.HandleVerifyPOPassword(po, func() {
			po.Passwords = slices.DeleteFunc(po.Passwords, func(item *PasswordItem) bool {
				if nil == item {
					return false
				}
				return item.Id == pi.Id
			})
			pm.HandleSavePasswordObject(po, true)
		})
	})
	content.Add(container.NewBorder(nil, nil, nil, deleteButton, container.NewHBox(nameUi, editButton)))

	uriUi := container.NewBorder(nil, nil, nil, ui.IconCopy(func() {
		config.UiDefault().UtilToClipboard(pi.Uri)
	}), container.NewHBox(
		widget.NewRichTextFromMarkdown("### URI: "),
		ui.NewContainerWithSize(200, 0, ui.NewLabelWrap(pi.Uri)),
	))

	DescUi := container.NewVBox(
		ui.NewContainerWithSize(200, 0, ui.NewLabelWrap(pi.Description)),
	)
	column1 := container.NewVBox(uriUi, DescUi)

	accountUi := container.NewBorder(nil, nil, nil, ui.IconCopy(func() {
		config.UiDefault().UtilToClipboard(pi.Account)
	}), ui.NewContainerWithSize(200, 0, ui.NewLabelWrap(pi.Account)))

	column2 := container.NewVBox(accountUi, pm.UiPassword(pi.Password, po))

	column3 := container.NewVBox()
	if len(pi.Extra) > 0 {
		for _, pe := range pi.Extra {
			peUi := container.NewBorder(nil, nil, ui.NewContainerWithSize(48, 0, ui.NewRichTextFromMarkdownWrap(fmt.Sprintf("### %s: ", pe.Name))), ui.IconCopy(func() {
				config.UiDefault().UtilToClipboard(pe.Value)
			}), ui.NewContainerWithSize(200, 0, ui.NewLabelWrap(pe.Value)))
			column3.Add(peUi)
		}
	}

	grid := container.NewBorder(nil, nil, column1, column3, column2)
	content.Add(grid)
	content.Add(widget.NewSeparator())
	return
}
func (pm *PasswordManager) UiPassword(password string, po *PasswordObject) (content *fyne.Container) {
	if "" == password {
		content = container.NewStack(widget.NewLabel(""))
		return
	}

	var visibilityTimer *time.Timer
	passwordView := ui.NewLabelWrap("******")
	visibilityAction := container.NewStack()

	/*planShow := widget.NewToolbarAction(theme.VisibilityIcon(), func() {})
	planHide := widget.NewToolbarAction(theme.VisibilityOffIcon(), func() {})*/
	var planShow *fyne.Container
	var planHide *fyne.Container

	planHideHandle := func() {
		passwordView.SetText("******")
		visibilityAction.RemoveAll()
		visibilityAction.Add(planShow)
		visibilityAction.Refresh()
		if nil != visibilityTimer {
			visibilityTimer.Stop()
			visibilityTimer = nil
		}
	}
	planShowHandle := func() {
		pm.HandleVerifyPOPassword(po, func() {
			passwordView.SetText(password)
			visibilityAction.RemoveAll()
			visibilityAction.Add(planHide)
			visibilityAction.Refresh()
			visibilityTimer = time.AfterFunc(config.PasswordPlainViveDuration(), func() {
				planHideHandle()
			})
		})

	}
	planShow = ui.IconVisibility(planShowHandle)
	planHide = ui.IconVisibilityOff(planHideHandle)

	visibilityAction.Add(planShow)

	copyAction := ui.IconCopy(func() {
		pm.HandleVerifyPOPassword(po, func() {
			config.UiDefault().UtilToClipboard(password)
		})
	})
	toolbar := container.NewHBox(visibilityAction, copyAction)
	content = container.NewBorder(nil, nil, nil, toolbar, ui.NewContainerWithSize(200, 0, passwordView))

	return
}
