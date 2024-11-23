package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/service/password_manager"
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
				config.UiDefault().Window.Show()
			}))
		desk.SetSystemTrayMenu(m)
	}

	config.UiDefault().Window.SetCloseIntercept(func() {
		config.UiDefault().Window.Hide()
	})
	config.UiDefault().ShowAndRun()
}
