package password_manager

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/password-manager/config"
	"os"
	"runtime"
	"slices"
)

func (pm *PasswordManager) HandleAppendAndViewPasswordObject(po *PasswordObject) {
	config.UiDefault.RefreshMainLeft(pm.UiMenuTree())
	if nil == po {
		config.UiDefault.RefreshMainContent(pm.UiWelcome())
		return
	}
	menuKey := utilEnc.Md5(po.SavePath)
	pm.menuTree.Select(menuKey)
	return
}
func (pm *PasswordManager) HandleOpenPasswordObject() {
	defaultPath := ""
	if "darwin" != runtime.GOOS {
		defaultPath = config.SelfPath
	}
	config.UiDefault.DialogOpenFile(func(reader fyne.URIReadCloser) {
		defer reader.Close()
		enData, err := os.ReadFile(reader.URI().Path())
		if nil != err {
			config.UiDefault.DialogError(err)
			return
		}

		var d dialog.Dialog
		d = config.UiDefault.Dialog("请输入密码", pm.UiInputPasswordForm(func(value string) {
			if "" == value {
				config.UiDefault.DialogError(fmt.Errorf("密码不能为空"))
				return
			}
			d.Hide()
			secret := UtilPasswordToSecret(value)

			encryptor := utilEnc.NewAesEncryptor(secret)
			var ivLength int
			ivLength, err = encryptor.GetBlockSize()
			if nil != err {
				config.UiDefault.DialogError(err)
				return
			}
			if ivLength <= 0 {
				config.UiDefault.DialogError(fmt.Errorf("密码错误"))
				return
			}

			iv := []byte(secret)[:ivLength]
			var deData []byte
			deData, err = encryptor.DecryptByte(enData, iv)
			if nil != err {
				config.UiDefault.DialogError(fmt.Errorf("密码错误! "))
				return
			}

			po := &PasswordObject{}
			err = json.Unmarshal(deData, po)
			if nil != err {
				config.UiDefault.DialogError(fmt.Errorf("密码错误!! "))
				return
			}

			po.Secret = secret
			po.SavePath = reader.URI().Path()
			pm.PasswordObjects = append(pm.PasswordObjects, po)
			pm.HandleAppendAndViewPasswordObject(po)
		}))

	}, []string{".pwe"}, defaultPath)

}
func (pm *PasswordManager) HandleNewPasswordObject() {
	config.UiDefault.RefreshMainContent(pm.UiPasswordObjectForm(nil))
}
func (pm *PasswordManager) HandleEditPasswordObject(po *PasswordObject) {
	config.UiDefault.RefreshMainContent(pm.UiPasswordObjectForm(po))
}
func (pm *PasswordManager) HandleViewPasswordObject(po *PasswordObject) {

	if nil == po {
		config.UiDefault.RefreshMainContent(widget.NewLabel("密码本不存在"))
	}
	config.UiDefault.RefreshMainContent(pm.UiPasswordObject(po))
}
func (pm *PasswordManager) HandleSavePasswordObject(po *PasswordObject, openNow bool) {
	if nil == po {
		config.UiDefault.DialogError(fmt.Errorf("密码本为空"))
		return
	}

	if "" == po.Secret {
		var d dialog.Dialog
		d = config.UiDefault.Dialog("设置密码", pm.UiInputPasswordConfirmForm(func(value string) {
			po.Secret = UtilPasswordToSecret(value)
			d.Hide()
			pm.HandleSavePasswordObject(po, openNow)
		}))
		return
	}

	enData, err := po.Encode()
	if nil != err {
		config.UiDefault.DialogError(err)
		return
	}

	if "" == po.SavePath {
		defaultPath := ""
		if "darwin" != runtime.GOOS {
			defaultPath = config.SelfPath
		}
		config.UiDefault.DialogSaveFile(func(writer fyne.URIWriteCloser) {
			defer writer.Close()
			_, err = writer.Write(enData)
			if nil != err {
				config.UiDefault.DialogError(err)
				return
			}
			po.SavePath = writer.URI().Path()
			if openNow {
				pm.PasswordObjects = append(pm.PasswordObjects, po)
				pm.HandleAppendAndViewPasswordObject(po)
			}

			config.UiDefault.DialogInfo("提示", "保存成功")
		}, defaultPath, po.Name+".pwe")
		return
	}

	err = os.WriteFile(po.SavePath, enData, 0666)
	if nil != err {
		config.UiDefault.DialogError(err)
		return
	}
	if openNow {
		pm.HandleAppendAndViewPasswordObject(po)
	}

}
func (pm *PasswordManager) HandleRemovePasswordObject(po *PasswordObject) {
	if nil == po {
		return
	}
	if len(pm.PasswordObjects) <= 0 {
		return
	}
	slices.DeleteFunc(pm.PasswordObjects, func(object *PasswordObject) bool {
		if nil == object {
			return false
		}
		return object.SavePath == po.SavePath
	})
	var firstPo *PasswordObject
	if len(pm.PasswordObjects) > 0 {
		firstPo = pm.PasswordObjects[0]
	}
	pm.HandleAppendAndViewPasswordObject(firstPo)
}

func (pm *PasswordManager) HandleEditPasswordItem(pi *PasswordItem, po *PasswordObject) {
	var d dialog.Dialog
	d = config.UiDefault.Dialog("密码详情", pm.UiPasswordItemForm(pi, po, func() {
		pm.HandleSavePasswordObject(po, true)
		d.Hide()
	}))

}
