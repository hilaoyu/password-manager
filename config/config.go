package config

import (
	"github.com/hilaoyu/go-utils/utils"
	"github.com/hilaoyu/password-manager/ui"
	"time"
)

var (
	SelfPath  = utils.GetSelfPath()
	uiDefault *ui.Ui
)

func ConfigureUiDefault() (userInterface *ui.Ui) {
	userInterface = ui.NewUi("密码管理器", 1200, 600)
	userInterface.Init()
	uiDefault = userInterface
	return
}

func UiDefault() *ui.Ui {
	return uiDefault
}

func PasswordVerifyDuration() time.Duration {
	return time.Duration(5) * time.Minute
}

func PasswordPlainViveDuration() time.Duration {
	return time.Duration(10) * time.Second
}
