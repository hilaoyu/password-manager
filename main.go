package main

import (
	"github.com/hilaoyu/password-manager/config"
	"github.com/hilaoyu/password-manager/service/password_manager"
)

func init() {
	config.ConfigureUiDefault()
}

func main() {

	pm := password_manager.NewPasswordManager()
	config.UiDefault.RefreshTop(pm.UiTop())
	config.UiDefault.RefreshMainLeft(pm.UiMenuTree())
	config.UiDefault.RefreshMainContent(pm.UiWelcome())
	config.UiDefault.ShowAndRun()
}
