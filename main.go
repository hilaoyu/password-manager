package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/service/password_manager"
	"github.com/hilaoyu/password-manager/tools"
)

func init() {
	config.ConfigureUiDefault()
}

func main() {

	pm := password_manager.NewPasswordManager()
	config.UiDefault().RefreshTop(pm.UiTop())
	config.UiDefault().RefreshMainLeft(pm.UiMenuTree())
	config.UiDefault().RefreshMainContent(pm.UiWelcome())
	pm.ListenFileDropIn()
	if desk, ok := config.UiDefault().App.(desktop.App); ok {
		m := fyne.NewMenu("PM",
			fyne.NewMenuItem("显示", func() {
				config.UiDefault().WindowMain.Show()
			}),
			fyne.NewMenuItem("生成密码", func() {
				config.UiDefault().NweWindowAndShow("生成密码", tools.ToolPasswordGenerate())
			}),
			fyne.NewMenuItem("AES加解密", func() {
				config.UiDefault().NweWindowAndShow("AES加解密", tools.ToolAesEncrypt())
			}),
			fyne.NewMenuItem("AES加密请求", func() {
				config.UiDefault().NweWindowAndShow("AES加密请求", tools.ToolAesRequest())
			}),
		)
		desk.SetSystemTrayMenu(m)
	}

	config.UiDefault().WindowMain.SetCloseIntercept(func() {
		config.UiDefault().WindowMain.Hide()
	})
	config.UiDefault().ShowAndRun()
}
