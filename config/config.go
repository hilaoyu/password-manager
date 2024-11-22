package config

import (
	"github.com/hilaoyu/go-utils/utils"
	"github.com/hilaoyu/password-manager/ui"
)

var (
	SelfPath  = utils.GetSelfPath()
	UiDefault *ui.Ui
)

func ConfigureUiDefault() (userInterface *ui.Ui) {
	userInterface = ui.NewUi("密码管理器", 800, 600)
	userInterface.Init()
	UiDefault = userInterface
	return
}
