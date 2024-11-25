package password_manager

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hilaoyu/go-utils/utilEnc"
	"github.com/hilaoyu/go-utils/utilTime"
)

type PasswordManager struct {
	menuTree              *widget.Tree
	PasswordObjects       []*PasswordObject
	CurrentPasswordObject *PasswordObject
}

func NewPasswordManager() (pm *PasswordManager) {
	pm = &PasswordManager{
		PasswordObjects: []*PasswordObject{},
	}
	return
}

type PasswordItemExtra struct {
	Name  string ` json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}
type PasswordItem struct {
	Id          string               `json:"id,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Uri         string               `json:"uri,omitempty"`
	Account     string               `json:"account,omitempty"`
	Password    string               `json:"password,omitempty"`
	Extra       []*PasswordItemExtra `json:"extra,omitempty"`
}
type PasswordObject struct {
	Name          string          `json:"name,omitempty"`
	Description   string          `json:"description,omitempty"`
	SavePath      string          `json:"-"`
	Secret        string          `json:"-"`
	Passwords     []*PasswordItem `json:"passwords,omitempty"`
	searchKeyword string
	verifyTimer   *utilTime.Timer
	pwScroll      *container.Scroll
}

func (po *PasswordObject) Encode() (enData []byte, err error) {

	if "" == po.Secret {
		err = fmt.Errorf("未设置加密密钥")
		return
	}
	jsonByte, err := json.Marshal(po)
	if nil != err {
		return
	}
	encryptor := utilEnc.NewAesEncryptor(po.Secret)
	ivLength, err := encryptor.GetBlockSize()
	if nil != err {
		return
	}

	iv := []byte(po.Secret)[:ivLength]
	enData, err = encryptor.EncryptByte(jsonByte, iv)

	return
}
func (po *PasswordObject) Clone() (passwordObject *PasswordObject, err error) {

	jsonByte, err := json.Marshal(po)
	if nil != err {
		return
	}
	passwordObject = &PasswordObject{}
	err = json.Unmarshal(jsonByte, passwordObject)
	return
}
